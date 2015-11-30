package diff

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var (
	rootPath      = "$root"
	DefaultConfig = &DiffConfig{AllowUnexpectedKeys: true}
)

// During deepValueEqual, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  uintptr
	a2  uintptr
	typ reflect.Type
}

type DiffConfig struct {
	AllowUnexpectedKeys bool
}

type Differences []*Mismatch

func (d *Differences) Append(m *Mismatch) {
	*d = append(*d, m)
}

func (d *Differences) toString() []string {
	s := make([]string, len(*d))
	for _, m := range *d {
		s = append(s, fmt.Sprint(m))
	}
	return s
}

func (d Differences) Error() string {
	return strings.Join(d.toString(), "\n")
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqual(path string, v1, v2 reflect.Value, visited map[visit]bool, depth int, d *Differences, conf *DiffConfig) (ok bool) {
	mismatchf := func(typ mismatchType, a ...interface{}) {
		d.Append(newMismatch(v1, v2, path, typ, a...))
	}

	if !v1.IsValid() || !v2.IsValid() {
		if v1.IsValid() == v2.IsValid() {
			return true
		}
		mismatchf(mValidty)
		return false
	}

	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := v1.UnsafeAddr()
		addr2 := v2.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are identical ...
		if addr1 == addr2 {
			return true
		}

		// ... or already seen
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	v1Kind := v1.Kind()
	//Do type check only when object is not a structure, so we can get deep diff
	if v1Kind != reflect.Struct {
		if v1.Type() != v2.Type() {
			mismatchf(mType, v1.Type(), v2.Type())
			return false
		}
	}

	switch v1Kind {
	case reflect.Array:
		if v1.Len() != v2.Len() {
			// can't happen!
			mismatchf(mLen, v1.Len(), v2.Len())
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if ok := deepValueEqual(
				fmt.Sprintf("%s[%d]", path, i),
				v1.Index(i), v2.Index(i), visited, depth+1, d, conf); !ok {
				return false
			}
		}
		return true
	case reflect.Slice:
		// We treat a nil slice the same as an empty slice.
		if v1.Len() != v2.Len() {
			mismatchf(mLen, v1.Len(), v2.Len())
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if ok := deepValueEqual(
				fmt.Sprintf("%s[%d]", path, i),
				v1.Index(i), v2.Index(i), visited, depth+1, d, conf); !ok {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			if v1.IsNil() != v2.IsNil() {
				mismatchf(mNilVsNonNil)
				return false
			}
			return true
		}
		return deepValueEqual(path, v1.Elem(), v2.Elem(), visited, depth+1, d, conf)
	case reflect.Ptr:
		return deepValueEqual("(*"+path+")", v1.Elem(), v2.Elem(), visited, depth+1, d, conf)
	case reflect.Struct:
		if v1.NumField() > v2.NumField() {
			mismatchf(mLen, v1.NumField(), v2.NumField())
		} else if v2.NumField() > v1.NumField() && conf.AllowUnexpectedKeys == false {
			mismatchf(mLen, v1.NumField(), v2.NumField())
			return false
		}

		result := true
		for i, n := 0, v1.NumField(); i < n; i++ {
			fieldName := v1.Type().Field(i).Name
			path := path + `["` + fieldName + `"]`

			fieldNotFound := true
			for x, m := 0, v2.NumField(); x < m; x++ {
				if fieldName == v2.Type().Field(x).Name {
					fieldNotFound = false
					break
				}
			}

			if fieldNotFound {
				mismatchf(mFieldNotFound, path)
				result = false
			} else if ok := deepValueEqual(path, v1.Field(i), v2.Field(i), visited, depth+1, d, conf); !ok {
				result = false
			}
		}
		return result
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			mismatchf(mNilVsNonNil)
			return false
		}
		if v1.Len() > v2.Len() {
			mismatchf(mLen, v1.Len(), v2.Len())
		} else if v2.Len() > v1.Len() && conf.AllowUnexpectedKeys == false {
			mismatchf(mLen, v1.Len(), v2.Len())
			return false
		}

		if v1.Pointer() == v2.Pointer() {
			return true
		}

		result := true
		for _, v1k := range v1.MapKeys() {
			var p string
			if v1k.CanInterface() {
				p = path + "[" + fmt.Sprintf("%#v", v1k.Interface()) + "]"
			} else {
				p = path + "[someKey]"
			}

			var keyFound interface{}

			for _, v2k := range v2.MapKeys() {
				if reflect.DeepEqual(interfaceOf(v1k), interfaceOf(v2k)) {
					keyFound = interfaceOf(v1k)
					break
				}
			}

			if keyFound == nil {
				mismatchf(mKeyNotFound, p)
				result = false
			} else if ok := deepValueEqual(p, v1.MapIndex(v1k), v2.MapIndex(v1k), visited, depth+1, d, conf); !ok {
				result = false
			}
		}
		return result
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		mismatchf(mNonNilFunc)
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v1.Int() != v2.Int() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v1.Uint() != v2.Uint() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.Float32, reflect.Float64:
		if v1.Float() != v2.Float() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.Complex64, reflect.Complex128:
		if v1.Complex() != v2.Complex() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.String:
		if v1.String() != v2.String() {
			mismatchf(mUnequal)
			return false
		}
		return true
	case reflect.Chan, reflect.UnsafePointer:
		if v1.Pointer() != v2.Pointer() {
			mismatchf(mUnequal)
			return false
		}
		return true
	default:
		panic("unexpected type " + v1.Type().String())
	}
}

// DeepDiff tests for deep equality. It uses normal == equality where
// possible but will scan elements of arrays, slices, maps, and fields
// of structs. In maps, keys are compared with == but elements use deep
// equality. DeepDiff correctly handles recursive types. Functions are
// equal only if they are both nil.
//
// DeepDiff differs from reflect.DeepDiff in that an empty slice is
// equal to a nil slice. Additional fileds/keys from structs/maps of a2
// are not treated as mismatch as this matching needs to be loose for
// non
func DeepDiff(a1, a2 interface{}, conf *DiffConfig) (bool, Differences) {
	var d Differences
	if conf == nil {
		conf = DefaultConfig
	}

	mismatchf := func(typ mismatchType, a ...interface{}) {
		d.Append(newMismatch(reflect.ValueOf(a1), reflect.ValueOf(a2), rootPath, typ, a...))
	}

	if a1 == nil || a2 == nil {
		if a1 == a2 {
			return true, nil
		}
		mismatchf(mNilVsNonNil)
		return false, d
	}
	v1 := reflect.ValueOf(a1)
	v2 := reflect.ValueOf(a2)

	return deepValueEqual(rootPath, v1, v2, make(map[visit]bool), 0, &d, conf), d
}

// interfaceOf returns v.Interface() even if v.CanInterface() == false.
// This enables us to call fmt.Printf on a value even if it's derived
// from inside an unexported field.
// See https://code.google.com/p/go/issues/detail?id=8965
// for a possible future alternative to this hack.
func interfaceOf(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}
	return bypassCanInterface(v).Interface()
}

type flag uintptr

var flagRO flag

// constants copied from reflect/value.go
const (
	// The value of flagRO up to and including Go 1.3.
	flagRO1p3 = 1 << 0

	// The value of flagRO from Go 1.4.
	flagRO1p4 = 1 << 5
)

var flagValOffset = func() uintptr {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	return field.Offset
}()

func flagField(v *reflect.Value) *flag {
	return (*flag)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + flagValOffset))
}

// bypassCanInterface returns a version of v that
// bypasses the CanInterface check.
func bypassCanInterface(v reflect.Value) reflect.Value {
	if !v.IsValid() || v.CanInterface() {
		return v
	}
	*flagField(&v) &^= flagRO
	return v
}

// Sanity checks against future reflect package changes
// to the type or semantics of the Value.flag field.
func init() {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	if field.Type.Kind() != reflect.TypeOf(flag(0)).Kind() {
		panic("reflect.Value flag field has changed kind")
	}
	var t struct {
		a int
		A int
	}
	vA := reflect.ValueOf(t).FieldByName("A")
	va := reflect.ValueOf(t).FieldByName("a")
	flagA := *flagField(&vA)
	flaga := *flagField(&va)

	// Infer flagRO from the difference between the flags
	// for the (otherwise identical) fields in t.
	flagRO = flagA ^ flaga
	if flagRO != flagRO1p3 && flagRO != flagRO1p4 {
		panic("reflect.Value read-only flag has changed semantics")
	}
}

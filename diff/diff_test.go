package diff

import (
	"reflect"
	"testing"
)

type Basic struct {
	x int
	y float32
}

type BasicV2 struct {
	x int
	y float32
	s string
}

type DiffEqualTest struct {
	a, b interface{}
	eq   bool
	msg  string
}

type DiffInequalTest struct {
	a, b interface{}
	diff *Mismatch
}

// Simple functions for DeepDiff tests.
var (
	fn1 func()             // nil.
	fn2 func()             // nil.
	fn3 = func() { fn1() } // Not nil.
)

var equalTests = []DiffEqualTest{
	// Equalities
	{nil, nil, true, ""},
	{1, 1, true, ""},
	{int32(1), int32(1), true, ""},
	{0.5, 0.5, true, ""},
	{float32(0.5), float32(0.5), true, ""},
	{"hello", "hello", true, ""},
	{make([]int, 10), make([]int, 10), true, ""},
	{&[3]int{1, 2, 3}, &[3]int{1, 2, 3}, true, ""},
	{Basic{1, 0.5}, Basic{1, 0.5}, true, ""},
	{Basic{1, 2}, BasicV2{1, 2, "text"}, true, ""},
	{error(nil), error(nil), true, ""},
	{map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, true, ""},
	{map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one", 3: "three"}, true, ""},
	{fn1, fn2, true, ""},

	// Nil vs empty: they're the same (difference from normal DeepDiff)
	{[]int{}, []int(nil), true, ""},
	{[]int{}, []int{}, true, ""},
	{[]int(nil), []int(nil), true, ""},
}

var inequalTests = []*DiffInequalTest{
	newInequalTest(1, 2, 1, 2, rootPath, mUnequal),
	newInequalTest(int32(1), int32(2), int32(1), int32(2), rootPath, mUnequal),
	newInequalTest(0.5, 0.6, 0.5, 0.6, rootPath, mUnequal),
	newInequalTest(float32(0.5), float32(0.6), float32(0.5), float32(0.6), rootPath, mUnequal),
	newInequalTest("hello", "hey", "hello", "hey", rootPath, mUnequal),
	newInequalTest(make([]int, 10), make([]int, 11), make([]int, 10), make([]int, 11), rootPath, mLen, 10, 11),
	newInequalTest(&[3]int{1, 2, 3}, &[3]int{1, 2, 4}, 3, 4, "(*"+rootPath+")[2]", mUnequal),
	newInequalTest(Basic{1, 0.5}, Basic{1, 0.6}, 0.5, 0.6, rootPath+"[\"y\"]", mUnequal),
	newInequalTest(Basic{1, 0}, Basic{2, 0}, 1, 2, rootPath+"[\"x\"]", mUnequal),
	newInequalTest(Basic{1, 2}, BasicV2{1, 2, "text"}, Basic{1, 2}, BasicV2{1, 2, "text"}, rootPath, mLen, 2, 3),
	newInequalTest(map[int]string{1: "one", 3: "two"}, map[int]string{2: "two", 1: "one"}, map[int]string{1: "one", 3: "two"}, map[int]string{2: "two", 1: "one"}, rootPath, mKeyNotFound, rootPath+"[3]"),
	newInequalTest(map[int]string{1: "one", 2: "txo"}, map[int]string{2: "two", 1: "one"}, "txo", "two", rootPath+"[2]", mUnequal),
	newInequalTest(map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one", 3: "three"}, map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one", 3: "three"}, rootPath, mLen, 2, 3),
	newInequalTest(nil, 1, nil, 1, rootPath, mNilVsNonNil),
	newInequalTest(fn1, fn3, fn1, fn3, rootPath, mNonNilFunc),
	newInequalTest([]interface{}{nil}, []interface{}{"a"}, nil, "a", rootPath+"[0]", mNilVsNonNil),
	newInequalTest(1, 1.0, 1, 1.0, rootPath, mType, reflect.TypeOf(1), reflect.TypeOf(1.0)),
	newInequalTest([]int{1, 2, 3}, [3]int{1, 2, 3}, []int{1, 2, 3}, [3]int{1, 2, 3}, rootPath, mType, "[]int", "[3]int"),
}

func newInequalTest(a, b, diffa, diffb interface{}, path string, typ mismatchType, typMsgArgs ...interface{}) *DiffInequalTest {
	return &DiffInequalTest{
		a:    a,
		b:    b,
		diff: newMismatch(reflect.ValueOf(diffa), reflect.ValueOf(diffb), path, typ, typMsgArgs...),
	}
}

func TestEqualities(t *testing.T) {
	t.Parallel()
	for _, test := range equalTests {
		r, diffs := DeepDiff(test.a, test.b, nil)
		if r != test.eq {
			t.Errorf("DeepDiff(%v, %v) = %v, want %v", test.a, test.b, r, test.eq)
		}
		if test.eq {
			if diffs != nil {
				t.Errorf("DeepDiff(%v, %v): unexpected error message %q when equal", test.a, test.b, diffs[0])
			}
		} else {
			if test.msg != diffs[0].String() {
				t.Errorf("DeepDiff(%v, %v); unexpected error %q, want %q", test.a, test.b, diffs[0], test.msg)
			}
		}
	}
}

func TestInequalities(t *testing.T) {
	t.Parallel()
	for _, test := range inequalTests {
		r, diffs := DeepDiff(test.a, test.b, &DiffConfig{AllowUnexpectedKeys: false})
		if r != false {
			t.Errorf("DeepDiff(%v, %v) = %v, want %v", test.a, test.b, r, false)
			continue
		}

		if test.diff.Reason() != diffs[0].Reason() {
			t.Errorf("DeepDiff(%v, %v); \ngot: %s, \nwant %s", test.a, test.b, diffs[0], test.diff)
		}
	}
}

type Recursive struct {
	x int
	r *Recursive
}

func TestDeepDiffRecursiveStruct(t *testing.T) {
	a, b := new(Recursive), new(Recursive)
	*a = Recursive{12, a}
	*b = Recursive{12, b}
	if ok, _ := DeepDiff(a, b, nil); !ok {
		t.Error("DeepDiff(recursive same) = false, want true")
	}
}

type _Complex struct {
	a int
	b [3]*_Complex
	c *string
	d map[float64]float64
}

func TestDeepDiffComplexStruct(t *testing.T) {
	m := make(map[float64]float64)
	stra, strb := "hello", "hello"
	a, b := new(_Complex), new(_Complex)
	*a = _Complex{5, [3]*_Complex{a, b, a}, &stra, m}
	*b = _Complex{5, [3]*_Complex{b, a, a}, &strb, m}
	if ok, _ := DeepDiff(a, b, nil); !ok {
		t.Error("DeepDiff(complex same) = false, want true")
	}
}

func TestDeepDiffComplexStructInequality(t *testing.T) {
	m := make(map[float64]float64)
	stra, strb := "hello", "helloo" // Difference is here
	a, b := new(_Complex), new(_Complex)
	*a = _Complex{5, [3]*_Complex{a, b, a}, &stra, m}
	*b = _Complex{5, [3]*_Complex{b, a, a}, &strb, m}
	if ok, _ := DeepDiff(a, b, nil); ok {
		t.Error("DeepDiff(complex different) = true, want false")
	}
}

type UnexpT struct {
	m map[int]int
}

func TestDeepDiffUnexportedMap(t *testing.T) {
	// Check that DeepDiff can look at unexported fields.
	x1 := UnexpT{map[int]int{1: 2}}
	x2 := UnexpT{map[int]int{1: 2}}
	if ok, _ := DeepDiff(&x1, &x2, nil); !ok {
		t.Error("DeepDiff(x1, x2) = false, want true")
	}

	y1 := UnexpT{map[int]int{2: 3}}
	if ok, _ := DeepDiff(&x1, &y1, nil); ok {
		t.Error("DeepDiff(x1, y1) = true, want false")
	}
}

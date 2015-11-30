package diff

import (
	"fmt"
	"reflect"
)

type mismatchType int

type Mismatch struct {
	v1, v2 reflect.Value
	path   string
	how    string
	typ    mismatchType
}

const (
	mType mismatchType = iota
	mLen
	mUnequal
	mValidty
	mFieldUnexpected
	mFieldNotFound
	mKeyNotFound
	mKeyUnexpected
	mNilVsNonNil
	mNonNilFunc
)

var (
	diffMsg       = "mismatch at %s: %s \nexpected: \n\t%#v \nrecieved \n\t%#v"
	diffMsgReason = "mismatch at %s: %s"
)

var typeMsgs = map[mismatchType]string{
	mType:            "type mismatch expected %s recieved %s",
	mLen:             "length mismatch, expected %d recieved %d",
	mUnequal:         "unequal",
	mValidty:         "validity mismatch",
	mFieldUnexpected: "unexpected field %s",
	mFieldNotFound:   "field %s not found",
	mKeyNotFound:     "key %s not found",
	mKeyUnexpected:   "unexpected key %s",
	mNilVsNonNil:     "nil vs non-nil mismatch",
	mNonNilFunc:      "non-nil functions",
}

func newMismatch(v1, v2 reflect.Value, path string, typ mismatchType, typMsgArgs ...interface{}) *Mismatch {
	return &Mismatch{
		v1:   v1,
		v2:   v2,
		path: path,
		typ:  typ,
		how:  fmt.Sprintf(typeMsgs[typ], typMsgArgs...),
	}
}

func (m *Mismatch) String() string {
	return fmt.Sprintf(diffMsg, m.path, m.how, interfaceOf(m.v1), interfaceOf(m.v2))
}

func (m *Mismatch) Reason() string {
	return fmt.Sprintf(diffMsgReason, m.path, m.how)
}

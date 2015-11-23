package diff

import (
	"fmt"
	"reflect"
)

type mismatchType int

type mismatch struct {
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
	diffMsg       = "mismatch at %s: %s \nexpected: \n\t%#v \nobtained \n\t%#v"
	diffMsgReason = "mismatch at %s: %s"
)

var typeMsgs = map[mismatchType]string{
	mType:            "type mismatch %s vs %s",
	mLen:             "length mismatch, %d vs %d",
	mUnequal:         "unequal",
	mValidty:         "validity mismatch",
	mFieldUnexpected: "unexpected field %s",
	mFieldNotFound:   "field %s not found",
	mKeyNotFound:     "key %s not found",
	mKeyUnexpected:   "unexpected key %s",
	mNilVsNonNil:     "nil vs non-nil mismatch",
	mNonNilFunc:      "non-nil functions",
}

func newMismatch(v1, v2 reflect.Value, path string, typ mismatchType, typMsgArgs ...interface{}) *mismatch {
	return &mismatch{
		v1:   v1,
		v2:   v2,
		path: path,
		typ:  typ,
		how:  fmt.Sprintf(typeMsgs[typ], typMsgArgs...),
	}
}

func (m *mismatch) String() string {
	return fmt.Sprintf(diffMsg, m.path, m.how, interfaceOf(m.v1), interfaceOf(m.v2))
}

func (m *mismatch) Reason() string {
	return fmt.Sprintf(diffMsgReason, m.path, m.how)
}

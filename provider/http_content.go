package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type httpContent interface {
	GetData() ([]byte, error)
	GetBody() interface{}
	SetBody(content interface{}) error
}

type jsonContent struct {
	data      map[string]interface{}
	sliceData []interface{}
}

func (c *jsonContent) GetData() ([]byte, error) {
	if len(c.data) > 0 {
		return json.Marshal(c.data)
	} else if len(c.sliceData) > 0 {
		return json.Marshal(c.sliceData)
	} else {
		return nil, nil
	}
}

func (c *jsonContent) GetBody() interface{} {
	if c.data != nil {
		return c.data
	} else if c.sliceData != nil {
		return c.sliceData
	} else {
		return nil
	}
}

func (c *jsonContent) SetBody(content interface{}) error {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		return c.setJSONStringBody(v.String())
	case reflect.Map, reflect.Struct:
		return c.setStructBody(v.Interface())
	case reflect.Slice:
		c.setSliceBody(v)
	default:
		return errors.New("content is not valid json")
	}
	return nil
}

func (c *jsonContent) setJSONStringBody(content string) error {
	if content == "" {
		return nil
	}

	var val interface{}
	d := json.NewDecoder(strings.NewReader(content))
	d.UseNumber()
	if err := d.Decode(&val); err != nil {
		return err
	}
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Map:
		return c.setStructBody(val)
	case reflect.Slice:
		c.setSliceBody(v)
	default:
		return errors.New("conent is not valid json")
	}
	return nil
}

func (c *jsonContent) setStructBody(content interface{}) error {
	if marshalContent, err := json.Marshal(content); err != nil {
		return err
	} else {
		c.data = make(map[string]interface{})
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&c.data); err != nil {
			return err
		}
	}
	return nil
}

func (c *jsonContent) setSliceBody(v reflect.Value) {
	c.sliceData = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		c.sliceData[i] = v.Index(i).Interface()
	}
}

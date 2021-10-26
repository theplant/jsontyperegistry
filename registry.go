package jsontyperegistry

import (
	"encoding/json"
	"errors"
	"reflect"
)

var reg = New()

type obj struct {
	Type  string
	Value interface{}
}

func MustRegisterType(v interface{}) {
	err := reg.Register(v)
	if errors.Is(err, ErrDuplicateEntry) {
		err = nil
	}

	if err != nil {
		panic(err)
	}
}

func MustJSONString(v interface{}) string {

	t := GetLongTypeName(v)

	r, err := json.Marshal(&obj{
		Type:  t,
		Value: v,
	})

	if err != nil {
		panic(err)
	}
	return string(r)
}

type obj2 struct {
	Type  string
	Value json.RawMessage
}

func MustNewWithJSONString(v string) interface{} {
	var val obj2
	err := json.Unmarshal([]byte(v), &val)
	if err != nil {
		panic(err)
	}
	t, err := reg.GetType(val.Type)
	if err != nil {
		panic(err)
	}

	i := reflect.New(t).Interface()

	err = json.Unmarshal(val.Value, i)
	if err != nil {
		panic(err)
	}
	return reflect.ValueOf(i).Elem().Interface()
}

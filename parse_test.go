package go_json_parse_example

import (
	"fmt"
	"reflect"
	"testing"
)

func TestName(t *testing.T) {
	type args struct {
		s    string
		want interface{}
	}

	for idx, v := range []args{
		{s: `null`, want: nil},
		{s: ` null `, want: nil},
		{s: `"str"`, want: "str"},
		{s: ` "str  "   `, want: "str  "},
		{s: ` "s\"tr  "   `, want: "s\"tr  "},
		{s: ` true `, want: true},
		{s: `  false   `, want: false},
		{s: ` [ ] `, want: []interface{}{}},
		{s: `[ "str  "   , false, null, {"a":"a", "b":1, "c": -1, "d": null, "e":[1, false]} ]`, want: []interface{}{"str  ", false, nil, map[string]interface{}{"a": "a", "b": int64(1), "c": int64(-1), "d": nil, "e": []interface{}{int64(1), false}}}},
		{s: `{}`, want: map[string]interface{}{}},
		{s: `{"a":"a", "b":1, "c": -1, "d": null, "e":[1, false]}`, want: map[string]interface{}{"a": "a", "b": int64(1), "c": int64(-1), "d": nil, "e": []interface{}{int64(1), false}}},
	} {
		parser := &jsonParser{data: []rune(v.s), idx: 0}
		parser.removeSpace()
		data, err := parser.parse()
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(data, v.want) {
			panic(fmt.Sprintf("%d - %s should get %#v, but got %#v", idx, v.s, v.want, data))
		}
	}
}

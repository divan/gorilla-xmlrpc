// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"fmt"
	"reflect"
	"strings"
)

func RPCRequest2XML(method string, rpc interface{}) (string, error) {
	buffer := "<methodCall><methodName>"
	buffer += method
	buffer += "</methodName>"
	params, err := RPCParams2XML(rpc)
	buffer += params
	buffer += "</methodCall>"
	return buffer, err
}

func RPCResponse2XML(rpc interface{}) (string, error) {
	buffer := "<methodResponse>"
	params, err := RPCParams2XML(rpc)
	buffer += params
	buffer += "</methodResponse>"
	return buffer, err
}

func RPCParams2XML(rpc interface{}) (string, error) {
	var err error
	buffer := "<params>"
	for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
		var xml string
		buffer += "<param>"
		xml, err = RPC2XML(reflect.ValueOf(rpc).Elem().Field(i).Interface())
		buffer += xml
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer, err
}

func RPC2XML(value interface{}) (string, error) {
	out := "<value>"
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		out += fmt.Sprintf("<int>%d</int>", value.(int))
	case reflect.Float64:
		out += fmt.Sprintf("<double>%f</double>", value.(float64))
	case reflect.String:
		out += String2XML(value.(string))
	case reflect.Bool:
		out += Bool2XML(value.(bool))
	case reflect.Struct:
		out += Struct2XML(value)
	case reflect.Slice, reflect.Array:
		out += Array2XML(value)
	}
	out += "</value>"
	return out, nil
}

func Bool2XML(value bool) string {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	return fmt.Sprintf("<boolean>%s</boolean>", b)
}

func String2XML(value string) string {
	value = strings.Replace(value, "&", "&amp;", -1)
	value = strings.Replace(value, "\"", "&quot;", -1)
	value = strings.Replace(value, "<", "&lt;", -1)
	value = strings.Replace(value, ">", "&gt;", -1)
	return fmt.Sprintf("<string>%s</string>", value)
}

func Struct2XML(value interface{}) (out string) {
	out += "<struct>"
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		field_type := reflect.TypeOf(value).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		field_value, _ := RPC2XML(field.Interface())
		field_name := fmt.Sprintf("<name>%s</name>", name)
		out += fmt.Sprintf("<member>%s%s</member>", field_name, field_value)
	}
	out += "</struct>"
	return
}

func Array2XML(value interface{}) (out string) {
	out += "<array><data>"
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		item_xml, _ := RPC2XML(reflect.ValueOf(value).Index(i).Interface())
		out += item_xml
	}
	out += "</data></array>"
	return
}

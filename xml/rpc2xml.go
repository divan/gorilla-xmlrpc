// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func rpcRequest2XML(method string, rpc interface{}) (string, error) {
	buffer := "<methodCall><methodName>"
	buffer += method
	buffer += "</methodName>"
	params, err := rpcParams2XML(rpc)
	buffer += params
	buffer += "</methodCall>"
	return buffer, err
}

func rpcResponse2XML(rpc interface{}) (string, error) {
	buffer := "<methodResponse>"
	params, err := rpcParams2XML(rpc)
	buffer += params
	buffer += "</methodResponse>"
	return buffer, err
}

func rpcParams2XML(rpc interface{}) (string, error) {
	var err error
	buffer := "<params>"
	for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
		var xml string
		buffer += "<param>"
		xml, err = rpc2XML(reflect.ValueOf(rpc).Elem().Field(i).Interface())
		buffer += xml
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer, err
}

func rpc2XML(value interface{}) (string, error) {
	out := "<value>"
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		out += fmt.Sprintf("<int>%d</int>", value.(int))
	case reflect.Float64:
		out += fmt.Sprintf("<double>%f</double>", value.(float64))
	case reflect.String:
		out += string2XML(value.(string))
	case reflect.Bool:
		out += bool2XML(value.(bool))
	case reflect.Struct:
		if reflect.TypeOf(value).String() != "time.Time" {
			out += struct2XML(value)
		} else {
			out += time2XML(value.(time.Time))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			out += array2XML(value)
		} else {
			out += base642XML(value.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			out += "<nil/>"
		}
	}
	out += "</value>"
	return out, nil
}

func bool2XML(value bool) string {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	return fmt.Sprintf("<boolean>%s</boolean>", b)
}

func string2XML(value string) string {
	value = strings.Replace(value, "&", "&amp;", -1)
	value = strings.Replace(value, "\"", "&quot;", -1)
	value = strings.Replace(value, "<", "&lt;", -1)
	value = strings.Replace(value, ">", "&gt;", -1)
	return fmt.Sprintf("<string>%s</string>", value)
}

func struct2XML(value interface{}) (out string) {
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
		field_value, _ := rpc2XML(field.Interface())
		field_name := fmt.Sprintf("<name>%s</name>", name)
		out += fmt.Sprintf("<member>%s%s</member>", field_name, field_value)
	}
	out += "</struct>"
	return
}

func array2XML(value interface{}) (out string) {
	out += "<array><data>"
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		item_xml, _ := rpc2XML(reflect.ValueOf(value).Index(i).Interface())
		out += item_xml
	}
	out += "</data></array>"
	return
}

func time2XML(t time.Time) string {
	/*
		// TODO: find out whether we need to deal
		// here with TZ
		var tz string;
		zone, offset := t.Zone()
		if zone == "UTC" {
			tz = "Z"
		} else {
			tz = fmt.Sprintf("%03d00", offset / 3600 )
		}
	*/
	return fmt.Sprintf("<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func base642XML(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("<base64>%s</base64>", str)
}

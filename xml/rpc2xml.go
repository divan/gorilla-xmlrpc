// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build jex
//go:generate jex


package xml

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"

	. "github.com/anjensan/jex"
)

func rpcRequest2XML_(method string, rpc interface{}) string {
	buffer := "<methodCall><methodName>"
	buffer += method
	buffer += "</methodName>"
	buffer += rpcParams2XML_(rpc)
	buffer += "</methodCall>"
	return buffer
}

func rpcResponse2XML_(rpc interface{}) string {
	buffer := "<methodResponse>"
	buffer += rpcParams2XML_(rpc)
	buffer += "</methodResponse>"
	return buffer
}

func rpcParams2XML_(rpc interface{}) string {
	buffer := "<params>"
	for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
		buffer += "<param>"
		buffer += rpc2XML_(reflect.ValueOf(rpc).Elem().Field(i).Interface())
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer
}

func rpc2XML_(value interface{}) string {
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
			out += struct2XML_(value)
		} else {
			out += time2XML(value.(time.Time))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			out += array2XML_(value)
		} else {
			out += base642XML(value.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			out += "<nil/>"
		}
	default:
		THROW(fmt.Errorf("unsupported type %T", value))
	}
	out += "</value>"
	return out
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

func struct2XML_(value interface{}) (out string) {
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
		field_value := rpc2XML_(field.Interface())
		field_name := fmt.Sprintf("<name>%s</name>", name)
		out += fmt.Sprintf("<member>%s%s</member>", field_name, field_value)
	}
	out += "</struct>"
	return
}

func array2XML_(value interface{}) (out string) {
	out += "<array><data>"
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		item_xml := rpc2XML_(reflect.ValueOf(value).Index(i).Interface())
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

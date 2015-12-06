// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

func rpcRequest2XML(out io.Writer, method string, rpc interface{}) error {
	_, err := fmt.Fprintf(out, "<methodCall><methodName>%s</methodName>", method)
	if err = rpcParams2XML(out, rpc); err != nil {
		return err
	}
	_, err = io.WriteString(out, "</methodCall>")
	return err
}

func rpcResponse2XML(out io.Writer, rpc interface{}) error {
	_, err := io.WriteString(out, "<methodResponse>")
	if err = rpcParams2XML(out, rpc); err != nil {
		return err
	}
	_, err = io.WriteString(out, "</methodResponse>")
	return err
}

func rpcParams2XML(out io.Writer, rpc interface{}) error {
	_, err := io.WriteString(out, "<params>")
	if err != nil {
		return err
	}
	if m, ok := rpc.(map[string]interface{}); ok {
		io.WriteString(out, "<struct>")
		for k, v := range m {
			fmt.Fprintf(out, "<param><name>%s</name>", k)
			err = rpc2XML(out, v)
			io.WriteString(out, "</param>")
			if err != nil {
				break
			}
		}
		io.WriteString(out, "</struct>")
	} else {
		for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
			io.WriteString(out, "<param>")
			err = rpc2XML(out, reflect.ValueOf(rpc).Elem().Field(i).Interface())
			io.WriteString(out, "</param>")
			if err != nil {
				break
			}
		}
	}
	_, _ = io.WriteString(out, "</params>")
	return err
}

func rpc2XML(out io.Writer, value interface{}) error {
	_, err := io.WriteString(out, "<value>")
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		_, err = fmt.Fprintf(out, "<int>%d</int>", value.(int))
	case reflect.Float64:
		_, err = fmt.Fprintf(out, "<double>%f</double>", value.(float64))
	case reflect.String:
		err = string2XML(out, value.(string))
	case reflect.Bool:
		err = bool2XML(out, value.(bool))
	case reflect.Struct:
		if reflect.TypeOf(value).String() != "time.Time" {
			err = struct2XML(out, value)
		} else {
			err = time2XML(out, value.(time.Time))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			err = array2XML(out, value)
		} else {
			err = base642XML(out, value.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			_, err = io.WriteString(out, "<nil/>")
		}
	}
	_, err = io.WriteString(out, "</value>")
	return err
}

func bool2XML(out io.Writer, value bool) error {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	_, err := fmt.Fprintf(out, "<boolean>%s</boolean>", b)
	return err
}

var strRepl = strings.NewReplacer(
	"&", "&amp;",
	`"`, "&quot;",
	"<", "&lt;",
	">", "&gt;",
)

func string2XML(out io.Writer, value string) error {
	_, err := fmt.Fprintf(out, "<string>%s</string>", strRepl.Replace(value))
	return err
}

func struct2XML(out io.Writer, value interface{}) error {
	_, err := io.WriteString(out, "<struct>")
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		field_type := reflect.TypeOf(value).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		_, err = fmt.Fprintf(out, "<member><name>%s</name>", name)
		err = rpc2XML(out, field.Interface())
		_, err = io.WriteString(out, "</member>")
	}
	_, err = io.WriteString(out, "</struct>")
	return err
}

func array2XML(out io.Writer, value interface{}) error {
	_, err := io.WriteString(out, "<array><data>")
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		err = rpc2XML(out, reflect.ValueOf(value).Index(i).Interface())
	}
	_, err = io.WriteString(out, "</data></array>")
	return err
}

func time2XML(out io.Writer, t time.Time) error {
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
	_, err := fmt.Fprintf(out, "<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return err
}

func base642XML(out io.Writer, data []byte) error {
	_, _ = io.WriteString(out, "<base64>")
	w := base64.NewEncoder(base64.StdEncoding, out)
	_, _ = w.Write(data)
	err := w.Close()
	_, _ = io.WriteString(out, "</base64>")
	return err
}

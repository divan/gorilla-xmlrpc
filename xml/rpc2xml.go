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

func rpcRequest2XML(w io.Writer, method string, rpc interface{}) error {
	_, err := fmt.Fprintf(w, "<methodCall><methodName>%s</methodName>", method)
	if err = rpcParams2XML(w, rpc); err != nil {
		return err
	}
	_, err = io.WriteString(w, "</methodCall>")
	return err
}

func rpcResponse2XML(w io.Writer, rpc interface{}) error {
	_, err := io.WriteString(w, "<methodResponse>")
	if err = rpcParams2XML(w, rpc); err != nil {
		return err
	}
	_, err = io.WriteString(w, "</methodResponse>")
	return err
}

func rpcParams2XML(w io.Writer, rpc interface{}) error {
	_, err := io.WriteString(w, "<params>")
	if err != nil {
		return err
	}
	if m, ok := rpc.(map[string]interface{}); ok {
		io.WriteString(w, "<struct>")
		for k, v := range m {
			fmt.Fprintf(w, "<param><name>%s</name>", k)
			err = rpc2XML(w, v)
			io.WriteString(w, "</param>")
			if err != nil {
				break
			}
		}
		io.WriteString(w, "</struct>")
	} else {
		for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
			io.WriteString(w, "<param>")
			err = rpc2XML(w, reflect.ValueOf(rpc).Elem().Field(i).Interface())
			io.WriteString(w, "</param>")
			if err != nil {
				break
			}
		}
	}
	_, _ = io.WriteString(w, "</params>")
	return err
}

func rpc2XML(w io.Writer, value interface{}) error {
	_, err := io.WriteString(w, "<value>")
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		_, err = fmt.Fprintf(w, "<int>%d</int>", value.(int))
	case reflect.Float64:
		_, err = fmt.Fprintf(w, "<double>%f</double>", value.(float64))
	case reflect.String:
		err = string2XML(w, value.(string))
	case reflect.Bool:
		err = bool2XML(w, value.(bool))
	case reflect.Struct:
		if reflect.TypeOf(value).String() != "time.Time" {
			err = struct2XML(w, value)
		} else {
			err = time2XML(w, value.(time.Time))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			err = array2XML(w, value)
		} else {
			err = base642XML(w, value.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			_, err = io.WriteString(w, "<nil/>")
		}
	}
	_, err = io.WriteString(w, "</value>")
	return err
}

func bool2XML(w io.Writer, value bool) error {
	var b string
	if value {
		b = "1"
	} else {
		b = "0"
	}
	_, err := fmt.Fprintf(w, "<boolean>%s</boolean>", b)
	return err
}

var strRepl = strings.NewReplacer(
	"&", "&amp;",
	`"`, "&quot;",
	"<", "&lt;",
	">", "&gt;",
)

func string2XML(w io.Writer, value string) error {
	_, err := fmt.Fprintf(w, "<string>%s</string>", strRepl.Replace(value))
	return err
}

func struct2XML(w io.Writer, value interface{}) error {
	_, err := io.WriteString(w, "<struct>")
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		field_type := reflect.TypeOf(value).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		_, err = fmt.Fprintf(w, "<member><name>%s</name>", name)
		err = rpc2XML(w, field.Interface())
		_, err = io.WriteString(w, "</member>")
	}
	_, err = io.WriteString(w, "</struct>")
	return err
}

func array2XML(w io.Writer, value interface{}) error {
	_, err := io.WriteString(w, "<array><data>")
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		err = rpc2XML(w, reflect.ValueOf(value).Index(i).Interface())
	}
	_, err = io.WriteString(w, "</data></array>")
	return err
}

func time2XML(w io.Writer, t time.Time) error {
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
	_, err := fmt.Fprintf(w, "<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return err
}

func base642XML(w io.Writer, data []byte) error {
	_, _ = io.WriteString(w, "<base64>")
	bw := base64.NewEncoder(base64.StdEncoding, w)
	_, _ = bw.Write(data)
	err := bw.Close()
	_, _ = io.WriteString(w, "</base64>")
	return err
}

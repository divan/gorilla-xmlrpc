// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

func rpcRequest2XML(w io.Writer, method string, rpc interface{}) error {
	ew := NewErrWriter(w)
	fmt.Fprintf(ew, "<methodCall><methodName>%s</methodName>", method)
	err := rpcParams2XML(ew, rpc)
	io.WriteString(ew, "</methodCall>")
	if err != nil {
		return err
	}
	return ew.Err()
}

func rpcResponse2XML(w io.Writer, rpc interface{}) error {
	ew := NewErrWriter(w)
	io.WriteString(ew, "<methodResponse>")
	err := rpcParams2XML(ew, rpc)
	io.WriteString(ew, "</methodResponse>")
	if err != nil {
		return err
	}
	return ew.Err()
}

func rpcParams2XML(w io.Writer, rpc interface{}) error {
	ew := NewErrWriter(w)
	io.WriteString(ew, "<params>")
	var err error
	if m, ok := rpc.(map[string]interface{}); ok {
		io.WriteString(w, "<struct>")
		for k, v := range m {
			fmt.Fprintf(w, "<param><name>")
			xml.EscapeText(w, []byte(k))
			fmt.Fprintf(w, "</name>")
			err = rpc2XML(w, v)
			io.WriteString(w, "</param>")
			if err != nil {
				break
			}
		}
		io.WriteString(w, "</struct>")

	} else {
		for i := 0; i < reflect.ValueOf(rpc).Elem().NumField(); i++ {
			io.WriteString(ew, "<param>")
			err = rpc2XML(ew, reflect.ValueOf(rpc).Elem().Field(i).Interface())
			io.WriteString(ew, "</param>")
			if err != nil {
				break
			}
		}
	}
	io.WriteString(ew, "</params>")
	if err != nil {
		return err
	}
	return ew.Err()
}

func rpc2XML(w io.Writer, value interface{}) error {
	ew := NewErrWriter(w)
	io.WriteString(ew, "<value>")
	var err error
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		fmt.Fprintf(ew, "<int>%d</int>", value.(int))
	case reflect.Float64:
		fmt.Fprintf(ew, "<double>%f</double>", value.(float64))
	case reflect.String:
		err = string2XML(ew, value.(string))
	case reflect.Bool:
		err = bool2XML(ew, value.(bool))
	case reflect.Struct:
		if reflect.TypeOf(value).String() != "time.Time" {
			err = struct2XML(ew, value)
		} else {
			err = time2XML(ew, value.(time.Time))
		}
	case reflect.Map:
		fmt.Fprintf(ew, "<struct>")
		for k, v := range value.(map[string]interface{}) {
			fmt.Fprintf(ew, "<member><name>")
			xml.EscapeText(ew, []byte(k))
			fmt.Fprintf(ew, "</name><value>")
			err = rpc2XML(ew, v)
			fmt.Fprintf(ew, "</member>")
			if err != nil {
				break
			}
		}
		fmt.Fprintf(ew, "</struct>")
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(value).String() != "[]uint8" {
			err = array2XML(ew, value)
		} else {
			err = base642XML(ew, value.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(value).IsNil() {
			io.WriteString(ew, "<nil/>")
		}
	}
	io.WriteString(ew, "</value>")
	if err != nil {
		return err
	}
	return ew.Err()
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
	ew := NewErrWriter(w)
	io.WriteString(ew, "<struct>")
	var err error
	for i := 0; i < reflect.TypeOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		field_type := reflect.TypeOf(value).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		fmt.Fprintf(ew, "<member><name>%s</name>", name)
		err = rpc2XML(ew, field.Interface())
		io.WriteString(ew, "</member>")
		if err != nil {
			break
		}
	}
	io.WriteString(ew, "</struct>")
	if err != nil {
		return err
	}
	return ew.Err()
}

func array2XML(w io.Writer, value interface{}) error {
	ew := NewErrWriter(w)
	io.WriteString(ew, "<array><data>")
	var err error
	for i := 0; i < reflect.ValueOf(value).Len(); i++ {
		if err = rpc2XML(ew, reflect.ValueOf(value).Index(i).Interface()); err != nil {
			break
		}
	}
	io.WriteString(ew, "</data></array>")
	if err != nil {
		return err
	}
	return ew.Err()
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
	ew := NewErrWriter(w)
	_, _ = io.WriteString(ew, "<base64>")
	bw := base64.NewEncoder(base64.StdEncoding, ew)
	bw.Write(data)
	err := bw.Close()
	io.WriteString(ew, "</base64>")
	if err != nil {
		return err
	}
	return ew.Err()
}

var _ = io.Writer((*errWriter)(nil))

type errWriter struct {
	w   io.Writer
	err error
}

func NewErrWriter(w io.Writer) *errWriter {
	if w == nil {
		return nil
	}
	if ew, ok := w.(*errWriter); ok {
		return ew
	}
	return &errWriter{w: w}
}

func (ew *errWriter) Write(p []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	var n int
	n, ew.err = ew.w.Write(p)
	return n, ew.err
}

func (ew *errWriter) Err() error {
	return ew.err
}

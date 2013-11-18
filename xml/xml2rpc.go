// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Types used for unmarshalling
type Response struct {
	Name   xml.Name `xml:"methodResponse"`
	Params []Param  `xml:"params>param"`
}

type Param struct {
	Value Value `xml:"value"`
}

type Value struct {
	Array    []Value  `xml:"array>data>value"`
	Struct   []Member `xml:"struct>member"`
	String   string   `xml:"string"`
	Int      string   `xml:"int"`
	Int4     string   `xml:"i4"`
	Double   string   `xml:"double"`
	Boolean  string   `xml:"boolean"`
	DateTime string   `xml:"dateTime.iso8601"`
	Base64   string   `xml:"base64"`
}

type Member struct {
	Name  string `xml:"name"`
	Value Value  `xml:"value"`
}

func XML2RPC(xmlraw string, rpc interface{}) (err error) {
	// Unmarshal raw XML into the temporal structure
	var ret Response
	err = xml.Unmarshal([]byte(xmlraw), &ret)
	if err != nil {
		return
	}

	// Structures should have equal number of fields
	if reflect.TypeOf(rpc).Elem().NumField() != len(ret.Params) {
		return errors.New("Wrong number of arguments")
	}

	// Now, convert temporal structure into the
	// passed rpc variable, according to it's structure
	for i, param := range ret.Params {
		field := reflect.ValueOf(rpc).Elem().Field(i)
		err = Value2Field(param.Value, &field)
		if err != nil {
			return
		}
	}

	return
}

func Value2Field(value Value, field *reflect.Value) (err error) {
	var val interface{}
	switch {
	case value.Int != "":
		val, _ = strconv.Atoi(value.Int)
	case value.Int4 != "":
		val, _ = strconv.Atoi(value.Int4)
	case value.Double != "":
		val, _ = strconv.ParseFloat(value.Double, 64)
	case value.String != "":
		val = value.String
	case value.Boolean != "":
		val = XML2Bool(value.Boolean)
	case value.DateTime != "":
		val, err = XML2DateTime(value.DateTime)
	case len(value.Struct) != 0:
		s := value.Struct
		for i := 0; i < len(s); i++ {
			f := field.FieldByName(s[i].Name)
			err = Value2Field(s[i].Value, &f)
		}
	case len(value.Array) != 0:
		a := value.Array
		f := *field
		slice := reflect.MakeSlice(reflect.TypeOf(f.Interface()),
			len(a), len(a))
		for i := 0; i < len(a); i++ {
			item := slice.Index(i)
			err = Value2Field(a[i], &item)
		}
		f = reflect.AppendSlice(f, slice)
		val = f.Interface()
	}

	if val != nil {
		if reflect.TypeOf(val) != reflect.TypeOf(field.Interface()) {
			return errors.New(fmt.Sprintf("Fields type mismatch: %s != %s",
				reflect.TypeOf(val),
				reflect.TypeOf(field.Interface())))
		}

		field.Set(reflect.ValueOf(val))
	}
	return
}

func XML2Bool(value string) bool {
	var b bool
	switch value {
	case "1", "true", "TRUE", "True":
		b = true
	case "0", "false", "FALSE", "False":
		b = false
	}
	return b
}

func XML2DateTime(value string) (time.Time, error) {
	var (
		year, month, day     int
		hour, minute, second int
	)
	_, err := fmt.Sscanf(value, "%04d%02d%02dT%02d:%02d:%02d",
		&year, &month, &day,
		&hour, &minute, &second)
	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return t, err
}

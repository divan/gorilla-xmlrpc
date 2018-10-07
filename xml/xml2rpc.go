// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build jex
//go:generate jex

package xml

import . "github.com/anjensan/jex"

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/rogpeppe/go-charset/charset"
	_ "github.com/rogpeppe/go-charset/data"

	"github.com/anjensan/jex/ex"
)

// Types used for unmarshalling
type response struct {
	Name   xml.Name   `xml:"methodResponse"`
	Params []param    `xml:"params>param"`
	Fault  faultValue `xml:"fault,omitempty"`
}

type param struct {
	Value value `xml:"value"`
}

type value struct {
	Array    []value  `xml:"array>data>value"`
	Struct   []member `xml:"struct>member"`
	String   string   `xml:"string"`
	Int      string   `xml:"int"`
	Int4     string   `xml:"i4"`
	Double   string   `xml:"double"`
	Boolean  string   `xml:"boolean"`
	DateTime string   `xml:"dateTime.iso8601"`
	Base64   string   `xml:"base64"`
	Raw      string   `xml:",innerxml"` // the value can be defualt string
}

type member struct {
	Name  string `xml:"name"`
	Value value  `xml:"value"`
}

func xml2RPC_(xmlraw string, rpc interface{}) {
	// Unmarshal raw XML into the temporal structure
	var ret response
	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlraw)))
	decoder.CharsetReader = charset.NewReader

	if TRY() {
		ERR = decoder.Decode(&ret)
	} else {
		THROW(FaultDecode)
	}

	if !ret.Fault.IsEmpty() {
		THROW(getFaultResponse(ret.Fault))
	}

	// Structures should have equal number of fields
	if reflect.TypeOf(rpc).Elem().NumField() != len(ret.Params) {
		THROW(FaultWrongArgumentsNumber)
	}

	// Now, convert temporal structure into the
	// passed rpc variable, according to it's structure
	for i, param := range ret.Params {
		field := reflect.ValueOf(rpc).Elem().Field(i)
		value2Field_(param.Value, &field)
	}
}

// getFaultResponse converts faultValue to Fault.
func getFaultResponse(fault faultValue) Fault {
	var (
		code int
		str  string
	)

	for _, field := range fault.Value.Struct {
		if field.Name == "faultCode" {
			code, _ = strconv.Atoi(field.Value.Int)
		} else if field.Name == "faultString" {
			str = field.Value.String
			if str == "" {
				str = field.Value.Raw
			}
		}
	}

	return Fault{Code: code, String: str}
}

func value2Field_(value value, field *reflect.Value) {
	ex.Check_(field.CanSet(), FaultApplicationError)
	var val interface{}

	switch {
	case value.Int != "":
		val, ERR = strconv.Atoi(value.Int)
	case value.Int4 != "":
		val, ERR = strconv.Atoi(value.Int4)
	case value.Double != "":
		val, ERR = strconv.ParseFloat(value.Double, 64)
	case value.String != "":
		val = value.String
	case value.Boolean != "":
		val = xml2Bool(value.Boolean)
	case value.DateTime != "":
		val = xml2DateTime_(value.DateTime)
	case value.Base64 != "":
		val = xml2Base64_(value.Base64)
	case len(value.Struct) != 0:
		if field.Kind() != reflect.Struct {
			fault := FaultInvalidParams
			fault.String += fmt.Sprintf("structure fields mismatch: %s != %s", field.Kind(), reflect.Struct.String())
			THROW(fault)
		}
		s := value.Struct
		for i := 0; i < len(s); i++ {
			// Uppercase first letter for field name to deal with
			// methods in lowercase, which cannot be used
			field_name := uppercaseFirst(s[i].Name)
			f := field.FieldByName(field_name)
			value2Field_(s[i].Value, &f)
		}
	case len(value.Array) != 0:
		a := value.Array
		f := *field
		slice := reflect.MakeSlice(reflect.TypeOf(f.Interface()),
			len(a), len(a))
		for i := 0; i < len(a); i++ {
			item := slice.Index(i)
			value2Field_(a[i], &item)
		}
		f = reflect.AppendSlice(f, slice)
		val = f.Interface()

	default:
		// value field is default to string, see http://en.wikipedia.org/wiki/XML-RPC#Data_types
		// also can be <nil/>
		if value.Raw != "<nil/>" {
			val = value.Raw
		}
	}

	if val != nil {
		if reflect.TypeOf(val) != reflect.TypeOf(field.Interface()) {
			fault := FaultInvalidParams
			fault.String += fmt.Sprintf(": fields type mismatch: %s != %s",
				reflect.TypeOf(val),
				reflect.TypeOf(field.Interface()))
			THROW(fault)
		}

		field.Set(reflect.ValueOf(val))
	}
}

func xml2Bool(value string) bool {
	var b bool
	switch value {
	case "1", "true", "TRUE", "True":
		b = true
	case "0", "false", "FALSE", "False":
		b = false
	}
	return b
}

func xml2DateTime_(value string) time.Time {
	var (
		year, month, day     int
		hour, minute, second int
	)
	_, ERR := fmt.Sscanf(value, "%04d%02d%02dT%02d:%02d:%02d",
		&year, &month, &day,
		&hour, &minute, &second)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func xml2Base64_(value string) []byte {
	r, ERR := base64.StdEncoding.DecodeString(value)
	return r
}

func uppercaseFirst(in string) (out string) {
	r, n := utf8.DecodeRuneInString(in)
	return string(unicode.ToUpper(r)) + in[n:]
}

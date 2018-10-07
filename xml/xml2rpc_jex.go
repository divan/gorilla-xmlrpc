//line xml2rpc.go:1
// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !jex
//jex:off

package xml

//line xml2rpc.go:8
import _jex "github.com/anjensan/jex/runtime"

//line xml2rpc.go:12
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
	Name	xml.Name		`xml:"methodResponse"`
	Params	[]param		`xml:"params>param"`
	Fault	faultValue		`xml:"fault,omitempty"`
}

type param struct {
	Value value `xml:"value"`
}

type value struct {
	Array	[]value		`xml:"array>data>value"`
	Struct	[]member		`xml:"struct>member"`
	String	string		`xml:"string"`
	Int	string		`xml:"int"`
	Int4	string		`xml:"i4"`
	Double	string		`xml:"double"`
	Boolean	string		`xml:"boolean"`
	DateTime	string		`xml:"dateTime.iso8601"`
	Base64	string		`xml:"base64"`
	Raw	string		`xml:",innerxml"`	// the value can be defualt string
}

type member struct {
	Name	string		`xml:"name"`
	Value	value		`xml:"value"`
}

func xml2RPC_(xmlraw string, rpc interface{}) {
	// Unmarshal raw XML into the temporal structure
	var ret response
	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlraw)))
	decoder.CharsetReader = charset.NewReader
//line xml2rpc.go:62
	_jex.TryCatch(func() {
//line xml2rpc.go:64
		var _jex_e1422 error
//line xml2rpc.go:64
		_jex_e1422 = decoder.Decode(&ret)
		_jex.Must(_jex_e1422)
	}, func(_jex_ex _jex.Exception) {
//line xml2rpc.go:66
		defer _jex.Suppress(_jex_ex)
//line xml2rpc.go:66
		panic(_jex.NewException(FaultDecode))
//line xml2rpc.go:68
	})

	if !ret.Fault.IsEmpty() {
//line xml2rpc.go:70
		panic(_jex.NewException(getFaultResponse(ret.Fault)))
//line xml2rpc.go:72
	}

	// Structures should have equal number of fields
	if reflect.TypeOf(rpc).Elem().NumField() != len(ret.Params) {
//line xml2rpc.go:75
		panic(_jex.NewException(FaultWrongArgumentsNumber))
//line xml2rpc.go:77
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
		code	int
		str	string
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
//line xml2rpc.go:113
		var _jex_e2513 error
		val, _jex_e2513 = strconv.Atoi(value.Int)
//line xml2rpc.go:114
		_jex.Must(_jex_e2513)
	case value.Int4 != "":
//line xml2rpc.go:115
		var _jex_e2574 error
		val, _jex_e2574 = strconv.Atoi(value.Int4)
//line xml2rpc.go:116
		_jex.Must(_jex_e2574)
	case value.Double != "":
//line xml2rpc.go:117
		var _jex_e2638 error
		val, _jex_e2638 = strconv.ParseFloat(value.Double, 64)
//line xml2rpc.go:118
		_jex.Must(_jex_e2638)
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
//line xml2rpc.go:130
			panic(_jex.NewException(fault))
//line xml2rpc.go:132
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
//line xml2rpc.go:166
			panic(_jex.NewException(fault))
//line xml2rpc.go:168
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
		year, month, day	int
		hour, minute, second	int
	)
	_, _jex_e4482 := fmt.Sscanf(value, "%04d%02d%02dT%02d:%02d:%02d",
		&year, &month, &day,
		&hour, &minute, &second)
//line xml2rpc.go:192
	_jex.Must(_jex_e4482)
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func xml2Base64_(value string) []byte {
	r, _jex_e4720 := base64.StdEncoding.DecodeString(value)
//line xml2rpc.go:197
	_jex.Must(_jex_e4720)
	return r
}

func uppercaseFirst(in string) (out string) {
	r, n := utf8.DecodeRuneInString(in)
	return string(unicode.ToUpper(r)) + in[n:]
}

//line xml2rpc.go:204
const _ = _jex.Unused

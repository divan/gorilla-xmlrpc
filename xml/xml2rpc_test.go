// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"reflect"
	"testing"
	"time"
)

type SubStructXml2Rpc struct {
	Foo  int
	Bar  string
	Data []int
}

type StructXml2Rpc struct {
	Int   int
	Float float64
	Str   string
	Bool  bool
	Sub   SubStructXml2Rpc
	Time  time.Time
}

func TestXML2RPC(t *testing.T) {
	req := new(StructXml2Rpc)
	err := XML2RPC("<methodCall><methodName>Some.Method</methodName><params><param><value><i4>123</i4></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><dateTime.iso8601>20120717T14:08:55</dateTime.iso8601></value></param></params></methodCall>", req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	expected_req := &StructXml2Rpc{123, 3.145926, "Hello, World!", false, SubStructXml2Rpc{42, "I'm Bar", []int{1, 2, 3}}, time.Date(2012, time.July, 17, 14, 8, 55, 0, time.Local)}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

type StructSpecialCharsXml2Rpc struct {
	String1 string
}

func TestXML2RPCSpecialChars(t *testing.T) {
	req := new(StructSpecialCharsXml2Rpc)
	err := XML2RPC("<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>", req)
	if err != nil {
		t.Error("XML2RPC conversion failed", err)
	}
	expected_req := &StructSpecialCharsXml2Rpc{" & \" < > "}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import "testing"

type SubStructRpc2Xml struct {
	Foo int
	Bar string
	Data []int
}

type StructRpc2Xml struct {
	Int   int
	Float float64
	Str   string
	Bool  bool
	Sub   SubStructRpc2Xml
}

func TestRPC2XML(t *testing.T) {
	req := &StructRpc2Xml{123, 3.145926, "Hello, World!", false, SubStructRpc2Xml{42, "I'm Bar", []int{1,2,3}}}
	xml, err := RPCRequest2XML("Some.Method", req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodCall><methodName>Some.Method</methodName><params><param><value><int>123</int></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param></params></methodCall>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

type StructSpecialCharsRpc2Xml struct {
	String1 string
}

func TestRPC2XMLSpecialChars(t *testing.T) {
	req := &StructSpecialCharsRpc2Xml{" & \" < > "}
	xml, err := RPCResponse2XML(req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}
	expected := "<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML Special chars conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

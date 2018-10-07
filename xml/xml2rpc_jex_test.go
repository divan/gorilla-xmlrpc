//line xml2rpc_test.go:1
// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !jex
//jex:off

package xml

//line xml2rpc_test.go:8
import _jex "github.com/anjensan/jex/runtime"

//line xml2rpc_test.go:12
import (
	"reflect"
	"testing"
	"time"

	"github.com/anjensan/jex/ex"
)

type SubStructXml2Rpc struct {
	Foo	int
	Bar	string
	Data	[]int
}

type StructXml2Rpc struct {
	Int	int
	Float	float64
	Str	string
	Bool	bool
	Sub	SubStructXml2Rpc
	Time	time.Time
	Base64	[]byte
}

func TestXML2RPC_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructXml2Rpc)
	xml2RPC_("<methodCall><methodName>Some.Method</methodName><params><param><value><i4>123</i4></value></param><param><value><double>3.145926</double></value></param><param><value><string>Hello, World!</string></value></param><param><value><boolean>0</boolean></value></param><param><value><struct><member><name>Foo</name><value><int>42</int></value></member><member><name>Bar</name><value><string>I'm Bar</string></value></member><member><name>Data</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param><param><value><dateTime.iso8601>20120717T14:08:55</dateTime.iso8601></value></param><param><value><base64>eW91IGNhbid0IHJlYWQgdGhpcyE=</base64></value></param></params></methodCall>", req)
	expected_req := &StructXml2Rpc{123, 3.145926, "Hello, World!", false, SubStructXml2Rpc{42, "I'm Bar", []int{1, 2, 3}}, time.Date(2012, time.July, 17, 14, 8, 55, 0, time.Local), []byte("you can't read this!")}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

type StructSpecialCharsXml2Rpc struct {
	String1 string
}

func TestXML2RPCSpecialChars_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructSpecialCharsXml2Rpc)
	xml2RPC_("<methodResponse><params><param><value><string> &amp; &quot; &lt; &gt; </string></value></param></params></methodResponse>", req)
	expected_req := &StructSpecialCharsXml2Rpc{" & \" < > "}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

type StructNilXml2Rpc struct {
	Ptr *int
}

func TestXML2RPCNil_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructNilXml2Rpc)
	xml2RPC_("<methodResponse><params><param><value><nil/></value></param></params></methodResponse>", req)
	expected_req := &StructNilXml2Rpc{nil}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

type StructXml2RpcSubArgs struct {
	String1	string
	String2	string
	Id	int
}

type StructXml2RpcHelloArgs struct {
	Args StructXml2RpcSubArgs
}

func TestXML2RPCLowercasedMethods_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructXml2RpcHelloArgs)
	xml2RPC_("<methodCall><params><param><value><struct><member><name>string1</name><value><string>I'm a first string</string></value></member><member><name>string2</name><value><string>I'm a second string</string></value></member><member><name>id</name><value><int>1</int></value></member></struct></value></param></params></methodCall>", req)
	args := StructXml2RpcSubArgs{"I'm a first string", "I'm a second string", 1}
	expected_req := &StructXml2RpcHelloArgs{args}
	if !reflect.DeepEqual(req, expected_req) {
		t.Error("XML2RPC conversion failed")
		t.Error("Expected", expected_req)
		t.Error("Got", req)
	}
}

func TestXML2PRCFaultCall_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructXml2RpcHelloArgs)
	data := `<?xmlversion="1.0"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>116</int></value></member><member><name>faultString</name><value><string>Error
Requiredattribute'user'notfound:
[{'User',"gggg"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
</string></value></member></struct></value></fault></methodResponse>`

	errstr := `Error
Requiredattribute'user'notfound:
[{'User',"gggg"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
`
//line xml2rpc_test.go:114
	_jex.TryCatch(func() {
//line xml2rpc_test.go:116
		xml2RPC_(data, req)
		t.Fatal("expected err to be not nil")
	}, func(_jex_ex _jex.Exception) {
//line xml2rpc_test.go:118
		defer _jex.Suppress(_jex_ex)
		fault, ok := _jex_ex.Err().(Fault)
		if !ok {
			t.Errorf("error should be of concrete type Fault, but got %v", _jex_ex.Err())
		} else {
			if fault.Code != 116 {
				t.Errorf("expected fault.Code to be %d, but got %d", 116, fault.Code)
			}
			if fault.String != errstr {
				t.Errorf("fault.String should be:\n\n%s\n\nbut got:\n\n%s\n", errstr, fault.String)
			}
		}
	})
}

func TestXML2PRCISO88591_(t *testing.T) {
	defer ex.Catch(errorReporter(t))
	req := new(StructXml2RpcHelloArgs)
	data := `<?xml version="1.0" encoding="ISO-8859-1"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>116</int></value></member><member><name>faultString</name><value><string>Error
Requiredattribute'user'notfound:
[{'User',"` + "\xd6\xf1\xe4" + `"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
</string></value></member></struct></value></fault></methodResponse>`

	errstr := `Error
Requiredattribute'user'notfound:
[{'User',"Öñä"},{'Host',"sss.com"},{'Password',"ssddfsdf"}]
`
//line xml2rpc_test.go:144
	_jex.TryCatch(func() {
//line xml2rpc_test.go:146
		xml2RPC_(data, req)
		t.Fatal("expected err to be not nil")
	}, func(_jex_ex _jex.Exception) {
//line xml2rpc_test.go:148
		defer _jex.Suppress(_jex_ex)
		fault, ok := _jex_ex.Err().(Fault)
		if !ok {
			t.Errorf("error should be of concrete type Fault, but got %v", _jex_ex.Err())
		} else {
			if fault.Code != 116 {
				t.Errorf("expected fault.Code to be %d, but got %d", 116, fault.Code)
			}
			if fault.String != errstr {
				t.Errorf("fault.String should be:\n\n%s\n\nbut got:\n\n%s\n", errstr, fault.String)
			}
		}
	})
}

//line xml2rpc_test.go:161
const _ = _jex.Unused

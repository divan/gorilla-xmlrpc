// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/rpc"
)

//////////////////////////////////
// Service 1
//////////////////////////////////
type Service1Request struct {
	A int
	B int
}

type Service1BadRequest struct {
	A int
	B int
	C int
}

type Service1Response struct {
	Result int
}

type Service1 struct {
}

func (t *Service1) Multiply(r *http.Request, req *Service1Request, res *Service1Response) error {
	res.Result = req.A * req.B
	return nil
}

//////////////////////////////////
// Service 2
//////////////////////////////////
type Service2Request struct {
	Name      string
	Age       int
	HasPermit bool
}

type Service2Response struct {
	Message string
	Status  int
}

type Service2 struct {
}

func (t *Service2) GetGreeting(r *http.Request, req *Service2Request, res *Service2Response) error {
	res.Message = "Hello, user " + req.Name + ". You're " + strconv.Itoa(req.Age) + " years old :-P"
	if req.HasPermit {
		res.Message += " And you has permit."
	} else {
		res.Message += " And you DON'T has permit."
	}
	res.Status = 42
	return nil
}

//////////////////////////////////
// Service 3
//////////////////////////////////

type Address struct {
	Number  int
	Street  string
	Country string
}

type Person struct {
	Name    string
	Surname string
	Age     int
	Address Address
}

type Info struct {
	Facebook string
	Twitter  string
	Phone    string
}

type Service3Request struct {
	Person Person
}

type Service3Response struct {
	Info Info
}

type Service3 struct {
}

func (t *Service3) GetInfo(r *http.Request, req *Service3Request, res *Service3Response) error {
	var i Info
	i.Facebook = "http://facebook.com/" + req.Person.Name
	i.Twitter = "http://twitter.com/" + req.Person.Name
	i.Phone = "+55-555-555-55-55"
	res.Info = i
	return nil
}

func execute(t *testing.T, s *rpc.Server, method string, req, res interface{}) error {
	if !s.HasMethod(method) {
		t.Fatal("Expected to be registered:", method)
	}

	buf, _ := EncodeClientRequest(method, req)
	body := bytes.NewBuffer(buf)
	r, _ := http.NewRequest("POST", "http://localhost:8080/", body)
	r.Header.Set("Content-Type", "text/xml")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	return DecodeClientResponse(w.Body, res)
}

func TestRPC2XMLConverter(t *testing.T) {
	req := &Service1Request{4, 2}
	xml, err := rpcRequest2XML("Some.Method", req)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}

	expected := "<methodCall><methodName>Some.Method</methodName><params><param><value><int>4</int></value></param><param><value><int>2</int></value></param></params></methodCall>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}

	req2 := &Service2Request{"Johnny", 33, true}
	xml, err = rpcRequest2XML("Some.Method", req2)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}

	expected = "<methodCall><methodName>Some.Method</methodName><params><param><value><string>Johnny</string></value></param><param><value><int>33</int></value></param><param><value><boolean>1</boolean></value></param></params></methodCall>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}

	address := Address{221, "Baker str.", "London"}
	person := Person{"Johnny", "Doe", 33, address}
	req3 := &Service3Request{person}
	xml, err = rpcRequest2XML("Some.Method", req3)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}

	expected = "<methodCall><methodName>Some.Method</methodName><params><param><value><struct><member><name>Name</name><value><string>Johnny</string></value></member><member><name>Surname</name><value><string>Doe</string></value></member><member><name>Age</name><value><int>33</int></value></member><member><name>Address</name><value><struct><member><name>Number</name><value><int>221</int></value></member><member><name>Street</name><value><string>Baker str.</string></value></member><member><name>Country</name><value><string>London</string></value></member></struct></value></member></struct></value></param></params></methodCall>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}

	res := &Service1Response{42}
	xml, err = rpcResponse2XML(res)
	if err != nil {
		t.Error("RPC2XML conversion failed", err)
	}

	expected = "<methodResponse><params><param><value><int>42</int></value></param></params></methodResponse>"
	if xml != expected {
		t.Error("RPC2XML conversion failed")
		t.Error("Expected", expected)
		t.Error("Got", xml)
	}
}

func TestServices(t *testing.T) {
	s := rpc.NewServer()
	s.RegisterCodec(NewCodec(), "text/xml")
	s.RegisterService(new(Service1), "")
	s.RegisterService(new(Service2), "")
	s.RegisterService(new(Service3), "")

	var res Service1Response
	if err := execute(t, s, "Service1.Multiply", &Service1Request{4, 2}, &res); err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}
	if res.Result != 8 {
		t.Errorf("Wrong response: %v.", res.Result)
	}

	var res2 Service2Response
	if err := execute(t, s, "Service2.GetGreeting", &Service2Request{"Johnny", 33, true}, &res2); err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}
	if res2.Message != "Hello, user Johnny. You're 33 years old :-P And you has permit." {
		t.Errorf("Wrong response: %v.", res2.Message)
	}
	if res2.Status != 42 {
		t.Errorf("Wrong response: %v.", res2.Status)
	}

	var res3 Service3Response
	address := Address{221, "Baker str.", "London"}
	person := Person{"Johnny", "Doe", 33, address}
	if err := execute(t, s, "Service3.GetInfo", &Service3Request{person}, &res3); err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}

	if res3.Info.Phone != "+55-555-555-55-55" {
		t.Errorf("Wrong response: %v.", res3.Info)
	}
}

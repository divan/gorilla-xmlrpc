// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/rpc"
)

//////////////////////////////////
// Service 1
//////////////////////////////////
type FaultTestRequest struct {
	A int
	B int
}

type FaultTestBadRequest struct {
	A int
	B int
	C int
}

type FaultTestResponse struct {
	Result int
}

type FaultTestBadResponse struct {
	Result string
}

type FaultTest struct {
}

func (t *FaultTest) Multiply(r *http.Request, req *FaultTestRequest, res *FaultTestResponse) error {
	res.Result = req.A * req.B
	return nil
}

func TestFaults(t *testing.T) {
	s := rpc.NewServer()
	s.RegisterCodec(NewCodec(), "text/xml")
	s.RegisterService(new(FaultTest), "")

	var err error

	var res1 FaultTestResponse
	err = execute(t, s, "FaultTest.Multiply", &FaultTestBadRequest{4, 2, 4}, &res1)
	if err == nil {
		t.Fatal("expected err to be not nil, but got:", err)
	}
	fault, ok := err.(Fault)
	if !ok {
		t.Fatal("expected error to be of concrete type Fault, but got", err)
	}
	if fault.Code != -32602 {
		t.Errorf("wrong fault code: %d", fault.Code)
	}
	if fault.String != "Wrong Arguments Number" {
		t.Errorf("wrong fault string: %s", fault.String)
	}

	var res2 FaultTestBadResponse
	err = execute(t, s, "FaultTest.Multiply", &FaultTestRequest{4, 2}, &res2)
	if err == nil {
		t.Fatal("expected err to be not nil, but got:", err)
	}
	fault, ok = err.(Fault)
	if !ok {
		t.Fatal("expected error to be of concrete type Fault, but got", err)
	}
	if fault.Code != -32602 {
		t.Errorf("wrong fault code: %d", fault.Code)
	}

	if !strings.HasPrefix(fault.String, "Invalid Method Parameters: fields type mismatch") {
		t.Errorf("wrong response: %s", fault.String)
	}
}

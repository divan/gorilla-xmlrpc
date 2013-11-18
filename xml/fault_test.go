// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"github.com/gorilla/rpc"
	"net/http"
	"strings"
	"testing"
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
		t.Error("Expected err to be not nil, but got:", err)
	}
	if err.Error() != "Wrong number of arguments" {
		t.Errorf("Wrong response: %v.", err.Error())
	}

	var res2 FaultTestBadResponse
	err = execute(t, s, "FaultTest.Multiply", &FaultTestRequest{4, 2}, &res2)
	if err == nil {
		t.Error("Expected err to be not nil, but got:", err)
	}
	if !strings.HasPrefix(err.Error(), "Fields type mismatch") {
		t.Errorf("Wrong response: %v.", err.Error())
	}
}

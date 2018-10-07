//line fault_test.go:1
// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !jex
//jex:off

package xml

//line fault_test.go:8
import _jex "github.com/anjensan/jex/runtime"

//line fault_test.go:12
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
	A	int
	B	int
}

type FaultTestBadRequest struct {
	A	int
	B	int
	C	int
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
//line fault_test.go:53
	_jex.TryCatch(func() {
//line fault_test.go:56
		var res1 FaultTestResponse
		execute_(t, s, "FaultTest.Multiply", &FaultTestBadRequest{4, 2, 4}, &res1)
		t.Fatal("expected err to be not nil")
	}, func(_jex_ex _jex.Exception) {
//line fault_test.go:59
		defer _jex.Suppress(_jex_ex)
		err := _jex_ex.Err()
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
	})
//line fault_test.go:71
	_jex.TryCatch(func() {
//line fault_test.go:74
		var res2 FaultTestBadResponse
		execute_(t, s, "FaultTest.Multiply", &FaultTestRequest{4, 2}, &res2)
		t.Fatal("expected err to be not nil")
	}, func(_jex_ex _jex.Exception) {
//line fault_test.go:77
		defer _jex.Suppress(_jex_ex)
		err := _jex_ex.Err()
		fault, ok := err.(Fault)
		if !ok {
			t.Fatal("expected error to be of concrete type Fault, but got", err)
		}
		if fault.Code != -32602 {
			t.Errorf("wrong fault code: %d", fault.Code)
		}
		if !strings.HasPrefix(fault.String, "Invalid Method Parameters: fields type mismatch") {
			t.Errorf("wrong response: %s", fault.String)
		}
	})

//line fault_test.go:92
}

//line fault_test.go:92
const _ = _jex.Unused

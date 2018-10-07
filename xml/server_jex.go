//line server.go:1
// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !jex
//jex:off

package xml

//line server.go:8
import _jex "github.com/anjensan/jex/runtime"

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/rpc"
//line server.go:19
)

// ----------------------------------------------------------------------------
// Codec
// ----------------------------------------------------------------------------

// NewCodec returns a new XML-RPC Codec.
func NewCodec() *Codec {
	return &Codec{
		aliases: make(map[string]string),
	}
}

// Codec creates a CodecRequest to process each request.
type Codec struct {
	aliases map[string]string
}

// RegisterAlias creates a method alias
func (c *Codec) RegisterAlias(alias, method string) {
	c.aliases[alias] = method
}

// NewRequest returns a CodecRequest.
func (c *Codec) NewRequest(r *http.Request) (_jex_r0 rpc.CodecRequest) {
//line server.go:43
	var _jex_ret bool
	var request ServerRequest
//line server.go:44
	var _jex_md954 _jex.MultiDefer
//line server.go:44
	defer _jex_md954.Run()
//line server.go:44
	_jex.TryCatch(func() {
//line server.go:46
		rawxml, _jex_e967 := ioutil.ReadAll(r.Body)
//line server.go:46
		_jex.Must(_jex_e967)
//line server.go:46
		{
//line server.go:46
			_f := r.Body.Close
			_jex_md954.Defer(func() {
//line server.go:47
				_f()
//line server.go:47
			})
//line server.go:47
		}
//line server.go:47
		_jex_e1031 := xml.Unmarshal(rawxml, &request)
//line server.go:49
		_jex.Must(_jex_e1031)

		request.rawxml = string(rawxml)
		if method, ok := c.aliases[request.Method]; ok {
			request.Method = method
		}
	}, func(_jex_ex _jex.Exception) {
//line server.go:55
		defer _jex.Suppress(_jex_ex)
//line server.go:55
		_jex_ret, _jex_r0 = true, &CodecRequest{err: _jex_ex.Err()}
		return
	})
//line server.go:57
	if _jex_ret {
//line server.go:57
		return
//line server.go:57
	}
	return &CodecRequest{request: &request}
}

// ----------------------------------------------------------------------------
// CodecRequest
// ----------------------------------------------------------------------------

type ServerRequest struct {
	Name	xml.Name		`xml:"methodCall"`
	Method	string		`xml:"methodName"`
	rawxml	string
}

// CodecRequest decodes and encodes a single request.
type CodecRequest struct {
	request	*ServerRequest
	err	error
}

// Method returns the RPC method for the current request.
//
// The method uses a dotted notation as in "Service.Method".
func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

// ReadRequest fills the request object for the RPC method.
//
// args is the pointer to the Service.Args structure
// it gets populated from temporary XML structure
func (c *CodecRequest) ReadRequest(args interface{}) error {
//line server.go:91
	_jex.TryCatch(func() {
//line server.go:93
		xml2RPC_(c.request.rawxml, args)
	}, func(_jex_ex _jex.Exception) {
//line server.go:94
		defer _jex.Suppress(_jex_ex)
		c.err = _jex_ex.Wrap()
	})
	return nil
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
//
// response is the pointer to the Service.Response structure
// it gets encoded into the XML-RPC xml string
func (c *CodecRequest) WriteResponse(w http.ResponseWriter, response interface{}, methodErr error) error {
	var xmlstr string
//line server.go:105
	_jex.TryCatch(func() {
		var _jex_e2598 error
//line server.go:106
		_jex_e2598 = c.err
		_jex.Must(_jex_e2598)
		xmlstr = rpcResponse2XML_(response)
	}, func(_jex_ex _jex.Exception) {
//line server.go:109
		defer _jex.Suppress(_jex_ex)
		var fault Fault
		switch f := _jex_ex.Err().(type) {
		case Fault:
			fault = f
		default:
			fault = FaultApplicationError
			fault.String += fmt.Sprintf(": %v", f)
		}
		xmlstr = fault2XML(fault)
	})
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xmlstr))
	return nil
}

//line server.go:123
const _ = _jex.Unused

// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"encoding/xml"
	"github.com/gorilla/rpc"
	"io/ioutil"
	"net/http"
)

// ----------------------------------------------------------------------------
// Codec
// ----------------------------------------------------------------------------

// NewCodec returns a new XML-RPC Codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Codec creates a CodecRequest to process each request.
type Codec struct {
}

// NewRequest returns a CodecRequest.
func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	return newCodecRequest(r)
}

// ----------------------------------------------------------------------------
// CodecRequest
// ----------------------------------------------------------------------------

type ServerRequest struct {
	Name   xml.Name `xml:"methodCall"`
	Method string   `xml:"methodName"`
	rawxml string
}

// CodecRequest decodes and encodes a single request.
type CodecRequest struct {
	request *ServerRequest
	err     error
}

func newCodecRequest(r *http.Request) rpc.CodecRequest {
	rawxml, err := ioutil.ReadAll(r.Body)
	var request ServerRequest
	err = xml.Unmarshal(rawxml, &request)
	request.rawxml = string(rawxml)

	r.Body.Close()
	return &CodecRequest{request: &request, err: err}
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
	c.err = XML2RPC(c.request.rawxml, args)
	return nil
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
//
// response is the pointer to the Service.Response structure
// it gets encoded into the XML-RPC xml string
func (c *CodecRequest) WriteResponse(w http.ResponseWriter, response interface{}, methodErr error) error {
	var xmlstr string
	if c.err != nil {
		fault := Fault{1, c.err.Error()}
		xmlstr = Fault2XML(fault)
	} else {
		xmlstr, _ = RPCResponse2XML(response)
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xmlstr))
	return nil
}

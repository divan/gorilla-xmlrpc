// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"io"
	"io/ioutil"
)

// EncodeClientRequest encodes parameters for a XML-RPC client request.
func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	xml, err := RPCRequest2XML(method, args)
	return []byte(xml), err
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse(r io.Reader, reply interface{}) (err error) {
	rawxml, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = XML2RPC(string(rawxml), reply)
	return
}

// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build jex
//go:generate jex


package xml

import (
	"io"
	"io/ioutil"

	. "github.com/anjensan/jex"
)

// EncodeClientRequest encodes parameters for a XML-RPC client request.
func EncodeClientRequest_(method string, args interface{}) []byte {
	return []byte(rpcRequest2XML_(method, args))
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse_(r io.Reader, reply interface{}) {
	var rawxml []byte
	if TRY() {
		rawxml, ERR = ioutil.ReadAll(r)
	} else {
		THROW(FaultSystemError)
	}
	xml2RPC_(string(rawxml), reply)
}

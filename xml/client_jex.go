//line client.go:1
// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !jex
//jex:off

//line client.go:9
package xml

//line client.go:9
import _jex "github.com/anjensan/jex/runtime"

import (
	"io"
	"io/ioutil"
//line client.go:16
)

// EncodeClientRequest encodes parameters for a XML-RPC client request.
func EncodeClientRequest_(method string, args interface{}) []byte {
	return []byte(rpcRequest2XML_(method, args))
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse_(r io.Reader, reply interface{}) {
	var rawxml []byte
//line client.go:26
	_jex.TryCatch(func() {
		var _jex_e627 error
		rawxml, _jex_e627 = ioutil.ReadAll(r)
//line client.go:28
		_jex.Must(_jex_e627)
	}, func(_jex_ex _jex.Exception) {
//line client.go:29
		defer _jex.Suppress(_jex_ex)
//line client.go:29
		panic(_jex.NewException(FaultSystemError))
//line client.go:31
	})
	xml2RPC_(string(rawxml), reply)
}

//line client.go:33
const _ = _jex.Unused

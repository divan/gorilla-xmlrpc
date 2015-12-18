// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import "io"

// EncodeClientRequest encodes parameters for a XML-RPC client request.
func EncodeClientRequest(w io.Writer, method string, args interface{}) error {
	return rpcRequest2XML(w, method, args)
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse(r io.Reader, reply interface{}) error {
	return xml2RPC(r, reply)
}

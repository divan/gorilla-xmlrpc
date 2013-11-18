// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

type Fault struct {
	Code   int    `xml:"faultCode"`
	String string `xml:"faultString"`
}

func Fault2XML(fault Fault) string {
	buffer := "<methodResponse><fault>"
	xml, _ := RPC2XML(fault)
	buffer += xml
	buffer += "</fault></methodResponse>"
	return buffer
}

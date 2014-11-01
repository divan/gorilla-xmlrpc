// Copyright 2013 Ivan Danyliuk
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"fmt"
)

type Fault struct {
	Code   int    `xml:"faultCode"`
	String string `xml:"faultString"`
}

func (f Fault) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

func Fault2XML(fault Fault) string {
	buffer := "<methodResponse><fault>"
	xml, _ := RPC2XML(fault)
	buffer += xml
	buffer += "</fault></methodResponse>"
	return buffer
}

type FaultValue struct {
	Value Value `xml:"value"`
}

// IsEmpty returns true if FaultValue contain fault.
//
// FaultValue should be a struct with 2 members.
func (f FaultValue) IsEmpty() bool {
	return len(f.Value.Struct) == 0
}

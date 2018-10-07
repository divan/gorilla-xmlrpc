//+build jex
//go:generate jex

package main

import . "github.com/anjensan/jex"

import (
    "log"
    "bytes"
    "net/http"
    "github.com/divan/gorilla-xmlrpc/xml"
)

func XmlRpcCall_(method string, args struct{Who string}) (reply struct{Message string}) {
    buf := xml.EncodeClientRequest_(method, &args)
    resp, ERR := http.Post("http://localhost:1234/RPC2", "text/xml", bytes.NewBuffer(buf))
    defer resp.Body.Close()
    xml.DecodeClientResponse_(resp.Body, &reply)
	return
}

func main() {
	if TRY() {
		reply := XmlRpcCall_("HelloService.Say", struct{Who string}{"User 1"})
		log.Printf("Response: %s\n", reply.Message)
	} else {
        log.Fatal(EX())
	}
}

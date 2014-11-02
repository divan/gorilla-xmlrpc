package main

import (
    "log"
    "bytes"
    "net/http"
    "github.com/divan/gorilla-xmlrpc/xml"
)

func XmlRpcCall(method string, args struct{Who string}) (reply struct{Message string}, err error) {
    buf, _ := xml.EncodeClientRequest(method, &args)

    resp, err := http.Post("http://localhost:1234/RPC2", "text/xml", bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    defer resp.Body.Close()

    err = xml.DecodeClientResponse(resp.Body, &reply)
    return
}

func main() {
    reply, err := XmlRpcCall("HelloService.Say", struct{Who string}{"User 1"})
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Response: %s\n", reply.Message)
}
# gorilla-xmlrpc #

[![Build Status](https://drone.io/github.com/divan/gorilla-xmlrpc/status.png)](https://drone.io/github.com/divan/gorilla-xmlrpc/latest)
[![GoDoc](https://godoc.org/github.com/divan/gorilla-xmlrpc/xml?status.svg)](https://godoc.org/github.com/divan/gorilla-xmlrpc/xml)

XML-RPC implementation for the Gorilla/RPC toolkit.

It implements both server and client.

It's built on top of gorilla/rpc package in Go(Golang) language and implements XML-RPC, according to [it's specification](http://xmlrpc.scripting.com/spec.html).
Unlike net/rpc from Go strlib, gorilla/rpc allows usage of HTTP POST requests for RPC.

### Installation ###
Assuming you already imported gorilla/rpc, use the following command:

    go get github.com/divan/gorilla-xmlrpc/xml

### Examples ###

#### Server Example ####

```go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/divan/gorilla-xmlrpc/xml"
)

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *struct{Who string}, reply *struct{Message string}) error {
    log.Println("Say", args.Who)
    reply.Message = "Hello, " + args.Who + "!"
    return nil
}

func main() {
    RPC := rpc.NewServer()
    xmlrpcCodec := xml.NewCodec()
    RPC.RegisterCodec(xmlrpcCodec, "text/xml")
    RPC.RegisterService(new(HelloService), "")
    http.Handle("/RPC2", RPC)

    log.Println("Starting XML-RPC server on localhost:1234/RPC2")
    log.Fatal(http.ListenAndServe(":1234", nil))
}
```

It's pretty self-explanatory and can be tested with any xmlrpc client, even raw curl request:

```bash
curl -v -X POST -H "Content-Type: text/xml" -d '<methodCall><methodName>HelloService.Say</methodName><params><param><value><struct><member><name>Who</name><value><string>XMLTest</string></value></member></struct></value></param><param><value><struct><member><name>Code</name><value><int>123</int></value></member></struct></value></param></params></methodCall>' http://localhost:1234/RPC2
```

#### Client Example ####

Implementing client is beyond the scope of this package, but with encoding/decoding handlers it should be pretty trivial. Here is an example which works with the server introduced above.

```go
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

```

### Implementation details ###

The main objective was to use standard encoding/xml package for XML marshalling/unmarshalling. Unfortunately, in current implementation there is no graceful way to implement common structre for marshal and unmarshal functions - marshalling doesn't handle interface{} types so far (though, it could be changed in the future).
So, marshalling is implemented manually.

Unmarshalling code first creates temporary structure for unmarshalling XML into, then converts it into the passed variable using *reflect* package.
If XML struct member's name is lowercased, it's first letter will be uppercased, as in Go/Gorilla field name must be exported(first-letter uppercased).

Marshalling code converts rpc directly to the string XML representation.

For the better understanding, I use terms 'rpc2xml' and 'xml2rpc' instead of 'marshal' and 'unmarshall'.

### Supported types ###

| XML-RPC          | Golang        |
| ---------------- | ------------- |
| int, i4          | int           |
| double           | float64       |
| boolean          | bool          |
| string           | string        |
| dateTime.iso8601 | time.Time     |
| base64           | []byte        |
| struct           | struct        |
| array            | []interface{} |
| nil              | nil           |

### TODO ###

*  Add more corner cases tests


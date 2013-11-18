# gorilla-xmlrpc #

This is an XML-RPC protocol implementation for the Gorilla/RPC toolkit.

It's built on top of gorilla/rpc package in Go(Golang) language and implements XML-RPC, according to [it's specifiaction](http://xmlrpc.scripting.com/spec.html).

So far it doesn't handle Faults/error correctly (as required by XML-RPC spec), but the work on it in progress.

**NOTE: I hope this code soon will be part of Gorilla toolkit, so the path and the name will slightly change**

### Installing ###
Assuming you already imported gorilla/rpc, use the following command:

    go get github.com/divan/gorilla-xmlrpc/xml

### Examples ###
	package main

	import (
		"log"
		"net/http"
		"github.com/gorilla/rpc"
		"github.com/divan/gorilla-xmlrpc/xml"
	)

	type HelloArgs struct {
		Who string
	}

	type HelloReply struct {
		Message string
		Status int
	}

	type HelloService struct{}

	func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
		log.Println("Say", args.Who)
		reply.Message = "Hello, " + args.Who + "!"
		reply.Status = 42
		return nil
	}

	func main() {
		RPC := rpc.NewServer()
		xmlrpcCodec := xml.NewCodec()
		RPC.RegisterCodec(xmlrpcCodec, "text/xml")
		RPC.RegisterService(new(HelloService), "")
		http.Handle("/RPC2", RPC)

		log.Println("Starting XML-RPC server on localhost:1234/RPC2")
		err := http.ListenAndServe(":1234", nil)
		if err != nil {
			log.Fatal("ListenAndServer: ", err)
		}
	}

It's pretty self-explanatory and can be tested with any xmlrpc client, even raw curl request:

   curl -v -X POST -H "Content-Type: text/xml" \
	       -d '<methodCall><methodName>HelloService.Say</methodName><params><param><value><struct><member><name>Who</name><value><string>XMLTest</string></value></member></struct></value></param><param><value><struct><member><name>Code</name><value><int>123</int></value></member></struct></value></param></params></methodCall>' \
		       http://localhost:1234/api

## TODO ##




# gorilla-xmlrpc #

This is an XML-RPC protocol implementation for the Gorilla/RPC toolkit.

It's built on top of gorilla/rpc package in Go(Golang) language and implements XML-RPC, according to [it's specifiaction](http://xmlrpc.scripting.com/spec.html).
Unlike Go standard net/rpc, gorilla/rpc allows usage HTTP POST requests for RPC.

So far it doesn't handle Faults/error correctly (as required by XML-RPC spec), but the work on it in progress.


### Installation ###
Assuming you already imported gorilla/rpc, use the following command:

    go get github.com/divan/gorilla-xmlrpc/xml

**NOTE: I hope this code soon will be part of Gorilla toolkit, so the path and the name will slightly change**

### Examples ###

#### Server Example ####

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

    curl -v -X POST -H "Content-Type: text/xml" -d '<methodCall><methodName>HelloService.Say</methodName><params><param><value><struct><member><name>Who</name><value><string>XMLTest</string></value></member></struct></value></param><param><value><struct><member><name>Code</name><value><int>123</int></value></member></struct></value></param></params></methodCall>' http://localhost:1234/RPC2


#### Client Example ####

Implementing client is beyound of scope of this package, but with encoding/decoding handlers it should be pretty trivial. Here is the example which works with the example server introduced above.

package main

	import (
		"log"
		"bytes"
		"net/http"
		"github.com/divan/gorilla-xmlrpc/xml"
	)

	type HelloArgs struct {
		Who string
	}

	type HelloReply struct {
		Message string
		Status int
	}

	func XmlRpcCall(method string, args HelloArgs) (reply HelloReply, err error) {
		buf, _ := xml.EncodeClientRequest(method, &args)
		body := bytes.NewBuffer(buf)

		resp, err := http.Post("http://localhost:1234/RPC2", "text/xml", body)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		xml.DecodeClientResponse(resp.Body, &reply)
		return
	}

	func main() {
		args := HelloArgs{"User1"}
		var reply HelloReply

		reply, err := XmlRpcCall("HelloService.Say", args)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Response: %s (%d)\n", reply.Message, reply.Status)
	}

### Implementation details ###

The main objective was to use standard encoding/xml package for XML marshalling/unmarshalling. Unfortunately, in current implementation there is no graceful way to implement common structre for marshal and unmarshal functions - marshalling doesn't handle interface{} types so far (though, it could be changed in the future).
So, marshalling is implemented manually.

Unmarshalling code first creates temporary structure for unmarshalling XML into, then converts it into the passed variable using *reflect* package.

Marshalling code converts rpc directly to the string XML representation.

For the better understanding, I use terms 'rpc2xml' and 'xml2xml' instead of 'marshal' and 'unmarshall'.

### TODO ###

*   Time / dateTime.iso8601 support
*   Base64  support
*   Fault support according to XML-RPC spec (it will require some changes in gorilla/rpc module, will be discussed)
*   Make/find tests that cover corner cases for XML-RPC


//line example_server.go:1
//+build !jex
//jex:off

package main

//line example_server.go:4
import _jex "github.com/anjensan/jex/runtime"

//line example_server.go:8
import (
	"log"
	"net/http"
	"github.com/gorilla/rpc"
	"github.com/divan/gorilla-xmlrpc/xml"
)

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *struct{ Who string }, reply *struct{ Message string }) error {
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

//line example_server.go:32
const _ = _jex.Unused

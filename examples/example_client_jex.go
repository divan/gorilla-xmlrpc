//line example_client.go:1
//+build !jex
//jex:off

package main

//line example_client.go:4
import _jex "github.com/anjensan/jex/runtime"

//line example_client.go:8
import (
	"log"
	"bytes"
	"net/http"
	"github.com/divan/gorilla-xmlrpc/xml"
)

func XmlRpcCall_(method string, args struct{ Who string }) (reply struct{ Message string }) {
	buf := xml.EncodeClientRequest_(method, &args)
	resp, _jex_e319 := http.Post("http://localhost:1234/RPC2", "text/xml", bytes.NewBuffer(buf))
//line example_client.go:17
	_jex.Must(_jex_e319)
	defer resp.Body.Close()
	xml.DecodeClientResponse_(resp.Body, &reply)
	return
}

func main() {
//line example_client.go:23
	_jex.TryCatch(func() {
//line example_client.go:25
		reply := XmlRpcCall_("HelloService.Say", struct{ Who string }{"User 1"})
		log.Printf("Response: %s\n", reply.Message)
	}, func(_jex_ex _jex.Exception) {
//line example_client.go:27
		defer _jex.Suppress(_jex_ex)
		log.Fatal(_jex_ex)
	})
}

//line example_client.go:30
const _ = _jex.Unused

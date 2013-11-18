# gorilla-xmlrpc #

This is an XML-RPC protocol implementation for the Gorilla/RPC toolkit.

It's built on top of gorilla/rpc package in Go(Golang) language and implements XML-RPC, according to [it's specifiaction](http://xmlrpc.scripting.com/spec.html).

So far it doesn't handle Faults/error correctly (as required by XML-RPC spec), but the work on it in progress.

**NOTE: I hope this code soon will be part of Gorilla toolkit, so the path and the name will slightly change**

### Installing ###
Assuming you already imported gorilla/rpc, use the following command:

     go get github.com/divan/gorilla-xmlrpc/xml

### Examples ###


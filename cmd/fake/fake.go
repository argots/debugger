// Package main illustrates how to run a debugger server that chrome
// dev tools can connect to.
//
// To start the server, run "go run ./cmd/fake/".  The server uses
// localhost:8222 port -- point chrome at this by visiting
// chrome://inspect/#devices, click on "Open dedicated DevTools for
// Node" and then type in "localhost:8222" in the text input.
package main

import (
	"log"
	"net/http"

	"github.com/argots/debugger"
	"github.com/argots/debugger/dispatch"
)

func main() {
	h := debugger.Handler(NewFakeServer())
	log.Println("Listening on :8222")
	log.Fatal(http.ListenAndServe(":8222", h))
}

var _ debugger.Driver = &FakeServer{}

type FakeServer struct {
	nw debugger.Network
	debugger.Dispatcher
	target *debugger.Target
}

func NewFakeServer() *FakeServer {
	target := &debugger.Target{
		Description: "Fake Debugger Description",
		FaviconURL:  "",
		Title:       "Fake Debugger Title",
		Type:        "Fake Type",
		URL:         "http://example.com/fakeUrl",
	}
	target.Init("http://localhost:8222", "FakeID")
	return &FakeServer{
		nw:         &debugger.FakeNetwork{},
		Dispatcher: dispatch.New(),
		target:     target,
	}
}

func (f *FakeServer) Name() string {
	return "FakeServer"
}

func (f *FakeServer) List() ([]*debugger.Target, error) {
	return []*debugger.Target{f.target}, nil
}

func (f *FakeServer) Network() debugger.Network {
	return f.nw
}

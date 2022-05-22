package http

import (
	"github.com/leor-w/kid/server"
	"net/http"
	"sync"
	"testing"
)

func TestHttp(t *testing.T) {
	srv := NewServer(server.Address(":8090"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	hd := srv.NewHandler(mux)

	if err := srv.Handle(hd); err != nil {
		t.Fatal(err)
	}

	wait := sync.WaitGroup{}
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	wait.Add(1)
	wait.Wait()
}

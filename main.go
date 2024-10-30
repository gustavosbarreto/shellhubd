package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func newServerMux(h http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h)
	return mux
}

func main() {
	listener, err := createListener()
	if err != nil {
		log.Fatalf("Erro ao criar listener: %v", err)
	}
	defer listener.Close()

	srv := &http.Server{
		Handler:     newServerMux(handleRequest),
		ConnContext: connContext,
	}

	fmt.Printf("Servidor rodando em %s...\n", runtime.GOOS)

	log.Fatal(srv.Serve(listener))
}

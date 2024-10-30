// +build windows

package main

import (
    "fmt"
    "net"
    "github.com/Microsoft/go-winio"
    "net/http"
)

// Criação do listener no Windows usando Named Pipes com go-winio
func createListener() (net.Listener, error) {
    pipeName := `\\.\pipe\shellhubd_pipe`
    return winio.ListenPipe(pipeName, nil)
}

// Handler específico para Windows: processa o request e pega informações do cliente via Named Pipes
func handleRequest(w http.ResponseWriter, r *http.Request) {
    if conn, ok := r.Body.(net.Conn); ok {
        if pc, ok := conn.(*winio.PipeConn); ok {
            token, err := pc.GetClientToken()
            if err != nil {
                fmt.Fprintf(w, "Erro ao obter token do cliente: %v", err)
                return
            }
            defer token.Close()

            // Verifica se o cliente é um administrador
            isAdmin, err := token.IsMember(winio.BuiltinAdministratorsSid)
            if err != nil {
                fmt.Fprintf(w, "Erro ao verificar se o cliente é administrador: %v", err)
                return
            }

            fmt.Fprintf(w, "Usuário é administrador: %v", isAdmin)
        } else {
            fmt.Fprintf(w, "Erro ao converter conexão para PipeConn")
        }
    } else {
        fmt.Fprintf(w, "Erro ao converter body para net.Conn")
    }
}

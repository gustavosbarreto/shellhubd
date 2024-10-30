//go:build linux
// +build linux

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"

	"golang.org/x/sys/unix"
)

type contextKey string

const connKey contextKey = "netConn"

// Criação do listener no Linux usando Unix Domain Sockets
func createListener() (net.Listener, error) {
	socketPath := "/tmp/shellhubd.sock"
	os.Remove(socketPath)
	return net.Listen("unix", socketPath)
}

func connContext(ctx context.Context, c net.Conn) context.Context {
	file, _ := c.(*net.UnixConn).File()
	credentials, _ := unix.GetsockoptUcred(int(file.Fd()), unix.SOL_SOCKET, unix.SO_PEERCRED)
	return context.WithValue(ctx, connKey, credentials)
}

// Handler específico para Linux: pega informações do cliente via Unix Domain Sockets
func handleRequest(w http.ResponseWriter, r *http.Request) {
	ucred, ok := r.Context().Value(connKey).(*unix.Ucred)
	if !ok {
		http.Error(w, "Conexão não encontrada", http.StatusInternalServerError)
		return
	}

	uid := ucred.Uid
	gid := ucred.Gid

	fmt.Fprintf(w, "UID: %d, GID: %d\n", uid, gid)

	wheel, err := user.LookupGroup("wheel")
	if err != nil {
		fmt.Fprintf(w, "Erro ao obter informações do grupo wheel: %v", err)
		return
	}

	reqUser, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		fmt.Fprintf(w, "Erro ao obter informações do usuário: %v", err)
		return
	}

	groups, err := reqUser.GroupIds()
	if err != nil {
		fmt.Fprintf(w, "Erro ao obter grupos do usuário: %v", err)
		return
	}

	isWheel := false

	for _, group := range groups {
		if group == wheel.Gid {
			isWheel = true
			break
		}
	}

	fmt.Fprintf(w, "Usuário é do grupo wheel: %v", isWheel)
}

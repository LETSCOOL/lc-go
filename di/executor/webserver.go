package executor

import (
	"context"
	"fmt"
	"golang.org/x/net/netutil"
	. "lc-go/lg"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultWebServerPort = 8000
)

type WebServer struct {
	address string
	port    int
	maxConn int
}

type RoutingPathFunc struct {
	path    string `http:"path"`
	method  string `http:"method"`
	handler func()
}

func (ws *WebServer) Route() {

}

//func (ws *WebServer) ServeHTTP(writer http.ResponseWriter,
//	request *http.Request) {
//	//io.WriteString(writer, sh.message)
//	ctx := request.Context()
//}

func (ws *WebServer) handleRequest(writer http.ResponseWriter, r *http.Request) {
	log.Printf("got request from %s\n", r.RemoteAddr)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("you got it"))
}

func (ws *WebServer) Run() {
	addr := fmt.Sprintf("%v:%d", ws.address, Ife(ws.port <= 0, DefaultWebServerPort, ws.port))
	router := http.NewServeMux()
	router.HandleFunc("/", ws.handleRequest)

	srv := http.Server{
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 10,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if ws.maxConn > 0 {
		listener = netutil.LimitListener(listener, ws.maxConn)
		//log.Printf("max connections set to %d\n", ws.maxConn)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			//
		}
	}(listener)

	log.Printf("listening on %s\n", listener.Addr().String())

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	log.Printf("interrupted, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v\n", err)
	}
}

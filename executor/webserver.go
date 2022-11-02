package executor

import (
	"context"
	"fmt"
	"golang.org/x/net/netutil"
	"lc-go/dij"
	. "lc-go/lg"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

const (
	DefaultWebServerPort = 8000
	TagName              = "http"
	WebConfigKey         = "webserver.config"
)

type WebServerConfig struct {
	address string
	port    int
	maxConn int
}

func LaunchWebServer(webServerType reflect.Type, others ...any) {
	ref := map[dij.DependencyKey]any{}
	config := WebServerConfig{}
	for _, other := range others {
		otherTyp := reflect.TypeOf(other)
		switch v := other.(type) {
		case WebServerConfig:
			config = v
			ref[WebConfigKey] = v
		default:
			log.Println("No ideal about type:", otherTyp)
		}
	}
	_, err := dij.CreateInstance(webServerType, &ref, "")
	if err != nil {
		log.Panic(err)
	}

	addr := fmt.Sprintf("%v:%d", config.address, Ife(config.port <= 0, DefaultWebServerPort, config.port))
	router := http.NewServeMux()
	//router.HandleFunc("/", ws.HandleRequest)
	// TODO: parse controller handle functions

	srv := http.Server{
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 10,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if config.maxConn > 0 {
		listener = netutil.LimitListener(listener, config.maxConn)
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

package executor

import (
	"context"
	"fmt"
	"golang.org/x/net/netutil"
	"lc-go/di"
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
)

type WebServerConfig struct {
	address string
	port    int
	maxConn int
}

type IsWebServer interface {
	GetConfig() WebServerConfig
	HandleRequest(writer http.ResponseWriter, r *http.Request)
}

type WebServer struct {
	config WebServerConfig
}

func (ws *WebServer) GetConfig() WebServerConfig {
	return ws.config
}

type RoutingPathFunc struct {
	path    string `http:"path"`
	method  string `http:"method"`
	handler func()
}

func (ws *WebServer) Config(config WebServerConfig) {
	ws.config = config
}

//func (ws *WebServerConfig) ServeHTTP(writer http.ResponseWriter,
//	request *http.Request) {
//	//io.WriteString(writer, sh.message)
//	ctx := request.Context()
//}

func (ws *WebServer) HandleRequest(writer http.ResponseWriter, r *http.Request) {
	log.Printf("got request from %s\n", r.RemoteAddr)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("you got it"))
}

func LaunchWebServer(webServerType reflect.Type, others ...any) {
	//wsTyp := reflect.TypeOf(*ws)
	ref := map[di.DependencyKey]any{}
	config := WebServerConfig{}
	for _, other := range others {
		otherTyp := reflect.TypeOf(other)
		switch v := other.(type) {
		case WebServerConfig:
			config = v
			ref["config"] = v
		default:
			log.Println("No ideal about type:", otherTyp)
		}
	}
	inst, err := di.CreateDependency(webServerType, &ref)
	if err != nil {
		log.Panic(err)
	}

	ws, ok := inst.(IsWebServer)
	if !ok {
		log.Fatalf("the type '%v' is not web server structure", reflect.TypeOf(inst))
	}

	addr := fmt.Sprintf("%v:%d", config.address, Ife(config.port <= 0, DefaultWebServerPort, config.port))
	router := http.NewServeMux()
	router.HandleFunc("/", ws.HandleRequest)

	srv := http.Server{
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 10,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if ws.GetConfig().maxConn > 0 {
		listener = netutil.LimitListener(listener, ws.GetConfig().maxConn)
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

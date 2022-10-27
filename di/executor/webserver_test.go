package executor

import (
	"reflect"
	"testing"
)

type TWebServer struct {
	WebServer
}

func (t *TWebServer) InjectDependency(p struct {
	config WebServerConfig
	ctl    *TWebCtl
}) {

}

type TWebCtl struct {
}

func (c *TWebCtl) Route(p struct {
	path   string `http:"path=/"`
	method string `http:"method"`
}) {

}

// go test lc-go/di/executor -v -run TestWeb
func TestWeb(t *testing.T) {
	//ws := &TWebServer{}
	config := WebServerConfig{
		address: "",
		port:    0,
		maxConn: 0,
	}
	//ws.Route()
	LaunchWebServer(reflect.TypeOf(TWebServer{}), config)
}

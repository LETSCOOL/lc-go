package executor

import (
	"lc-go/dij"
	"reflect"
	"testing"
)

type TWebServer struct {
	_ *TAuthRoute `di:"" http:"/auth,"`
	_ *TUserRoute `di:"" http:"/user,"`
}

type TAuthRoute struct {
}

func (a *TAuthRoute) GetToken() {

}

type TUserRoute struct {
}

func (u *TUserRoute) GetUser() {

}

// go test lc-go/executor -v -run TestWeb
func TestWeb(t *testing.T) {
	dij.EnableLog()
	LaunchWebServer(reflect.TypeOf(TWebServer{}))
}

package zgo

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
)

type ZgoRouter struct {
	*gin.Engine
}


type ZgoContext struct {
	Ctx *gin.Context
}

var ZCtx ZgoContext

func init()  {
	ZCtx = ZgoContext{Ctx:&gin.Context{}}
}

func Init(objfunc Objfunc)  {
	router := GinPreWork()
	//installRouter(ZgoRouter{router})
	objfunc(&ZgoRouter{router})
	GinPostWork(router)
}

func GinPreWork() *gin.Engine {
	gin.DefaultWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)
	return gin.Default()
}

func GinPostWork(router *gin.Engine) error {
	server := &http.Server{
		Addr: net.JoinHostPort(
			"0.0.0.0",
			"7777",
		),
		Handler: router,
	}
	return server.ListenAndServe()
}

func installRouter(z *ZgoRouter){
	z.Group("v1")
	z.GET("hello", hello)
}

type Objfunc func(z *ZgoRouter)

func hello(z *gin.Context)  {
	z.JSON(200, "hello" )
}
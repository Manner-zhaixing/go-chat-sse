package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/biz"
	"net/http"

	"go-chat-sse/internal/config"
	"go-chat-sse/internal/handler"
	"go-chat-sse/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	c.Log.Level = "info"
	logx.SetLevel(logx.InfoLevel)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// httpx.SetErrorHandler 仅在调用了 httpx.Error 处理响应时才有效。
	httpx.SetErrorHandler(func(err error) (int, any) {
		var e *biz.Error
		switch {
		case errors.As(err, &e):
			return http.StatusOK, biz.Fail(e)
		default:
			return http.StatusInternalServerError, nil
		}
	})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

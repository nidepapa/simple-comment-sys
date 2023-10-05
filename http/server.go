/* Date:10/4/23-2023
 * Author:fanzhao
 */
package http

import (
	"context"
	"log"
	"net/http"
	"simple-comment-sys/service"
)

var (
	svc     *service.Service
	authSvc *middleware.CooperationAuth
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
	var (
		cfg     bm.ServerConfig
		authCfg configs.Config
		ct      paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	// auth
	if err = paladin.Get("application.toml").UnmarshalTOML(&authCfg); err != nil {
		panic(err)
	}
	authSvc = middleware.NewAuth(nil, &authCfg)
	tip.Init(nil)
	engine = bm.DefaultServer(&cfg)
	initRouter(engine)
	ossapi.RegisterOssBMServer(engine, svc)
	fileapi.RegisterFileBMServer(engine, svc)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)

	c := e.Group("/x/coop/comment")
	{
		c.POST("/add", authSvc.Valid(), addComment)
		c.POST("/del", delComment)
		c.POST("/state", authSvc.Valid(), setCommentState)
		c.GET("/list", authSvc.Valid(), listComments)
		c.POST("/export", exportComments)
		c.GET("/export/file", getExportCommentsFile)
	}

}

func ping(ctx context.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

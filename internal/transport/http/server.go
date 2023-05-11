package http

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	ErrResponseStatus     = "error"
	SuccessResponseStatus = "success"
)

type Router interface {
	Route(rg *gin.RouterGroup)
}

type GinServer struct {
	config *Config
	server *http.Server
	router *gin.Engine
}

func NewGinServer(config *Config) *GinServer {

	router := gin.Default()

	locationConf := location.Config{
		Host:   net.JoinHostPort(config.Host, config.Port),
		Scheme: config.Scheme,
	}
	corsConf := cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
	}
	router.Use(location.New(locationConf), cors.New(corsConf))

	serv := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: router,
	}

	return &GinServer{
		config: config,
		server: serv,
		router: router,
	}
}

func (gs *GinServer) gracefulPowerOff() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gs.server.Shutdown(ctx); err != nil {
		log.Printf("error shutting down server %s", err)
	} else {
		log.Println("Server stopping...")
	}

	log.Println("Timeout of 5 seconds")
	select {
	case <-ctx.Done():
	}
	log.Println("Server successfully stopped")
}

func (gs *GinServer) prepareHealthcheck() {
	gs.router.GET("/healthcheck", func(ctx *gin.Context) {
		message := "Connected"
		ctx.JSON(http.StatusOK, gin.H{"status": SuccessResponseStatus, "message": message})
	})
}

func (gs *GinServer) serve() {
	if err := gs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (gs *GinServer) StartGinServer() {

	// ctx for graceful shutdown server
	ctxServ, cancelServ := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelServ()

	gs.prepareHealthcheck()
	go gs.serve()

	select {
	case <-ctxServ.Done():
		cancelServ()
		gs.gracefulPowerOff()
	}
}

func (gs *GinServer) StartAllRoutes(routers ...Router) {
	for _, r := range routers {
		r.Route(gs.router.Group("/"))
	}
}

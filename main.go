package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"test/internal/config"
	"test/internal/handler"
	"test/internal/model"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	router := gin.Default()
	config.LoadConfig()
	model.InitDB()
	u := handler.NewUserHandler()
	u.Register(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	srv := &http.Server{
		Addr:    ":" + config.Cfg.Listen,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("press ctrl+c again to force shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

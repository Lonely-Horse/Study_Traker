package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	mu sync.RWMutex
)

func main() {
	filename := "studylog.json"

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:8081", "地址")
	flag.Parse()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		HandlerDashboard(w, r, filename)
	})

	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		HandlerLogs(w, r, filename)
	})

	mux.HandleFunc("/api/skill", func(w http.ResponseWriter, r *http.Request) {
		HandlerSkill(w, r, filename)
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("server listening on %s\n", addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("error:%s", err)
			return
		}
	}()

	<-ctx.Done()
	fmt.Println("\n服务接收到退出操作，正在退出...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println("程序完成退出操作!")

}

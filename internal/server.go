package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func HttpServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/_health", func(writer http.ResponseWriter, request *http.Request) {
		info := make(map[string]string)
		info["build"] = Build
		info["compile"] = Compile
		info["version"] = Version

		marshal, _ := json.Marshal(info)
		_, _ = fmt.Fprintf(writer, string(marshal))
	})

	conf := GetConfig()
	svr := &http.Server{
		Addr:         ":" + conf.HttpPort,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 5 * 60,
	}

	go func() {
		log.Printf("agent server listen at http://127.0.0.1%s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error ListenAndServe : %s", err)
			return
		}
	}()

	gracefulShutdown(ctx, svr)
}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	now := time.Now()

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("start server shutdown...")
	if err := server.Shutdown(timeout); err != nil {
		log.Fatal("server shutdown:", err)
	}
	log.Printf("graceful shutdown server %v\n", time.Since(now))
}

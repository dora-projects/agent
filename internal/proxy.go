package internal

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	conf := Rdb.Get(ctx, "test").String()
	fmt.Println(conf)

	//todo 从配置中获取
	staticPath := "./static/build"
	indexPath := "index.html"

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(staticPath, path)
	accept := r.Header.Get("Accept")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// 返回 index.html
		if strings.Contains(accept, "text/html") || strings.Contains(accept, "text/plain") {
			http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err != nil {
		// 其它文件错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回静态文件
	http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
}

func Proxy() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	svr := &http.Server{
		Addr:         ":" + viper.GetString("port"),
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

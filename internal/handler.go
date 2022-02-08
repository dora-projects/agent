package internal

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func errorRes(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request) {
	conf := GetConfig()
	subHost := GetSubHost(r.Host, conf.WebHost)

	ctx := context.Background()
	bytes, err := Rdb.Get(ctx, "agent:"+subHost).Bytes()
	if err != nil {
		errorRes(w, errors.WithMessage(err, "missing site config"))
		return
	}
	config, err := ParseSiteConfig(bytes)
	if err != nil {
		errorRes(w, err)
		return
	}

	// api proxy
	if config.Proxy != nil {
		for k, v := range config.Proxy {
			fmt.Printf("%v %v  \n", k, v)
		}
	}

	// static
	staticPath := config.Filepath
	indexPath := config.Index

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(staticPath, path)
	accept := r.Header.Get("Accept")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// return index.html
		if strings.Contains(accept, "text/html") || strings.Contains(accept, "text/plain") {
			http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
			return
		}
		errorRes(w, err)
		return
	}

	// file error
	if err != nil {
		errorRes(w, err)
		return
	}

	// static server
	http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
}

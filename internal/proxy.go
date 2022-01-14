package internal

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func spaHandler(w http.ResponseWriter, r *http.Request) {
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", spaHandler)
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}

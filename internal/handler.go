package internal

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func errorRes(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery

	var realUrl string

	director := func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		realUrl = req.URL.String()

		//if _, ok := req.Header["User-Agent"]; !ok {
		//	// explicitly disable User-Agent so it's not set to default value
		//	req.Header.Set("User-Agent", "")
		//}
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(response *http.Response) error {
			response.Header.Set("X-Real-Url", realUrl)

			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
		}}

	return proxy
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

	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// api proxy
	if config.Proxy != nil {
		var targetVal string
		for k, v := range config.Proxy {
			if strings.HasPrefix(path, k) {
				targetVal = v
				break
			}
		}

		// todo: support path rewrite
		// https://webpack.js.org/configuration/dev-server/#devserverproxy
		if targetVal != "" {
			target, _ := url.Parse(targetVal)
			// create the reverse proxy
			proxy := newReverseProxy(target)
			proxy.ServeHTTP(w, r)
			return
		}
	}

	// static
	staticPath := config.Filepath
	indexPath := config.Index

	relPath := filepath.Join(staticPath, path)
	accept := r.Header.Get("Accept")

	_, err = os.Stat(relPath)
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

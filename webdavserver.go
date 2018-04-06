// command webdavserver provides access to given directory via WebDAV protocol
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/artyom/autoflags"
	"golang.org/x/net/webdav"
)

func main() {
	p := struct {
		Dir  string `flag:"dir,directory to serve"`
		Addr string `flag:"addr,address to listen"`
		Auth string `flag:"auth,basic auth credentials in user:password format (or set WEBDAV_AUTH env)"`
	}{
		Dir:  ".",
		Addr: "127.0.0.1:8080",
		Auth: os.Getenv("WEBDAV_AUTH"),
	}
	autoflags.Parse(&p)
	log.Fatal(serve(p.Dir, p.Addr, p.Auth))
}

func serve(dir, addr, auth string) error {
	var handler http.Handler
	webdavHandler := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}
	handler = webdavHandler
	if auth != "" {
		fields := strings.SplitN(auth, ":", 2)
		if len(fields) != 2 {
			return errors.New("invalid auth format (want user:password)")
		}
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok || u != fields[0] || p != fields[1] {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			webdavHandler.ServeHTTP(w, r)
		})
	}
	return (&http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}).ListenAndServe()
}

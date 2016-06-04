command webdavserver provides access to given directory via WebDAV protocol

Use `go get github.com/artyom/webdavserver` to install.

	Usage of webdavserver:
	  -addr string
		address to listen (default "127.0.0.1:8080")
	  -auth string
		basic auth credentials in user:password format (or set WEBDAV_AUTH env)
	  -dir string
		directory to serve (default ".")

Most useful together with [leproxy](https://github.com/artyom/leproxy) for HTTPS support.

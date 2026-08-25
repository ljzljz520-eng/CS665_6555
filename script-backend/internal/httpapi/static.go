package httpapi

import (
	"io/fs"
	"net/http"
)

func (a *API) MountFrontend(frontend fs.FS) {
	a.mux.Handle("/", http.FileServer(http.FS(frontend)))
}

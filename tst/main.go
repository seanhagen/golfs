package gitlfstest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router

func init() {
	r := httprouter.New()
	r.POST("/:org/:repo/locks/verify", VerifyLock)
	r.POST("/:org/:repo/objects/batch", ObjectBatch)

	router = r
}

func GitLFSTest(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

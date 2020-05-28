package golfs

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func setup() (*httprouter.Router, error) {
	s, err := setupService()
	if err != nil {
		log.Printf("unable to get configuration: %v", err)
		return nil, err
	}

	r := httprouter.New()

	// object send/download
	r.POST("/:host/:org/:repo/objects/batch", s.objectBatch)

	//locks
	r.POST("/:host/:org/:repo/locks", s.createLock)
	r.GET("/:host/:org/:repo/locks", s.fetchLocks)
	r.POST("/:host/:org/:repo/locks/verify", s.verifyLocks)
	r.POST("/:host/:org/:repo/locks/:id/unlock", s.deleteLock)

	return r, nil
}

func GOLFS(w http.ResponseWriter, r *http.Request) {
	rt, err := setup()
	if err != nil {
		log.Printf("unable to do function setup: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rt.ServeHTTP(w, r)
}

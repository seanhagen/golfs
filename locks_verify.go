package golfs

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// verifyLocks ...
func (s *service) verifyLocks(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	log.Printf("request to verify locks -- %v", r.URL.Path)
	if !s.locking {
		// not locking, fake response
	}

	w.WriteHeader(http.StatusInternalServerError)
}

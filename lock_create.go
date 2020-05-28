package golfs

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// createLock ...
func (s *service) createLock(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	log.Printf("request to create lock -- %v", r.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
}

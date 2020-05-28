package golfs

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// fetchLocks ...
func (s *service) fetchLocks(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	w.WriteHeader(http.StatusInternalServerError)
}

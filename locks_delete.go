package golfs

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// deleteLock ...
func (s *service) deleteLock(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	w.WriteHeader(http.StatusInternalServerError)
}

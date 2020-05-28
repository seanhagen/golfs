package golfs

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type batchReq struct {
	Operation string `json:"operation"`
	Objects   []struct {
		OID  string `json:"oid"`
		Size int64  `json:"size"`
	} `json:"objects"`
	Transfers []string `json:"transfers"`
	Ref       struct {
		Name string `json:"name"`
	} `json:"ref"`
}

type batchResp struct {
	Transfer string `json:"transfer"`
	Objects  []struct {
		OID           string `json:"oid"`
		Size          int64  `json:"size"`
		Authenticated bool   `json:"authenticated"`
		Actions       struct {
			Download *downloadAction `json:"download,omitempty"`
		} `json:"actions"`
	} `json:"objects"`
}

type downloadAction struct {
	HREF      string
	Header    map[string]string
	ExpiresIn int64
	ExpiresAt string
}

func (s *service) objectBatch(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	log.Printf("request to object batch -- %v", r.URL.Path)
	w.WriteHeader(http.StatusInternalServerError)
}

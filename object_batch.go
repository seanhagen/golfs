package golfs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/davecgh/go-spew/spew"
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

type batchRespObjAction struct {
	Download *downloadAction `json:"download,omitempty"`
}

type batchRespObj struct {
	OID           string             `json:"oid"`
	Size          int64              `json:"size"`
	Authenticated bool               `json:"authenticated"`
	Actions       batchRespObjAction `json:"actions"`
}

type batchResp struct {
	Transfer string         `json:"transfer"`
	Objects  []batchRespObj `json:"objects"`
}

type downloadAction struct {
	HREF      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresIn int64             `json:"expires_in,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

func (s *service) objectBatch(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	ctx := r.Context()
	rid := ctx.Value("rid").(string)

	log.Printf("request to object batch -- %v", r.URL.Path)

	user := ctx.Value("auth-user").(string)
	pass := ctx.Value("auth-pass").(string)

	host := pr.ByName("host")
	org := pr.ByName("org")
	repo := pr.ByName("repo")

	required := []string{"write", "admin"}
	err := s.authUser(r.Context(), required, host, org, repo, user, pass)
	if err != nil {
		msg := "You must have push access to verify locks"
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#unauthorized-response-2"
		outputError(w, rid, msg, url, http.StatusForbidden)
		return
	}
	log.Printf("authorized user")

	dec := json.NewDecoder(r.Body)
	req := &batchReq{}
	dec.Decode(req)

	// b := s.gs.Bucket(s.bucket)
	// if _, err := b.Attrs(ctx); err != nil {
	// 	msg := fmt.Sprintf("the bucket '%s' doesn't exist", s.bucket)
	// 	outputError(w, rid, msg, "", http.StatusInternalServerError)
	// 	return
	// }
	// o := b.Object(req.Objects.OID)

	out := batchResp{
		Transfer: "basic",
		//Objects: []batchRespObj{},
	}

	objs := []batchRespObj{}

	log.Printf("setting up output")

	for _, o := range req.Objects {
		ea := time.Now().Add(s.lockTimeout)
		url, err := storage.SignedURL(s.bucket, o.OID, &storage.SignedURLOptions{
			Scheme:         storage.SigningSchemeV4,
			Method:         "PUT",
			GoogleAccessID: s.conf.Email,
			PrivateKey:     s.conf.PrivateKey,
			Expires:        ea,
		})

		if err != nil {
			msg := fmt.Sprintf("unable to create signed url for file: %v", err)
			outputError(w, rid, msg, "", http.StatusInternalServerError)
			return
		}

		or := batchRespObj{
			OID:           o.OID,
			Size:          o.Size,
			Authenticated: true,
			Actions: batchRespObjAction{
				Download: &downloadAction{
					HREF:      url,
					ExpiresAt: ea,
				},
			},
		}
		objs = append(objs, or)
	}
	out.Objects = objs

	log.Printf("got output: %v", spew.Sdump(out))

	enc := json.NewEncoder(w)
	enc.Encode(out)
}

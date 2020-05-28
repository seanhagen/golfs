package golfs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
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
	Download *action `json:"download,omitempty"`
	Upload   *action `json:"upload,omitempty"`
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

type action struct {
	HREF      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresIn *int64            `json:"expires_in,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

func (s *service) objectBatch(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	ctx := r.Context()
	rid := ctx.Value("rid").(string)

	log.Printf("request to object batch -- %v", r.URL.Path)

	dec := json.NewDecoder(r.Body)
	req := &batchReq{}
	dec.Decode(req)

	switch req.Operation {
	case "upload":
		log.Printf("object batch -- upload")
		s.objectBatchUpload(ctx, w, pr, req)
	case "download":
		log.Printf("object batch -- download")
		s.objectBatchDownload(ctx, w, pr, req)
	default:
		msg := fmt.Sprintf("invalid git-lfs operation '%v'", req.Operation)
		outputError(w, rid, msg, "", http.StatusBadRequest)
	}
}

// objectBatchUpload ...
func (s *service) objectBatchUpload(ctx context.Context, w http.ResponseWriter, pr httprouter.Params, req *batchReq) {
	rid := ctx.Value("rid").(string)
	user := ctx.Value("auth-user").(string)
	pass := ctx.Value("auth-pass").(string)

	host := pr.ByName("host")
	org := pr.ByName("org")
	repo := pr.ByName("repo")

	required := []string{"write", "admin"}
	err := s.authUser(ctx, required, host, org, repo, user, pass)
	if err != nil {
		msg := "You must have push access to verify locks"
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#unauthorized-response-2"
		outputError(w, rid, msg, url, http.StatusForbidden)
		return
	}

	if s.locking {
		// check for locks!
	}

	out := batchResp{
		Transfer: "basic",
	}

	objs := []batchRespObj{}

	for _, o := range req.Objects {
		ea := time.Now().Add(s.lockTimeout)
		obj := fmt.Sprintf("%v/%v/%v/%v/%v", host, org, repo, req.Ref.Name, o.OID)
		url, err := storage.SignedURL(s.bucket, obj, &storage.SignedURLOptions{
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
				Upload: &action{
					HREF:      url,
					ExpiresAt: ea,
				},
			},
		}
		objs = append(objs, or)
	}
	out.Objects = objs
	enc := json.NewEncoder(w)
	enc.Encode(out)
}

// objectBatchDownload ...
func (s *service) objectBatchDownload(ctx context.Context, w http.ResponseWriter, pr httprouter.Params, req *batchReq) {
	rid := ctx.Value("rid").(string)
	user := ctx.Value("auth-user").(string)
	pass := ctx.Value("auth-pass").(string)

	host := pr.ByName("host")
	org := pr.ByName("org")
	repo := pr.ByName("repo")

	required := []string{"read", "admin"}
	err := s.authUser(ctx, required, host, org, repo, user, pass)
	if err != nil {
		msg := "You must have push access to verify locks"
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#unauthorized-response-2"
		outputError(w, rid, msg, url, http.StatusForbidden)
		return
	}

	if s.locking {
		// check for locks!
	}

	out := batchResp{
		Transfer: "basic",
	}

	objs := []batchRespObj{}

	for _, o := range req.Objects {
		ea := time.Now().Add(s.lockTimeout)
		log.Printf("building download url for %v/%v - %v - %v", org, repo, req.Ref.Name, o.OID)
		obj := fmt.Sprintf("%v/%v/%v/%v/%v", host, org, repo, req.Ref.Name, o.OID)
		url, err := storage.SignedURL(s.bucket, obj, &storage.SignedURLOptions{
			Scheme:         storage.SigningSchemeV4,
			Method:         "GET",
			GoogleAccessID: s.conf.Email,
			PrivateKey:     s.conf.PrivateKey,
			Expires:        ea,
		})
		if err != nil {
			msg := fmt.Sprintf("unable to create signed url for file: %v", err)
			outputError(w, rid, msg, "", http.StatusInternalServerError)
			return
		}
		log.Printf("built url: %v", url)
		or := batchRespObj{
			OID:           o.OID,
			Size:          o.Size,
			Authenticated: true,
			Actions: batchRespObjAction{
				Download: &action{
					HREF:      url,
					ExpiresAt: ea,
				},
			},
		}
		objs = append(objs, or)
	}
	out.Objects = objs
	enc := json.NewEncoder(w)
	enc.Encode(out)
}

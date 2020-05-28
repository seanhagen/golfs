package golfs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const contentType = "application/vnd.git-lfs+json"

const lockEntityType = "LFSEntity"

func GOLFS(w http.ResponseWriter, r *http.Request) {
	log.Printf("request made to: %v", r.URL.Path)

	w.Header().Set("Content-Type", contentType)
	rid := r.Header.Get("X-Cloud-Trace-Context")

	s, err := setupService(r.Context())
	if err != nil {
		msg := fmt.Sprintf("unable to get configuration: %v", err)
		outputError(w, rid, msg, "", http.StatusInternalServerError)
		return
	}

	user, pass, ok := r.BasicAuth()
	if !ok {
		msg := "This route requires HTTP Basic Auth with your GitHub username & password"
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/authentication.md#specified-in-url"
		outputError(w, rid, msg, url, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "auth-user", user)
	ctx = context.WithValue(ctx, "auth-pass", pass)
	ctx = context.WithValue(ctx, "rid", rid)
	s.ServeHTTP(w, r.WithContext(ctx))
}

func outputError(w http.ResponseWriter, rid, msg, url string, status int) {
	log.Printf("error encountered: %v", msg)
	w.WriteHeader(status)

	r := lockError{Message: msg, DocumentationURL: url, RequestId: rid}
	wr := json.NewEncoder(w)
	err := wr.Encode(r)
	if err != nil {
		log.Printf("unable to marshal JSON response: %v", err)
		return
	}
}

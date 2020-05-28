package golfs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
)

// verifyLocks handles the /locks/verify request -- see the Git-LFS docs: https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#list-locks-for-verification
func (s *service) verifyLocks(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	ctx := r.Context()
	rid := ctx.Value("rid").(string)

	if !s.locking {
		// not locking, return 404
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#not-found-response"
		outputError(w, rid, "Locking disabled", url, http.StatusNotFound)
		return
	}

	// test because of httprouter routing thing
	if t := pr.ByName("id"); t != "verify" {
		msg := fmt.Sprintf("Invalid route: %v", t)
		url := "https://github.com/git-lfs/git-lfs/blob/master/docs/api/locking.md#list-locks-for-verification"
		outputError(w, rid, msg, url, http.StatusForbidden)
		return
	}

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

	dec := json.NewDecoder(r.Body)
	req := &lockVerifyRequest{}
	if err = dec.Decode(req); err != nil {
		outputError(w, rid, "unable to unmarshal request body", "", http.StatusInternalServerError)
		return
	}

	q := datastore.NewQuery(lockEntityType).
		Namespace(s.namespace).
		Filter("Org =", org).
		Filter("Repo =", repo).
		Filter("RefName =", req.Ref.Name)

	out := []storedLock{}
	_, err = s.ds.GetAll(r.Context(), q, &out)
	if err != nil {
		outputError(w, rid, fmt.Sprintf("unable to query datastore: %v", err), "", http.StatusInternalServerError)
		return
	}

	ours := []lock{}
	theirs := []lock{}

	for _, l := range out {
		if l.Owner == user {
			ours = append(ours, lock{
				ID:       l.UUID,
				Path:     l.Path,
				LockedAt: l.LockedAt,
				Owner:    name{Name: l.Owner},
			})
		} else {
			theirs = append(theirs, lock{
				ID:       l.UUID,
				Path:     l.Path,
				LockedAt: l.LockedAt,
				Owner:    name{Name: l.Owner},
			})
		}
	}

	output := verifyLockResponse{
		Ours:   ours,
		Theirs: theirs,
	}

	wr := json.NewEncoder(w)
	wr.Encode(output)
}

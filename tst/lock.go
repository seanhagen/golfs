package gitlfstest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

/*
{
  "lock": {
    "id": "some-uuid",
    "path": "foo/bar.zip",
    "locked_at": "2016-05-17T15:49:06+00:00",
    "owner": {
      "name": "Jane Doe"
    }
  }
}
*/

type LockResponse struct {
	Lock struct {
		ID       uuid.UUID `json:"id"`
		Path     string    `json:"path"`
		LockedAt time.Time `json:"locked_at"`
		Owner    struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"lock"`
}

/*
{
  "path": "foo/bar.zip",
  "ref": {
    "name": "refs/heads/my-feature"
  }
}
*/

type LockRequest struct {
	Path string `json:"path"`
	Ref  struct {
		Name string `json:"name"`
	} `json:"ref"`
}

func VerifyLock(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Printf("lock reequest input: %v -- %v", string(body), spew.Sdump(pr))

	if err != nil {
		log.Printf("unable to read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := &LockRequest{}
	if err = json.Unmarshal(body, res); err != nil {
		log.Printf("unable to unmarshal request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("got lock request: %v", spew.Sdump(res))

	w.WriteHeader(http.StatusInternalServerError)
}

package golfs

import (
	"time"
)

type name struct {
	Name string `json:"name"`
}

type lockRequest struct {
	Path string `json:"path"`
	Ref  name   `json:"ref"`
}

type lock struct {
	ID       string    `json:"id"`
	Path     string    `json:"path"`
	LockedAt time.Time `json:"locked_at"`
	Owner    name      `json:"owner"`
}

type lockResponse struct {
	Lock lock `json:"lock"`
}

type lockVerifyRequest struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int64  `json:"limit,omitempty"`
	Ref    name    `json:"ref"`
}

type lockError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	RequestId        string `json:"request_id"`
}

type verifyLockResponse struct {
	Ours       []lock `json:"ours"`
	Theirs     []lock `json:"theirs"`
	NextCursor string `json:"next_cursor"`
}

/*
{
  "ours": [
    {
      "id": "some-uuid",
      "path": "/path/to/file",
      "locked_at": "2016-05-17T15:49:06+00:00",
      "owner": {
        "name": "Jane Doe"
      }
    }
  ],
  "theirs": [],
  "next_cursor": "optional next ID"
}
*/

type storedLock struct {
	UUID string

	Org  string
	Repo string

	Path     string
	LockedAt time.Time
	Owner    string
}

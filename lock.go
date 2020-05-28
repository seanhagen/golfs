package golfs

import (
	"time"

	"github.com/gofrs/uuid"
)

type lockRequest struct {
	Path string `json:"path"`
	Ref  struct {
		Name string `json:"name"`
	} `json:"ref"`
}

type lockResponse struct {
	Lock struct {
		ID       uuid.UUID `json:"id"`
		Path     string    `json:"path"`
		LockedAt time.Time `json:"locked_at"`
		Owner    struct {
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"lock"`
}

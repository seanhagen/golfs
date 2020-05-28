package golfs

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

func (s *service) authUser(host, org, repo, user, pass string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	spew.Dump(client)
}

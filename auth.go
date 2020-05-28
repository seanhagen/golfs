package golfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

// authUser ...
func (s *service) authUser(ctx context.Context, require []string, host, org, repo, user, pass string) error {
	if host != "github.com" {
		return fmt.Errorf("only GitHub supported currently")
	}
	return s.authGithubUser(ctx, require, org, repo, user, pass)
}

// this token requires the 'repo' scope
func (s *service) authGithubUser(ctx context.Context, require []string, org, repo, user, pass string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pass},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	perm, _, err := client.Repositories.GetPermissionLevel(ctx, org, repo, user)
	if err != nil {
		return err
	}

	if perm.Permission == nil {
		return fmt.Errorf("user has no permissions")
	}

	for _, v := range require {
		if v == *perm.Permission {
			return nil
		}
	}

	return fmt.Errorf("missing required scope -- any of: %v", strings.Join(require, ","))
}

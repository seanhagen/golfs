package golfs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var ErrNoNamespace = fmt.Errorf("a Google Cloud Datastore namespace is required")
var ErrNoGithubToken = fmt.Errorf("a GitHub API token is required")

const defaultTimeout = time.Minute * 5
const defaultLocking = true

type service struct {
	// google datastore namespace
	namespace string
	// is locking enabled
	locking bool
	// how long do locks last when enabled
	lockTimeout time.Duration
	// the token used to authorize calls to the GitHub API
	githubToken string
}

func setupService() (*service, error) {
	viper.SetEnvPrefix("GOLFS")
	viper.SetDefault("lock_timeout", defaultTimeout.String())
	viper.SetDefault("locking", defaultLocking)
	viper.AutomaticEnv()

	ns := viper.GetString("ds_namespace")
	if ns == "" {
		return nil, ErrNoNamespace
	}

	gt := viper.GetString("github_token")
	if gt == "" {
		return nil, ErrNoGithubToken
	}

	dur := defaultTimeout
	if to := viper.GetString("lock_timeout"); to != "" {
		dt, err := time.ParseDuration(to)
		if err != nil {
			return nil, err
		}
		dur = dt
	}

	c := &service{
		namespace:   ns,
		githubToken: gt,

		lockTimeout: dur,
		locking:     viper.GetBool("locking"),
	}

	return c, nil
}

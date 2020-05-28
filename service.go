package golfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

var ErrNoNamespace = fmt.Errorf("a Google Cloud Datastore namespace is required")
var ErrNoBucket = fmt.Errorf("a Google Cloud Storage bucket name is required")

const defaultTimeout = time.Minute * 5
const defaultLocking = true

type service struct {
	// google datastore namespace
	namespace string
	// google cloud storage bucket
	bucket string

	// is locking enabled
	locking bool
	// how long do locks last when enabled
	lockTimeout time.Duration

	router *httprouter.Router

	conf *jwt.Config

	ds *datastore.Client

	gs *storage.Client
}

func setupService(ctx context.Context) (*service, error) {
	viper.SetEnvPrefix("GOLFS")
	viper.SetDefault("lock_timeout", defaultTimeout.String())
	viper.SetDefault("locking", defaultLocking)
	viper.AutomaticEnv()

	ns := viper.GetString("ds_namespace")
	if ns == "" {
		return nil, ErrNoNamespace
	}

	bk := viper.GetString("bucket")
	if bk == "" {
		return nil, ErrNoBucket
	}

	dur := defaultTimeout
	if to := viper.GetString("lock_timeout"); to != "" {
		dt, err := time.ParseDuration(to)
		if err != nil {
			return nil, err
		}
		dur = dt
	}

	prj := os.Getenv("GCP_PROJECT")
	dsClient, err := datastore.NewClient(ctx, prj)
	if err != nil {
		log.Printf("unable to create datastore client, project: '%v'", prj)
		return nil, err
	}

	gsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("unable to create storage client, project: '%v'", prj)
		return nil, err
	}

	files, err := ioutil.ReadDir("./serverless_function_source_code")
	if err != nil {
		log.Printf("ioutil.ListFiles: %v", err)
	}
	log.Printf("files in './serverless_function_source_code")
	for _, f := range files {
		log.Printf("\t%v", f.Name())
	}
	log.Printf("\n\n")

	path, err := filepath.Abs("./serverless_function_source_code/account.json")
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path to account.json: %v", err)
	}

	jsonKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	s := &service{
		namespace: ns,
		bucket:    bk,

		lockTimeout: dur,
		locking:     viper.GetBool("locking"),

		conf: conf,

		ds: dsClient,
		gs: gsClient,
	}

	r := httprouter.New()

	// object send/download
	r.POST("/:host/:org/:repo/objects/batch", s.objectBatch)

	//locks
	r.POST("/:host/:org/:repo/locks", s.createLock)
	r.GET("/:host/:org/:repo/locks", s.fetchLocks)
	r.POST("/:host/:org/:repo/locks/:id/unlock", s.deleteLock)
	r.POST("/:host/:org/:repo/locks/:id", s.verifyLocks)

	s.router = r

	return s, nil
}

// ServeHTTP ...
func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

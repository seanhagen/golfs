package gitlfstest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
{
	"operation":"upload",
	"objects":[
		{"oid":"12993ba6227990ece4e77a149b1bc6edb462b2744cbeaf3b1964a33b0ca50b23","size":543262}
	],
	"transfers":["lfs-standalone-file","basic"],
	"ref":{"name":"refs/heads/setup"}
}
*/

func ObjectBatch(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	user, pass, ok := r.BasicAuth()
	log.Printf("auth(%v): %v -> %v", ok, user, pass)
	log.Printf("url: %v - %v - %v", r.URL.Host, r.URL.Path, r.URL.String())
	//log.Printf("headers: %v", spew.Sdump(r.Header))
	log.Printf("received: %v", string(data))

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "woop")

}

package main

import (
	"fmt"
	"net/http"
	"os"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "vcap_services is %s", os.Getenv("VCAP_SERVICES"))
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

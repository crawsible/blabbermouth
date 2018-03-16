package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
)

const credhubURI string = "https://credhub.service.cf.internal:8844"

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "vcap_services is %s\n", os.Getenv("VCAP_SERVICES"))

	ch, err := credhubClient()
	if err != nil {
		panic(err)
	}

	paths, err := ch.FindAllPaths()
	if err != nil {
		fmt.Fprintf(w, "unable to talk to CredHub\n")
		return
	}

	pathsBytes, err := json.MarshalIndent(paths, "", "\t")
	if err != nil {
		fmt.Fprintf(w, "got odd response from CredHub\n")
	}

	fmt.Fprintf(w, "credhub paths I have access to:\n")
	fmt.Fprintf(w, string(pathsBytes)+"\n")
}

func credhubClient() (*credhub.CredHub, error) {
	if os.Getenv("CF_INSTANCE_CERT") == "" || os.Getenv("CF_INSTANCE_KEY") == "" {
		return nil, fmt.Errorf("Missing CF_INSTANCE_CERT and/or CF_INSTANCE_KEY")
	}
	if os.Getenv("CF_SYSTEM_CERT_PATH") == "" {
		return nil, fmt.Errorf("Missing CF_SYSTEM_CERT_PATH")
	}

	systemCertsPath := os.Getenv("CF_SYSTEM_CERT_PATH")
	caCerts := []string{}
	files, err := ioutil.ReadDir(systemCertsPath)
	if err != nil {
		return nil, fmt.Errorf("Can't read contents of system cert path: %v", err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".crt") {
			contents, err := ioutil.ReadFile(filepath.Join(systemCertsPath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("Can't read contents of cert in system cert path: %v", err)
			}
			caCerts = append(caCerts, string(contents))
		}
	}

	return credhub.New(
		credhubURI,
		credhub.ClientCert(os.Getenv("CF_INSTANCE_CERT"), os.Getenv("CF_INSTANCE_KEY")),
		credhub.CaCerts(caCerts...),
	)
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

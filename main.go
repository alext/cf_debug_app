package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(handler))
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	env := newEnvMap(os.Environ())
	data := make(map[string]interface{})
	data["ENV"] = env
	data["es_stats"] = elasticsearch_status(env)

	output, err := json.MarshalIndent(data, "  ", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

type service struct {
	Name        string
	Label       string
	Tags        []string
	Plan        string
	Credentials struct {
		Hostname string
		Ports    map[string]string
	}
}

func elasticsearch_status(env envMap) interface{} {
	serviceInfoJson, ok := env["VCAP_SERVICES"]
	if !ok {
		return "Missing VCAP_SERVICES in environment"
	}

	serviceInfo := make(map[string][]service)
	err := json.Unmarshal([]byte(serviceInfoJson), &serviceInfo)
	if err != nil {
		return "Error parsing VCAP_SERVICES Json : " + err.Error()
	}

	esServicesInfo, ok := serviceInfo["elasticsearch13"]
	if !ok || len(esServicesInfo) < 1 {
		return "No elasticsearch services available"
	}
	esServiceInfo := esServicesInfo[0]

	esURL := fmt.Sprintf("http://%s:%s", esServiceInfo.Credentials.Hostname, esServiceInfo.Credentials.Ports["9200/tcp"])

	resp, err := http.Get(esURL + "/_cluster/stats")
	if err != nil {
		return "Error querying ES status : " + err.Error()
	}
	defer resp.Body.Close()
	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "Error parsing ES Json : " + err.Error()
	}
	return data
}

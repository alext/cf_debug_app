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
	data["vcap_application"] = extractJSONData(env, "VCAP_APPLICATION")
	data["vcap_services"] = extractJSONData(env, "VCAP_SERVICES")

	output, err := json.MarshalIndent(data, "  ", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func extractJSONData(env envMap, key string) interface{} {
	raw, ok := env[key]
	if !ok {
		return fmt.Sprintf("Env var %s not found", key)
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON : %s", err.Error())
	}
	return data
}

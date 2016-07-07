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
	data["ENV"] = expandJSONFields(env)

	output, err := json.MarshalIndent(data, "  ", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}

func expandJSONFields(env envMap) map[string]interface{} {
	data := make(map[string]interface{}, len(env))
	for key, value := range env {
		if value[0] == '{' || value[0] == '[' {
			var valueData interface{}
			err := json.Unmarshal([]byte(value), &valueData)
			if err != nil {
				data[key] = value
			} else {
				data[key] = valueData
			}
		} else {
			data[key] = value
		}
	}
	return data
}

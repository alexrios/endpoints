package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ConfigFile struct {
	Addr      string `json:"address"`
	Responses []Response
}
type Response struct {
	Status   int
	Path     string
	JsonBody string
}

func main() {
	file, err := ioutil.ReadFile("endpoints.json")
	if err != nil {
		log.Fatal(err)
	}
	configFile := ConfigFile{}
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, response := range configFile.Responses {
		closureResponse := response
		fmt.Println(closureResponse.Path, "will return status", closureResponse.Status, "with body from file ", closureResponse.JsonBody)
		http.HandleFunc(closureResponse.Path, func(writer http.ResponseWriter, request *http.Request) {
			if 200 != closureResponse.Status {
				writer.WriteHeader(closureResponse.Status)
			}
			writer.Header().Add("Content-Type", "application/json")
			jsonBytes, err := ioutil.ReadFile(closureResponse.JsonBody)
			if err != nil {
				log.Fatal(err)
			}
			writer.Write(jsonBytes)
		})
	}
	log.Fatal(http.ListenAndServe(configFile.Addr, nil))
}

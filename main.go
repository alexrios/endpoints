package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
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
	router := mux.NewRouter()

	for _, response := range configFile.Responses {
		closureResponse := response
		fmt.Println(closureResponse.Path, "will return status", closureResponse.Status, "with body from file ", closureResponse.JsonBody)
		router.HandleFunc(closureResponse.Path, func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			if 200 != closureResponse.Status {
				writer.WriteHeader(closureResponse.Status)
			}
			writer.Header().Add("Content-Type", "application/json")
			jsonBytes, err := ioutil.ReadFile(closureResponse.JsonBody)
			if err != nil {
				log.Fatal(err)
			}

			t := template.Must(template.New("letter").Parse(string(jsonBytes)))
			_ = t.Execute(writer, vars)
		})
	}
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(configFile.Addr, nil))
}

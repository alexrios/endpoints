package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

type ConfigFile struct {
	Addr      string     `json:"address"`
	Responses []Response `json:"responses"`
}
type Response struct {
	Status   int    `json:"status"`
	Path     string `json:"path"`
	Latency  string `json:"latency"`
	JsonBody string `json:"json_body"`
}

func main() {
	file, err := ioutil.ReadFile("endpoints.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	configFile := ConfigFile{}
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := mux.NewRouter()
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	router.Use(cors)

	for _, response := range configFile.Responses {
		closure := response
		log.Println(closure.Path, "will return status", closure.Status, "with body from file", closure.JsonBody)
		router.HandleFunc(closure.Path, func(writer http.ResponseWriter, request *http.Request) {
			if closure.Status != 200 {
				writer.WriteHeader(closure.Status)
			}
			writer.Header().Add("Content-Type", "application/json")

			tmpl, err := template.ParseFiles(closure.JsonBody)
			if err != nil {
				log.Error(err.Error())
				writer.WriteHeader(500)
				return
			}
			if closure.Latency != "" {
				duration, err := time.ParseDuration(closure.Latency)
				if err != nil {
					log.Error(err.Error())
				}
				time.Sleep(duration)
			}

			err = tmpl.Execute(writer, mux.Vars(request))
			if err != nil {
				log.Error(err.Error())
			}
		})
	}
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(configFile.Addr, nil))
}

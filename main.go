package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

type ConfigFile struct {
	Addr      string     `json:"address"`
	Responses []Response `json:"responses"`
}
type Response struct {
	Status   int    `json:"status"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Latency  string `json:"latency"`
	JsonBody string `json:"json_body"`
}

func main() {
	if !fileExists("endpoints.json") {
		err := ioutil.WriteFile("endpoints.json", []byte(DefaultFile), 0644)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("customBody.json", []byte(CustomBody), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

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
		if closure.Status == 0 {
			closure.Status = 200
		}
		if closure.Method == "" {
			closure.Method = "GET"
		}

		log.Info("[", closure.Method, "] ", closure.Path, " -> ", closure.Status, " with body -> ", closure.JsonBody)
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
		}).Methods(closure.Method)
	}
	http.Handle("/", router)
	if configFile.Addr == "" {
		configFile.Addr = ":8080"
	}
	log.Info("Listen at ", configFile.Addr)
	log.Fatal(http.ListenAndServe(configFile.Addr, nil))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

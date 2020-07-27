package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

func newHandleFunc(fs afero.Fs, r Response) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if r.Status != 200 {
			writer.WriteHeader(r.Status)
		}

		tmpl, err := parseFile(fs, r.JsonBody)
		if err != nil {
			log.Error(err)
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		if r.Latency != "" {
			duration, err := time.ParseDuration(r.Latency)
			if err != nil {
				log.Error(err)
			}
			time.Sleep(duration)
		}

		writer.Header().Add("Content-Type", "application/json")
		err = tmpl.Execute(writer, mux.Vars(request))
		if err != nil {
			log.Error(err)
		}
	}
}

func parseFile(fs afero.Fs, filename string) (*template.Template, error) {
	b, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	s := string(b)
	name := filepath.Base(filename)
	tmpl := template.New(name)
	_, err = tmpl.Parse(s)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

package main

import (
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func newHandleFunc(fs afero.Fs, path string, r Response) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if r.Status != 200 {
			writer.WriteHeader(r.Status)
		}

		tmpl, err := parseFile(fs, path, r.JsonBody)
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

func parseFile(fs afero.Fs, path string, filename string) (*template.Template, error) {
	if path != "" {
		filename = filepath.Join(path, filename)
	}
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

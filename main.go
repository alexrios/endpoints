package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const ShutdownTimeout = "30s"

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "location of config file (default: .)")
	flag.Parse()

	logger := log.New()
	var appFs = afero.NewOsFs()
	var cfgFile string
	if cfgPath == "" {
		cfgFile = DefaultConfigurationFileName
	} else {
		cfgFile = filepath.Join(cfgPath, DefaultConfigurationFileName)
	}
	if err := run(appFs, cfgPath, cfgFile, logger); err != nil {
		log.WithField("fn", "main").Error(err)
		os.Exit(1)
	}
}

func run(fs afero.Fs, path string, cfgFile string, log *log.Logger) error {
	firstRun, err := isFirstRun(fs, cfgFile)
	if err != nil {
		return err
	}
	if firstRun {
		err = configureFirstRun(fs, cfgFile)
		if err != nil {
			return err
		}
	}

	configFile, err := loadConfig(fs, cfgFile)
	if err != nil {
		return err
	}

	//Config CORS
	router := mux.NewRouter()
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)
	router.Use(cors)
	//Set handlers
	for _, response := range configFile.Responses {
		closure := response
		log.Info("[", closure.Method, "] ", closure.Path, " -> ", closure.Status, " with body -> ", closure.JsonBody)
		router.HandleFunc(closure.Path, newHandleFunc(fs, path, closure)).Methods(closure.Method)
	}
	http.Handle("/", router)

	server := http.Server{
		Addr:    configFile.Addr,
		Handler: router,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API listening on %s\n", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatal(err)

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown\n", sig)

		// Give outstanding requests a deadline for completion.
		duration, err := time.ParseDuration(ShutdownTimeout)
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	log.Fatal(server.ListenAndServe())
	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("starting...")

	r := mux.NewRouter()
	s := &server{}
	s.registerRoute(r)
	fmt.Print("routes registered...")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Print("serving...")
	log.Fatal(srv.ListenAndServe())
}

type server struct {
}

func (s *server) registerRoute(router *mux.Router) {
	router.HandleFunc("/{module:(?:[^/]+/)+[^/]+}/@v/list", s.handleList)
}

func (s *server) handleList(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	module := vars["module"]
	username := req.Header.Get("username")
	password := req.Header.Get("password")
	output, err := s.runFetcherImage(req.Context(), "--module", module, "--username", username, "--password", password, "--host", "github.com")
	if err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.Write([]byte(output))
}

func (s *server) runFetcherImage(ctx context.Context, args ...string) (string, error) {
	cm := []string{"run", "fetcher"}
	cm = append(cm, args...)
	cmd := exec.CommandContext(ctx, "docker", cm...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execute docker command has errors: %v, %s", err, output)
	}
	return string(output), nil
}

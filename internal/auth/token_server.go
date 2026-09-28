package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type authHandler struct {
	code     string
	exitChan chan struct{}
}

func GetAuthCode() (string, error) {
	ah := authHandler{
		exitChan: make(chan struct{}),
	}
	port := "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", ah.handleGetCode)

	tokenServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Getting auth code on port %s", port)
		err := tokenServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return "", fmt.Errorf(
			"failed to start capture server: %w", err,
		)
	case <-ah.exitChan:
		log.Println("Received shutdown trigger")
	}
	// <-ah.exitChan
	err := tokenServer.Shutdown(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to capture authorization code: %w", err)
	}

	if ah.code == "" {
		return "", fmt.Errorf("no authorization code received")
	}

	return ah.code, nil
}

func (a *authHandler) handleGetCode(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(
			w,
			"Authorization rejected: "+errMsg,
			http.StatusBadRequest,
		)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(
			w,
			"Missing authorization code",
			http.StatusBadRequest,
		)
		return
	}

	a.code = code
	http.ServeFile(w, r, "./internal/auth/resources/redirect.html")

	go func() {
		a.exitChan <- struct{}{}
	}()
}

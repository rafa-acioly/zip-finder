package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rafa-acioly/zip-finder/core"
	"github.com/rafa-acioly/zip-finder/services"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	ctx, cancelCtx := context.WithTimeout(r.Context(), time.Second*1)
	defer cancelCtx()

	zipCodeServices := []core.Service{
		&services.ViaCep{},
		&services.PostMon{},
	}
	channel := make(chan core.ServiceResponse, len(zipCodeServices))

	baseService := services.BaseService{}
	for _, service := range zipCodeServices {
		go baseService.DoRequest("07400885", service, channel)
	}

	select {
	case response := <-channel:
		json.NewEncoder(w).Encode(response)
		return
	case <-ctx.Done():
		http.Error(
			w,
			http.StatusText(http.StatusNoContent),
			http.StatusNoContent,
		)
	}
}

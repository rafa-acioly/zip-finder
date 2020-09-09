package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rafa-acioly/zip-finder/core"
	"github.com/rafa-acioly/zip-finder/services"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{zipcode}", handler).Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))

	fmt.Println("Listening on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancelCtx := context.WithTimeout(r.Context(), time.Second*1)
	defer cancelCtx()

	zipCodeServices := []core.Service{
		&services.ViaCep{},
		&services.PostMon{},
		&services.RepublicaVirtual{},
	}

	channel := make(chan core.ServiceResponse, len(zipCodeServices))
	vars := mux.Vars(r)
	baseService := services.BaseService{}
	for _, service := range zipCodeServices {
		go baseService.DoRequest(vars["zipcode"], service, channel)
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

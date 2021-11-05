package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rafa-acioly/zip-finder/api/core"
	"github.com/rafa-acioly/zip-finder/api/services"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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

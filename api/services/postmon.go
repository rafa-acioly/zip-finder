package services

import (
	"encoding/json"
	"fmt"

	"github.com/rafa-acioly/zip-finder/api/core"
)

const (
	postmonURL       = "https://api.postmon.com.br/v1/cep/%s"
	postmonServiceID = "postmon"
)

type PostMon struct {
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Logradouro string `json:"logradouro"`
	Estado     string `json:"estado"`
}

func (p PostMon) ValidResponse(serviceResponseContent []byte) bool {
	if len(serviceResponseContent) == 0 {
		return false
	}

	return true
}

func (p PostMon) GetServiceEndpoint(zipCode string) string {
	return fmt.Sprintf(postmonURL, zipCode)
}

func (p *PostMon) ParseServiceResponse(serviceResponseContent []byte) error {
	err := json.Unmarshal(serviceResponseContent, p)
	if err != nil {
		return err
	}

	return nil
}

func (p PostMon) ConvertToDefaultResponse() core.ServiceResponse {
	return core.ServiceResponse{
		Service:    postmonServiceID,
		Cidade:     p.Cidade,
		Bairro:     p.Bairro,
		Logradouro: p.Logradouro,
		UF:         p.Estado,
	}
}

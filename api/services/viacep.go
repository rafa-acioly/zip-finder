package services

import (
	"encoding/json"
	"fmt"

	"github.com/rafa-acioly/zip-finder/api/core"
)

const (
	viacepURL       = "https://viacep.com.br/ws/%s/json/"
	viacepServiceID = "viacep"
)

// ViaCep represents the basic information structure retrieved by viacep API
type ViaCep struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

// ValidResponse verify if the key "erro" is present on the service response,
// the viacep endpoint will always return status "200"
func (v ViaCep) ValidResponse(content []byte) bool {
	response := make(map[string]interface{})
	_ = json.Unmarshal(content, &response)

	if _, ok := response["erro"]; ok {
		return false
	}

	return true
}

func (v ViaCep) GetServiceEndpoint(zipCode string) string {
	url := fmt.Sprintf(viacepURL, zipCode)
	return url
}

func (v *ViaCep) ParseServiceResponse(serviceResponseContent []byte) error {
	return json.Unmarshal(serviceResponseContent, v)
}

func (v ViaCep) ConvertToDefaultResponse() core.ServiceResponse {
	return core.ServiceResponse{
		Service:    viacepServiceID,
		Cidade:     v.Localidade,
		Bairro:     v.Bairro,
		Logradouro: v.Logradouro,
		UF:         v.UF,
	}
}

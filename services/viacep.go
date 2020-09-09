package services

import (
	"encoding/json"
	"fmt"

	"github.com/rafa-acioly/zip-finder/core"
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
// this validation is needed because viacep retrieves status "200" even when
// the zipcode does not exist or is invalid
func (v ViaCep) ValidResponse(content []byte) bool {
	response := make(map[string]interface{})
	_ = json.Unmarshal(content, &response)

	if _, ok := response["erro"]; ok {
		return true
	}

	return false
}

// GetServiceEndpoint retrieves the service full URL that should be called for viacep service
func (v ViaCep) GetServiceEndpoint(zipCode string) string {
	url := fmt.Sprintf(viacepURL, zipCode)
	return url
}

// ParseServiceResponse will convert the service response into a default response
func (v *ViaCep) ParseServiceResponse(responseContent []byte) error {
	err := json.Unmarshal(responseContent, v)
	if err != nil {
		return err
	}

	return nil
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

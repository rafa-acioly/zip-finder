package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rafa-acioly/zip-finder/core"
)

const (
	republicaVirtualURL = "https://republicavirtual.com.br/web_cep.php?cep=%s&formato=json"
	republicaVirtualID  = "republica_virtual"
)

type RepublicaVirtual struct {
	Bairro         string `json:"bairro"`
	Cidade         string `json:"cidade"`
	Logradouro     string `json:"logradouro"`
	TipoLogradouro string `json:"tipo_logradouro"`
	UF             string `json:"uf"`
}

// ValidResponse verify if the key "resultado" is not "0", if it is "0" that means
// that the zipcode is invalid or doesn't exist, the republicaVirtual endpoint
// will always return status "200"
func (r RepublicaVirtual) ValidResponse(serviceResponseContent []byte) bool {
	response := make(map[string]interface{})
	_ = json.Unmarshal(serviceResponseContent, &response)

	if response["resultado"] == "0" {
		return false
	}

	return true
}

func (r RepublicaVirtual) GetServiceEndpoint(zipCode string) string {
	return fmt.Sprintf(republicaVirtualURL, zipCode)
}

func (r *RepublicaVirtual) ParseServiceResponse(serviceResponseContent []byte) error {
	return json.Unmarshal(serviceResponseContent, r)
}

func (r RepublicaVirtual) ConvertToDefaultResponse() core.ServiceResponse {
	logradouro := strings.Trim(fmt.Sprintf("%s %s", r.TipoLogradouro, r.Logradouro), " ")
	return core.ServiceResponse{
		Service:    republicaVirtualID,
		Cidade:     r.Cidade,
		Bairro:     r.Bairro,
		Logradouro: logradouro,
		UF:         r.UF,
	}
}

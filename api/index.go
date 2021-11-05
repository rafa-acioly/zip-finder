package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	// MinZipCodeLength is the minimum number of number characters in a zipcode
	MinZipCodeLength = 8
	// ZipCodeRegex is the regex used to validate a zipcode
	ZipCodeRegex = `[^0-9]`

	postmonURL       = "https://api.postmon.com.br/v1/cep/%s"
	postmonServiceID = "postmon"

	republicaVirtualURL = "https://republicavirtual.com.br/web_cep.php?cep=%s&formato=json"
	republicaVirtualID  = "republica_virtual"

	viacepURL       = "https://viacep.com.br/ws/%s/json/"
	viacepServiceID = "viacep"
)

// Service defines the contract that every zip-code service should implement
type Service interface {
	ValidResponse(serviceResponseContent []byte) bool
	GetServiceEndpoint(zipCode string) string
	ParseServiceResponse(serviceResponseContent []byte) error
	ConvertToDefaultResponse() ServiceResponse
}

// BaseService contains the basic methods to call a zipcode service
type BaseService struct{}

// DoRequest ...
func (base BaseService) DoRequest(zipCode string, service Service, channel chan ServiceResponse) {

	response, err := http.Get(service.GetServiceEndpoint(zipCode))
	if err != nil {
		return
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil || !service.ValidResponse(content) {
		return
	}

	err = service.ParseServiceResponse(content)
	if err != nil {
		return
	}

	defaultResponse := service.ConvertToDefaultResponse()

	channel <- defaultResponse
}

// ValidZipCode retrieves if a zipcode contains eight characters and only numbers
func (base BaseService) ValidZipCode(zipCode string) bool {
	re := regexp.MustCompile(ZipCodeRegex)
	formattedZipCode := re.ReplaceAllString(zipCode, `$1`)

	if len(formattedZipCode) < MinZipCodeLength {
		return false
	}

	return true
}

type ServiceResponse struct {
	Service    string `json:"service"`
	Cidade     string `json:"cidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"uf"`
}

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

func (o *PostMon) ConvertToDefaultResponse() ServiceResponse {
	return ServiceResponse{
		Service:    postmonServiceID,
		Cidade:     o.Cidade,
		Bairro:     o.Bairro,
		Logradouro: o.Logradouro,
		UF:         o.Estado,
	}
}

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

func (r RepublicaVirtual) ConvertToDefaultResponse() ServiceResponse {
	logradouro := strings.Trim(fmt.Sprintf("%s %s", r.TipoLogradouro, r.Logradouro), " ")
	return ServiceResponse{
		Service:    republicaVirtualID,
		Cidade:     r.Cidade,
		Bairro:     r.Bairro,
		Logradouro: logradouro,
		UF:         r.UF,
	}
}

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

func (v ViaCep) ConvertToDefaultResponse() ServiceResponse {
	return ServiceResponse{
		Service:    viacepServiceID,
		Cidade:     v.Localidade,
		Bairro:     v.Bairro,
		Logradouro: v.Logradouro,
		UF:         v.UF,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancelCtx := context.WithTimeout(r.Context(), time.Second*1)
	defer cancelCtx()

	zipCodeServices := []Service{
		&ViaCep{},
		&PostMon{},
		&RepublicaVirtual{},
	}

	channel := make(chan ServiceResponse, len(zipCodeServices))
	vars := mux.Vars(r)
	baseService := BaseService{}
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

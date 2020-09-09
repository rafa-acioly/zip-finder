package core

// ServiceResponse is the default payload used to respond on the http calls
type ServiceResponse struct {
	Service    string `json:"service"`
	Cidade     string `json:"cidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	UF         string `json:"uf"`
}

// Service defines the contract that every zip-code service should implement
type Service interface {
	ValidResponse(serviceResponseContent []byte) bool
	GetServiceEndpoint(zipCode string) string
	ParseServiceResponse(serviceResponseContent []byte) error
	ConvertToDefaultResponse() ServiceResponse
}

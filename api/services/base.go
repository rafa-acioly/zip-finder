package services

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/rafa-acioly/zip-finder/api/core"
)

const (
	// MinZipCodeLength is the minimum number of number characters in a zipcode
	MinZipCodeLength = 8
	// ZipCodeRegex is the regex used to validate a zipcode
	ZipCodeRegex = `[^0-9]`
)

// BaseService contains the basic methods to call a zipcode service
type BaseService struct{}

// DoRequest ...
func (base BaseService) DoRequest(zipCode string, service core.Service, channel chan core.ServiceResponse) {

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

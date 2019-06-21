package typeahead

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/emicklei/go-restful"
)

// Provider is a simple struct to hold Provider name and details
type Provider struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IndexName   string `json:"index"`
	Logo        string `json:"logo"`
}

// TODO add in in typeahead functions like parameters, measurements, ???

// New builds out the services calls for type ahead
func New() *restful.WebService {
	service := new(restful.WebService)

	service.
		Path("/api/v1/typeahead").
		Doc("Typeahead services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/providers").To(ProvidersCall).
		Doc("Get list of providers for typeahead query to support query box typeahead").
		ReturnsError(400, "Unable to handle request", nil).
		Writes(Provider{}).
		Operation("ProvidersCall"))

	return service
}

// ProvidersCall returns the provider json package
func ProvidersCall(request *restful.Request, response *restful.Response) {
	ic, err := os.Open("./indexcatalog.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ic.Close()

	pa := []Provider{}

	jsonParser := json.NewDecoder(ic)
	jsonParser.Decode(&pa)

	data, err := json.Marshal(pa)
	if err != nil {
		fmt.Println(err.Error())
	}
	response.Write(data)
}

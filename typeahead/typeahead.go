package typeahead

import (
	"encoding/json"

	"github.com/emicklei/go-restful"
)

// Provider is a simple struct to hold Provider name and details
type Provider struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IndexName   string `json:"index"`
}

// New builds out the services calls for type ahead
func New() *restful.WebService {
	service := new(restful.WebService)

	service.
		Path("/api/v1/typeahead").
		Doc("Typeahead services to support user interfaces").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/providers").To(ProvidersCall).
		Doc("Get list of providers for typeahead query to support query box typeahead").
		ReturnsError(400, "Unable to handle request", nil).
		Writes(Provider{}).
		Operation("ProvidersCall"))
	return service
}

// TODO add in in typeahead functions like parameters, measurements, ???

// TODO  this is just a hard coded place holder for now...  replace later with
// something pulling from KV store or something like that...

// ProvidersCall returns the provider json package
func ProvidersCall(request *restful.Request, response *restful.Response) {

	pa := []Provider{}

	ocd := Provider{Name: "OpenCore", Description: "Open Core Data", IndexName: "ocd"}
	pa = append(pa, ocd)

	bcodmo := Provider{Name: "BCO-DMO", Description: "Biological and Chemical Oceanography Data Management Office", IndexName: "bcodmo"}
	pa = append(pa, bcodmo)

	le := Provider{Name: "LinkedEarth", Description: "EARTHCUBE Linked Earth Building Block", IndexName: "linkedearth"}
	pa = append(pa, le)

	// neotoma := Provider{Name: "Neotoma", Description: "Neotoma", IndexName: "neotoma"}
	// pa = append(pa, neotoma)

	ieda := Provider{Name: "IEDA", Description: "Interdisciplinary Earth Data Alliance ", IndexName: "ieda"}
	pa = append(pa, ieda)

	rwg := Provider{Name: "EarthCube RWG", Description: "EarthCube Council of Data Facilities Registry Working Group", IndexName: "rwg"}
	pa = append(pa, rwg)

	csdco := Provider{Name: "CSDCO", Description: "Neotoma", IndexName: "csdco"}
	pa = append(pa, csdco)

	data, _ := json.Marshal(pa)

	response.Write(data)
}

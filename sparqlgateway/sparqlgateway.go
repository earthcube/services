package sparqlgateway

import (
	"encoding/json"
	"log"

	restful "github.com/emicklei/go-restful"
)

type URLSet []string

// New fires up the services to query the graph
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/graph").
		Doc("P418 graph driven API calls").
		Consumes("application/x-www-form-urlencoded").
		Produces(restful.MIME_JSON) //Consumes(restful.M).

	// add in start point and length cursors
	service.Route(service.GET("/resdetails").To(ResourceCall).
		Doc("Call for details on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]ResourceResults{}).
		Operation("ResourceCall"))

	service.Route(service.POST("/ressetdetails").To(ResourceSetCall).
		Doc("Call for details on an array of resources from the triplestore (graph)").
		Param(service.BodyParameter("body", "The body containing an array of URIs to obtain parameter values from")).
		Writes([]ResourceResults{}).
		Operation("ResourceSetCall"))

	return service
}

// ResourceCall call for details on the resource from the graph
func ResourceCall(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := ResCall(resource)
	response.WriteEntity(sr)
}

// ResourceSetCall call for details on the resource array from the graph
func ResourceSetCall(request *restful.Request, response *restful.Response) {

	log.Println("Get request body")
	body, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
	}
	log.Println(body)

	// ja := []byte(`["<https://www.bco-dmo.org/dataset/3300>", "<http://opencoredata.org/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816>"]`)
	ja := []byte(body)

	var jas URLSet

	err = json.Unmarshal(ja, &jas)
	if err != nil {
		log.Println("error with unmarshal..   return http error")
	}

	sr := ResSetCall(jas)
	response.WriteEntity(sr)
}

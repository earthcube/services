package sparqlgateway

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
)

type URLSet []string

// New fires up the services to query the graph
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/graph").
		Doc("Graph query services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) //Consumes(restful.M).

	service.Route(service.POST("/ressetdetails").To(ResourceSetCall).
		Doc("Call for details on an array of resources from the triplestore (graph)").
		Param(service.BodyParameter("body", "The body containing an array of URIs to obtain parameter values from")).
		Writes([]ResourceResults{}).
		Consumes("application/x-www-form-urlencoded").
		Operation("ResourceSetCall"))

	service.Route(service.POST("/ressetpeople").To(ResourceSetPeopleCall).
		Doc("Call for people associated with an array of resources from the triplestore (graph)").
		Param(service.BodyParameter("body", "The body containing an array of URIs to obtain people relation values from")).
		Writes([]ResourceSetPeople{}).
		Consumes("application/x-www-form-urlencoded").
		Operation("ResourceSetPeopleCall"))

	return service
}

// Dev fires up the DEVELOPMENT services
func Dev() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/dev/graph").
		Doc("Graph query services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) //Consumes(restful.M).

	// add in start point and length cursors
	service.Route(service.GET("/resdetails").To(ResourceCall).
		Doc("Call for details on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]ResourceResults{}).
		Operation("ResourceCall"))

	service.Route(service.GET("/logo").To(Logo).
		Doc("Call for logo URL on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Logo"))

	return service
}

// Logo call for details on the resource from the graph
func Logo(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := LogoCall(resource)
	response.WriteEntity(sr)
}

// ResourceCall call for details on the resource from the graph
func ResourceCall(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := ResCall(resource)
	response.WriteEntity(sr)
}

// ResourceSetCall call for details on the resource array from the graph
func ResourceSetCall(request *restful.Request, response *restful.Response) {

	body, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ja := []byte(body)
	var jas URLSet
	err = json.Unmarshal(ja, &jas)
	if err != nil {
		log.Println("error with unmarshal..   return http error")
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	sr := ResSetCall(jas)
	response.WriteEntity(sr)
}

// ResourceSetPeopleCall call for details on the resource array from the graph
func ResourceSetPeopleCall(request *restful.Request, response *restful.Response) {

	body, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ja := []byte(body)
	var jas URLSet
	err = json.Unmarshal(ja, &jas)
	if err != nil {
		log.Println("error with unmarshal..   return http error")
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	sr := ResSetPeople(jas)
	response.WriteEntity(sr)
}

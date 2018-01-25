package sparqlgateway

import (
	restful "github.com/emicklei/go-restful"
)

// ResourceResults is a place holder struct
type ResourceResults struct {
	Val     string
	Desc    string
	PubName string
	PubURL  string
}

// New fires up the services to query the graph
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/graph").
		Doc("P418 graph driven API calls").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// add in start point and length cursors
	service.Route(service.GET("/resdetails").To(ResourceCall).
		Doc("Call for details on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]ResourceResults{}).
		Operation("SearchCall"))

	return service
}

// ResourceCall call for details on the resource from the graph
func ResourceCall(request *restful.Request, response *restful.Response) {

	phrase := request.QueryParameter("r")

	sr := ResCall(phrase)

	response.WriteEntity(sr)

}

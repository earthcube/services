package main

import (
	"log"
	"net/http"

	"earthcube.org/Project418/services/spatialsearch"
	"earthcube.org/Project418/services/textsearch"
	"earthcube.org/Project418/services/typeahead"
	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
)

func main() {
	wsContainer := restful.NewContainer()

	// CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	// Add the services
	wsContainer.Add(textsearch.New())    // text search services
	wsContainer.Add(spatialsearch.New()) // text spatial services
	wsContainer.Add(typeahead.New())     // text spatial services
	// wsContainer.Add(graphsearch.New())  // text graph services

	// Swagger
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		ApiPath:        "/apidocs.json",
		WebServicesUrl: "http://geodex.org"} // localhost:6789
	swagger.RegisterSwaggerService(config, wsContainer)

	// Start up
	log.Printf("Services on localhost:6789")
	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()
}

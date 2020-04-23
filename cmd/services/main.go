package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"earthcube.org/Project418/services/internal/sparqlgateway"
	"earthcube.org/Project418/services/internal/spatialsearch"
	"earthcube.org/Project418/services/internal/textsearch"
	"earthcube.org/Project418/services/internal/typeahead"
	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
)


func main() {

	var host string
	flag.StringVar(&host, "host", "geodex.org", "Web services host")
	flag.Parse()

	wsContainer := restful.NewContainer()

	// CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	// Add the services
	wsContainer.Add(textsearch.New())    // text search services
	wsContainer.Add(spatialsearch.New()) // spatial services
	wsContainer.Add(typeahead.New())     // typeahead services
	wsContainer.Add(sparqlgateway.New()) // graph services
	wsContainer.Add(sparqlgateway.Dev()) // DEV graph services

	// Swagger
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		ApiPath:        "/apidocs.json",
		WebServicesUrl: "http://" + host } 
	swagger.RegisterSwaggerService(config, wsContainer)

	// Start up
	log.Printf("Services on localhost:6789")
	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()
}

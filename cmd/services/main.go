package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"earthcube.org/Project418/services/internal/sparqlgateway"
	"earthcube.org/Project418/services/internal/spatialsearch"
	"earthcube.org/Project418/services/internal/textsearch"
	"earthcube.org/Project418/services/internal/typeahead"
	"earthcube.org/Project418/services/internal/utils"
	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr, can be any io.Writer
	// I override this and set output to file (io.Writer) in main
	log.SetOutput(os.Stdout)

	// Set log level
	// Will log anything that is info or above (warn, error, fatal, panic). Default.
	// only other level is debug
	log.SetLevel(log.DebugLevel) // Info level for deployment
}

func main() {
	// Set up our log file for runs...
	serviceslogfile := utils.GetEnv("SERVICES_LOGFILE","./runtime/log/serviceslog.txt" )
	f, err := os.OpenFile(serviceslogfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// env
	geodexBaseUrl := utils.GetEnv("GEODEX_BASEURL","http://geodex.org" )
	geodexPort := utils.GetEnv("GEODEX_PORT",":6789" )

	servicesBaseUrl := fmt.Sprintf("%s%s",geodexBaseUrl,geodexPort)

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
		WebServicesUrl: servicesBaseUrl} // localhost:6789
	swagger.RegisterSwaggerService(config, wsContainer)

	// Start up
	log.Printf("Services on %s %s", geodexBaseUrl, geodexPort)
	server := &http.Server{Addr: geodexPort, Handler: wsContainer}
	server.ListenAndServe()
}

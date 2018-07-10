package shacl

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/emicklei/go-restful"
	"github.com/piprate/json-gold/ld"
)

// New builds out the service calls, in this case for the SHACL services
func New() *restful.WebService {
	service := new(restful.WebService)

	service.Path("/api/beta/shacl").
		Doc("SHACL testing APIs...  much changing here...  do not trust... NOT performant!")
		// Consumes(restful.MIME_JSON).
		// Produces(restful.MIME_JSON)

	service.Route(service.GET("/test").To(Test).
		Doc("Simple test...   rather useless").
		ReturnsError(400, "Unable to handle request", nil).
		Writes("plain/text").
		Operation("Test"))

	service.Route(service.POST("/igsn").To(IGSNCall).
		Doc("Test a data graph against a IGSN sample ID shape constraint").
		Param(service.BodyParameter("body", "The body containing a JSON-LD schame.org based documents")).
		Consumes("multipart/form-data").
		Produces("plain/text").
		ReturnsError(400, "Unable to handle request", nil).
		Operation("IGSNCall"))

	service.Route(service.POST("/eval").To(EvalCall).
		Doc("Test a data graph against a shapes graph ").
		Param(service.FormParameter("datagraph", "The body containing a data graph in turtle")).
		Param(service.FormParameter("shapesgraph", "The body containing a shape graph in turtle")).
		Consumes("multipart/form-data").
		Produces("plain/text").
		ReturnsError(400, "Unable to handle request", nil).
		Operation("EvalCall"))

	return service
}

// EvalCall test a provided data graph against a known IGSN shape graph
func EvalCall(request *restful.Request, response *restful.Response) {

	infile, _, err := request.Request.FormFile("shapesgraph")
	if err != nil {
		return
	}
	sbnr := bufio.NewReader(infile)
	shapegraph, err := ioutil.ReadAll(sbnr)
	if err != nil {
		return
	}

	dinfile, _, err := request.Request.FormFile("datagraph")
	if err != nil {
		return
	}
	dbnr := bufio.NewReader(dinfile)
	datagraph, err := ioutil.ReadAll(dbnr)
	if err != nil {
		return
	}

	log.Printf("Datagraph: %s \n", string(datagraph))
	log.Printf("Shapesgraph: %s \n", string(shapegraph))

	results := shaclExec(string(datagraph), string(shapegraph)) // TODO  needs to return error too
	response.Write([]byte(results))
}

// IGSNCall test a provided data graph against a known IGSN shape graph
func IGSNCall(request *restful.Request, response *restful.Response) {
	datagraph, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v with %s \n", err, datagraph)
	}

	// TODO  validate the JSONLD..  error back if it doesn't work
	// ttl, err := jsonLDToTTL(datagraph)
	// if err != nil {
	// 	log.Println("Error converting the POST body RDF")
	// 	// TODO response with error here...
	// }

	shapegraph, err := ioutil.ReadFile("./scripts/shape.ttl")
	if err != nil {
		log.Println("Error reading local shape graph")
		// TODO response with error here
	}

	results := shaclExec(datagraph, string(shapegraph)) // TODO  needs to return error too
	response.Write([]byte(results))
}

func shaclExec(datagraph, shapegraph string) string {
	// temporary data graph file
	df, err := ioutil.TempFile("", "datagraph")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(df.Name())
	df.Write([]byte(datagraph))

	// temporary shape graph file
	sf, err := ioutil.TempFile("", "shapegraph")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(sf.Name())
	sf.Write([]byte(shapegraph))

	app := "./scripts/shacl-1.0.0/bin/shaclvalidate.sh"
	args := []string{"-datafile", df.Name(), "-shapesfile", sf.Name()}
	cmd := exec.Command(app, args...)
	//	cmd.Stdin = strings.NewReader("some input")
	log.Println(cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("combined out:\n%s\n", string(out))

	return string(out)
}

func jsonLDToTTL(jsonld string) (string, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/nquads"

	var myInterface interface{}
	err := json.Unmarshal([]byte(jsonld), &myInterface)
	if err != nil {
		log.Printf("Error when transforming JSON-LD document to interface: %v", err)
		return "", err
	}

	ttl, err := proc.ToRDF(myInterface, options) // returns triples but toss them, we just want to see if this processes with no err
	if err != nil {
		log.Printf("Error when transforming JSON-LD document to RDF: %v", err)
		return "", err
	}

	return ttl.(string), err
}

// Test test a provided data graph against a known IGSN shape graph
func Test(request *restful.Request, response *restful.Response) {
	app := "./scripts/shacl-1.0.0/bin/shaclvalidate.sh"
	args := []string{"-datafile", "example1.ttl", "-shapesfile", "shape.ttl"}
	cmd := exec.Command(app, args...)
	//	cmd.Stdin = strings.NewReader("some input")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("combined out:\n%s\n", string(out))

	response.Write(out)
}

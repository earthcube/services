package sparqlgateway

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
	sparql "github.com/knakk/sparql"
)

const tq = `
# Comments are ignored, except those tagging a query.

# tag: timesearch
PREFIX  xsd:    <http://www.w3.org/2001/XMLSchema#>
PREFIX  dc:     <http://purl.org/dc/elements/1.1/>
PREFIX  :       <.>

SELECT *
{
    GRAPH ?g
      {
       ?s <http://geoschemas.org/contexts/temporal.jsonldtemporalCoverage>/<http://www.w3.org/2006/time#hasEnd>/<http://www.w3.org/2006/time#inXSDDateTimeStamp> ?end .
       ?s <http://geoschemas.org/contexts/temporal.jsonldtemporalCoverage>/<http://www.w3.org/2006/time#hasBeginning>/<http://www.w3.org/2006/time#inXSDDateTimeStamp> ?begin .
       FILTER (?begin > "{{.Begin}}"^^xsd:dateTime)
       FILTER (?end < "{{.End}}"^^xsd:dateTime)
     }
}

`

type params struct {
	Begin string
	End   string
}

// Temporal DEV call for time defined data in the RDF graph
// FILTER (?begin > "2019-04-21T00:00:00Z"^^xsd:dateTime)
func Temporal(request *restful.Request, response *restful.Response) {
	b := request.QueryParameter("b")
	e := request.QueryParameter("e")

	sr := temporalQuery(b, e)
	// fmt.Println(sr)
	// response.WriteJson(string(sr), " ")
	response.AddHeader("Content-Type", "application/json")
	response.Write(sr)
}

// temporalQuery takes a single resource and returns the variable measured property value
func temporalQuery(b, e string) []byte {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(tq)
	bank := sparql.LoadBank(f)

	p := params{Begin: b, End: e}

	q, err := bank.Prepare("timesearch", p)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.WithFields(log.Fields{
		"SPARQL": q,
	}).Info("A SPARQL call in P418 services")

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(res)

	// rr := &ResourceResults
	j, err := json.MarshalIndent(res, " ", "")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(j))

	// return fmt.Sprintf("%v", res.Bindings())
	// for this one don't return the map..  return JSON of the results
	return j
}

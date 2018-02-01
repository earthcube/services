package sparqlgateway

import (
	"bytes"
	"log"
	"strings"
	"time"

	sparql "github.com/knakk/sparql"
)

// ResourceResults is a place holder struct
type ResourceResults struct {
	Val     string
	Desc    string
	PubName string
	PubURL  string
}

const queries = `
# Comments are ignored, except those tagging a query.

# tag: ResourceResults
prefix schema: <http://schema.org/>
SELECT ?val ?desc ?pubname ?puburl
WHERE
{
  BIND(<{{.RESID}}> AS ?ID)
  ?ID schema:publisher ?pub .
  ?pub schema:name ?pubname .
  ?pub schema:url ?puburl .
  ?ID schema:variableMeasured ?res  .
  ?res a schema:PropertyValue .
  ?res schema:value ?val   .
  ?res schema:description ?desc     
} 

# tag: ResourceSetResults
prefix schema: <http://schema.org/>
SELECT DISTINCT ?val ?desc ?pubname ?puburl
WHERE
{
VALUES ?ID
{  {{.}}
}
?ID schema:variableMeasured ?res .
OPTIONAL {
?res a schema:PropertyValue .
?res schema:value ?val .
?res schema:description ?desc
}
OPTIONAL {
?ID schema:publisher ?pub .
OPTIONAL { ?pub schema:name ?pubname }
OPTIONAL { ?pub schema:url ?puburl }
}
}

`

func getP418SPARQL() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://geodex.org/blazegraph/namespace/p418/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

// ResSetCall takes a single resource and returns the variable measured property value
func ResSetCall(resources URLSet) []ResourceResults {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("ResourceSetResults", strings.Join(resources, " "))
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	rra := []ResourceResults{}
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		rr := ResourceResults{Val: i["val"].Value, Desc: i["desc"].
			Value, PubName: i["pubname"].Value, PubURL: i["puburl"].Value}
		rra = append(rra, rr)
	}

	return rra
}

// ResCall takes a single resource and returns the variable measured property value
func ResCall(resource string) []ResourceResults {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("ResourceResults", struct{ RESID string }{resource})
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	rra := []ResourceResults{}
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		rr := ResourceResults{Val: i["val"].Value, Desc: i["desc"].
			Value, PubName: i["pubname"].Value, PubURL: i["puburl"].Value}
		rra = append(rra, rr)
	}

	return rra
}

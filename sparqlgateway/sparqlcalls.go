package sparqlgateway

import (
	"bytes"
	"log"
	"time"

	sparql "github.com/knakk/sparql"
)

const queries = `
# Comments are ignored, except those tagging a query.

# tag: resInfo
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

func ResCall(resource string) *sparql.Results {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("resInfo", struct{ RESID string }{resource})
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res
}

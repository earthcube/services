package vocabulary

import (
	"bytes"
	"log"
	"strings"
	"time"

	"github.com/knakk/sparql"
)

const queries = `
# Place the SAPRQL here..

#####  VOCABULARY ITEMS #####
## just list all the vocs used  (just a graph search?)
# tag: voclist 

## find a voc based on some set of terms...   
# tag: vocsearch 

## what info (descibe?) could we get on this?
# tag: vocinfo

######  TERM ITEMS  ########
##  (should we optionally type this?  scheme:MeasuredVariable   then allow others in the future)
# tag: termlist  

## Use the SPARQL free text BIF in blaze to search  (can this be done in pure SPARQL 1.1?)
## look for terms (measuredvar?) that matches a string
# tag: termsearch

## what info (descibe?) could we get on this?
# tag: terminfo


##### AGENT ITEMS #########
## list all type PERSON or type ORG 
# tag: agentlist

# tag: agentsearch

# tag: agentinfo


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

// List gets a list of all unique vocabularies (or parameters) in the graph
func List(resources URLSet) []ResourceSetPeople {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("ListVoc", strings.Join(resources, " "))
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	rra := []ResourceSetPeople{}
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		rr := ResourceSetPeople{G: i["g"].Value, Person: i["person"].Value, Rolename: i["rolename"].
			Value, Name: i["name"].Value, URL: i["url"].Value, Orcid: i["orcid"].Value}
		rra = append(rra, rr)
	}

	return rra
}

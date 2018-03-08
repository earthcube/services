package sparqlgateway

import (
	"bytes"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	sparql "github.com/knakk/sparql"
)

// LogoResults is a place holder struct
type LogoResults struct {
	Graph    string
	Type     string
	Resource string
	Logo     string
}

// ResourceResults is a place holder struct
type ResourceResults struct {
	Val     string
	Desc    string
	PubName string
	PubURL  string
}

// DetailResults is a place holder struct
type DetailResults struct {
	S             string
	Aname         string
	Name          string
	URL           string
	Description   string
	Citation      string
	Datepublished string
	Curl          string
	Keywords      string
	License       string
}

// ResourceSetPeople struct
type ResourceSetPeople struct {
	G        string
	Person   string
	Rolename string
	Name     string
	URL      string
	Orcid    string
}

func getP418SPARQL() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://geodex.org/blazegraph/namespace/p418/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

// ResSetPeople takes a single resource and returns the variable measured property value
func ResSetPeople(resources URLSet) []ResourceSetPeople {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("ResourceSetPeople", strings.Join(resources, " "))
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

// DetailsCall takes a single resource and returns the variable measured property value
func DetailsCall(resources string) DetailResults {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("detailsCall", resources)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// rra := []DetailResults{}
	bindings := res.Results.Bindings // map[string][]rdf.Term
	// for _, i := range bindings {
	// x := bindings["s"][0].Value

	rr := DetailResults{S: bindings[0]["s"].Value,
		Aname: bindings[0]["aname"].Value, URL: bindings[0]["url"].Value,
		Description: bindings[0]["description"].Value, Citation: bindings[0]["citation"].Value,
		Datepublished: bindings[0]["datepublished"].Value,
		Curl:          bindings[0]["curl"].Value, Keywords: bindings[0]["keywords"].Value, License: bindings[0]["license"].Value}
	// rra = append(rra, rr)
	// }

	return rr
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

// LogoCall takes a single resource and returns the variable measured property value
func LogoCall(resource string) []LogoResults {
	repo, err := getP418SPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("LogoCall", struct{ RESID string }{resource})
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	rra := []LogoResults{}
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		rr := LogoResults{Graph: i["graph"].Value, Type: i["type"].Value, Resource: i["resource"].Value, Logo: i["logo"].Value}
		rra = append(rra, rr)
	}

	return rra
}

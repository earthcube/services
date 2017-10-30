package textsearch

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	restful "github.com/emicklei/go-restful"
)

// Foo is a place holder struct
type Foo struct {
	Item string
}

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/textindex").
		Doc("P418 text search API").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/test").To(SampleCall).
		Doc("Testing call").
		Param(service.PathParameter("test", "TESTING: a simple service to roundtrip the code").DataType("string")).
		Writes(Foo{}).
		Operation("SampleCall"))

	service.Route(service.GET("/search/{term}").To(SearchCall).
		Doc("Search call").
		Param(service.PathParameter("search", "TESTING: first search test").DataType("string")).
		Writes(Foo{}).
		Operation("SearchCall"))

	return service
}

// SampleCall is a simple sample service for testing....
func SampleCall(request *restful.Request, response *restful.Response) {

	allitems := Foo{Item: "string to send"}
	response.WriteEntity(allitems)
}

// First test function..   opens each time..  not what we want..
// need to open indexes and maintain state
func SearchCall(request *restful.Request, response *restful.Response) {

	// Old func line
	// func searchCall(phrase string, searchIndex string) string {
	phrase := request.PathParameter("term")

	searchIndex := "" // use all indexes for testing now...

	log.Printf("Search Term: %s \n", phrase)

	// Open all the index files
	// TODO  really should only open the ones I already know will be in the index alias
	index1, err := openIndex("/Users/dfils/Data/OCDDataVolumes/indexes/abstracts.bleve")
	if err != nil {
		log.Printf("Error with index1 alias: %v", err)
	}
	index2, err := openIndex("/Users/dfils/Data/OCDDataVolumes/indexes/csdco.bleve")
	if err != nil {
		log.Printf("Error with index2 alias: %v", err)
	}
	index3, err := openIndex("/Users/dfils/Data/OCDDataVolumes/indexes/janus.bleve")
	if err != nil {
		log.Printf("Error with index3 alias: %v", err)
	}

	var index bleve.IndexAlias

	if searchIndex == "abstracts" {
		index = bleve.NewIndexAlias(index1)
		log.Println("abstract index only")
	}
	if searchIndex == "csdco" {
		index = bleve.NewIndexAlias(index2)
		log.Println("CSDCO index only")
	}
	if searchIndex == "jrso" {
		index = bleve.NewIndexAlias(index3)
		log.Println("JRSO index only")
	} else {
		index = bleve.NewIndexAlias(index1, index2, index3)
		log.Println("All indexes active")
	}

	// Set up query and search
	// query := bleve.NewMatchQuery(phrase)
	query := bleve.NewQueryStringQuery(phrase)
	search := bleve.NewSearchRequestOptions(query, 10, 0, false) // no explanation
	search.Highlight = bleve.NewHighlight()                      // need Stored and IncludeTermVectors in index
	// search.Highlight = bleve.NewHighlightWithStyle("html") // need Stored and IncludeTermVectors in index

	var jsonResults []byte // will hold our results

	// do search and get results
	searchResults, err := index.Search(search)
	if err != nil {
		log.Printf("Error in search call: %v", err)
	} else {
		hits := searchResults.Hits
		jsonResults, err = json.MarshalIndent(hits, " ", " ")
		if err != nil {
			log.Printf("Error with json marshal call: %v", err)
		}

		// testing print loop
		for k, item := range hits {
			fmt.Printf("\n%d: %s, %f, %s, %v\n", k, item.Index, item.Score, item.ID, item.Fragments)
			for key, frag := range item.Fragments {
				fmt.Printf("%s   %s\n", key, frag)
			}
		}
	}

	response.WriteEntity(string(jsonResults))
}

func openIndex(indexPath string) (bleve.Index, error) {
	var bleveIdx bleve.Index

	if bleveIdx == nil {
		var err error
		bleveIdx, err = bleve.OpenUsing(indexPath, map[string]interface{}{"read_only": true})
		if err != nil {
			return nil, err
		}
	}

	return bleveIdx, nil
}

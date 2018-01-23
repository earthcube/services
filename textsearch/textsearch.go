package textsearch

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	restful "github.com/emicklei/go-restful"
)

// Foo is a place holder struct
type OrganicResults struct {
	Position int
	Index    string
	Score    float64
	ID       string
}

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/textindex").
		Doc("P418 text search API").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

		// add in start point and length cursors
	service.Route(service.GET("/search").To(SearchCall).
		Doc("Search call").
		Param(service.QueryParameter("q", "Query string").DataType("string")).
		Param(service.QueryParameter("s", "Starting cursor point").DataType("int")).
		Param(service.QueryParameter("n", "Number of results to return").DataType("int")).
		Param(service.QueryParameter("i", "Index to use.  Currently one of; ocd, bcodmo, linkedearth").DataType("string")).
		Writes([]OrganicResults{}).
		Operation("SearchCall"))

	return service
}

// SearchCall First test function..   opens each time..  not what we want..
// need to open indexes and maintain state
func SearchCall(request *restful.Request, response *restful.Response) {

	// Old func line
	// func searchCall(phrase string, searchIndex string) string {
	phrase := request.QueryParameter("q")
	startPoint, err := strconv.ParseInt(request.QueryParameter("s"), 10, 32)
	if err != nil {
		log.Printf("Error with index1 alias: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	numToReturn, err := strconv.ParseInt(request.QueryParameter("n"), 10, 32)
	if err != nil {
		log.Printf("Error with index1 alias: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	searchIndex := "" // use all indexes for testing now...
	searchIndex = request.QueryParameter("i")
	// TODO.  Make a string array of search index options and then test here to make sure
	// searchIndex is NOT "" that it is in the array via contains call

	log.Printf("Search Term: %s \n", phrase)

	// Open all the index files
	// TODO  really should only open the ones I already know will be in the index alias
	index1, err := openIndex("indexes/bcodmo.bleve")
	if err != nil {
		log.Printf("Error with index1 alias: %v", err)
	}
	index2, err := openIndex("indexes/ocd.bleve")
	if err != nil {
		log.Printf("Error with index2 alias: %v", err)
	}
	index3, err := openIndex("indexes/linkedearth.bleve")
	if err != nil {
		log.Printf("Error with index3 alias: %v", err)
	}

	var index bleve.IndexAlias

	if searchIndex == "bcodmo" {
		index = bleve.NewIndexAlias(index1)
		log.Println("BCODMO index only")
	}
	if searchIndex == "ocd" {
		index = bleve.NewIndexAlias(index2)
		log.Println("OCD index only")
	}
	if searchIndex == "linkedearth" {
		index = bleve.NewIndexAlias(index3)
		log.Println("LinkedEarth index only")
	} else {
		index = bleve.NewIndexAlias(index1, index2, index3)
		log.Println("All indexes active")
	}

	// Set up query and search
	// query := bleve.NewMatchQuery(phrase)
	query := bleve.NewQueryStringQuery(phrase)
	search := bleve.NewSearchRequestOptions(query, int(numToReturn), int(startPoint), false) // no explanation
	// search.Highlight = bleve.NewHighlight()                      // need Stored and IncludeTermVectors in index
	search.Highlight = bleve.NewHighlightWithStyle("html") // need Stored and IncludeTermVectors in index

	// var jsonResults []byte // will hold our results
	var ora []OrganicResults

	// do search and get results
	searchResults, err := index.Search(search)
	if err != nil {
		log.Printf("Error in search call: %v", err)
	} else {
		hits := searchResults.Hits
		// jsonResults, err = json.MarshalIndent(hits, "", " ")
		if err != nil {
			log.Printf("Error with json marshal call: %v", err)
		}

		// testing print loop
		for k, item := range hits {
			ors := OrganicResults{Position: k, Index: item.Index, Score: item.Score, ID: item.ID}
			ora = append(ora, ors)
			fmt.Printf("\n%d: %s, %f, %s, %v\n", k, item.Index, item.Score, item.ID, item.Fragments)
			for key, frag := range item.Fragments {
				fmt.Printf("%s   %s\n", key, frag)
			}
		}
	}

	// response.WriteEntity(string(jsonResults))
	response.WriteEntity(ora)
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

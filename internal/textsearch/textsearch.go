package textsearch

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
	"gopkg.in/resty.v1"

	"github.com/blevesearch/bleve"
	restful "github.com/emicklei/go-restful"
)

type Provider struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	IndexName     string `json:"index"`
	IndexLocation string `json:"indexlocation"`
}

// OrganicResultsSet has top N results from each provider with scores
type OrganicResultsSet struct {
	OR        []OrganicResults `json:"or"`        // provider:results
	HighScore float64          `json:"highscore"` // provider:highestScore
	Index     string           `json:"index"`     // ordered string array based on score
}

// OrganicResults is a place holder struct
type OrganicResults struct {
	Position  int     `json:"position"`
	IndexPath string  `json:"indexpath"`
	Score     float64 `json:"score"`
	ID        string  `json:"URL"`
}

type byScore []OrganicResultsSet // for our custom sorting of orsa

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/textindex").
		Doc("Organic free text search services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

		// add in start point and length cursors
	service.Route(service.GET("/search").To(SearchCall).
		Doc("Query string search call").
		Param(service.QueryParameter("q", "Query string").DataType("string")).
		Param(service.QueryParameter("s", "Starting cursor point").DataType("int")).
		Param(service.QueryParameter("n", "Number of results to return").DataType("int")).
		Param(service.QueryParameter("i", "Index to use.  Currently 1 or more (comma sparated) of: ocd, bcodmo, ieda, neotoma, rwg, linkedearth").DataType("string")).
		Writes([]OrganicResults{}).
		Operation("SearchCall"))

	service.Route(service.GET("/searchset").To(SearchSetCall).
		Doc("Query string search call with grouped results").
		Param(service.QueryParameter("q", "Query string").DataType("string")).
		Param(service.QueryParameter("s", "Starting cursor point").DataType("int")).
		Param(service.QueryParameter("n", "Number of results to return").DataType("int")).
		Writes([]OrganicResults{}).
		Operation("SearchSetCall"))

	service.Route(service.POST("/nusearch").To(NuSearch).
		Doc("BETA: Query string search call ").
		Param(service.FormParameter("body", "The body containing query document")).
		Consumes("multipart/form-data").
		// Produces("plain/text").
		ReturnsError(400, "Unable to handle request", nil).
		Operation("NuSearch"))

	service.Route(service.GET("/getnusearch").To(GETNuSearch).
		Doc("BETA: Query string search call ").
		Param(service.QueryParameter("q", "Query string").DataType("string")).
		Operation("GETNuSearch"))

	return service
}

// GETNuSearch is a test of the Blast search package
func GETNuSearch(request *restful.Request, response *restful.Response) {
	phrase := request.QueryParameter("q")
	log.Printf("Body %s\n", phrase)

	// addressing   {code: 4, message: "context deadline exceeded"}
	// ctx := context.TODO()
	ctx, cancel := context.WithTimeout(context.Background(), 8000*time.Millisecond)
	defer cancel()

	// u := "http://blast:10002/rest/_search"
	u := "http://localhost:10002/rest/_search"

	resp, err := resty.R().
		SetBody(phrase).SetContext(ctx).
		Post(u)
	if err != nil {
		log.Print(err)
	}

	// If we excede the context time limit the first time try again...
	// look for "message": "context deadline exceeded"  in resp.String()
	s := fastjson.GetString([]byte(resp.String()), "error", "message")
	if s == "context deadline exceede" {
		log.Println("Try the search once more")
		resp, err = resty.R().
			SetBody(phrase).SetContext(ctx).
			Post(u)
		if err != nil {
			log.Print(err)
		}
	}

	log.Printf("\nInput: %s", phrase)
	log.Printf("\nError: %v", err)
	log.Printf("\nResponse Status Code: %v", resp.StatusCode())
	log.Printf("\nResponse Status: %v", resp.Status())
	log.Printf("\nResponse Time: %v", resp.Time())
	log.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	log.Printf("\nResponse Body: %v", resp.String()) // or resp.String() or string(resp.Body())

	response.Write([]byte(resp.String()))
}

// NuSearch is a test of the Blast search package
func NuSearch(request *restful.Request, response *restful.Response) {
	log.Println("In the nucall")
	log.Print(request.Request.Form)

	infile, _, err := request.Request.FormFile("body")
	if err != nil {
		log.Println(err)
		return
	}
	sbnr := bufio.NewReader(infile)
	body, err := ioutil.ReadAll(sbnr)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Body %s\n", string(body))

	resp, err := resty.R().
		SetBody(body).
		Post("http://blast:10002/rest/_search")
	if err != nil {
		log.Print(err)
	}

	log.Printf("\nInput: %s", body)
	log.Printf("\nError: %v", err)
	log.Printf("\nResponse Status Code: %v", resp.StatusCode())
	log.Printf("\nResponse Status: %v", resp.Status())
	log.Printf("\nResponse Time: %v", resp.Time())
	log.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	log.Printf("\nResponse Body: %v", resp.String()) // or resp.String() or string(resp.Body())

	response.Write([]byte(resp.String()))
}

// SearchCall First test function..   opens each time..  not what we want..
// need to open indexes and maintain state
func SearchCall(request *restful.Request, response *restful.Response) {
	phrase := request.QueryParameter("q")
	log.Printf("Search Term: %s \n", phrase)

	startPoint, err := strconv.ParseInt(request.QueryParameter("s"), 10, 32)
	if err != nil {
		log.Printf("Error with starting index value: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	numToReturn, err := strconv.ParseInt(request.QueryParameter("n"), 10, 32)
	if err != nil {
		log.Printf("Error with number requested value: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Neither of the index or number requested can be less than 1
	if numToReturn < 1 || startPoint < 0 {
		log.Printf("Requested index or return value of 0 or negative: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	im := indexMap() // index maps  (TODO later build from a config file, so not hard coded)

	// get the request index string "i" and parse it to an array
	searchIndex := "" // use all indexes for testing now...
	searchIndex = request.QueryParameter("i")

	// Pull and parse the string array
	var sia []string
	if searchIndex != "" {
		sia = strings.Split(searchIndex, ",")
		if len(sia) > 0 {
			for index := range sia {
				indexname := sia[index] // get an element from the array, then check it..
				if val, ok := im[indexname]; !ok {
					log.Printf("Requested unknown index %s, %s", searchIndex, val)
					response.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			}
		}
	}

	if len(sia) == 0 {
		log.Printf("We seem to have no index set..   SO use them all!  :)   ")
		im := indexMap() // index maps  (TODO later build from a config file, so not hard coded)
		for name := range im {
			sia = append(sia, name) // just put in the name, not the path..  I look that up later   (this could be written better!)
		}
	}

	index, err := getMultiIndexAlias(sia, im) // we have our index
	if err != nil {
		response.WriteErrorString(422, "Error getting a set of indexes to search on.  (getMultiIndexalias)")
		return
	}

	ora := getResultSet(index, phrase, numToReturn, startPoint)

	log.WithFields(log.Fields{
		"ora":    ora,
		"phrase": phrase,
		"start":  startPoint,
		"number": numToReturn,
	}).Info("An organic text call in P418 services SearchCall")

	response.WriteEntity(ora)
}

// SearchSetCall return a set of organic results from across all the providers
func SearchSetCall(request *restful.Request, response *restful.Response) {
	phrase := request.QueryParameter("q")

	startPoint, err := strconv.ParseInt(request.QueryParameter("s"), 10, 32)
	if err != nil {
		log.Printf("Error with index alias: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	numToReturn, err := strconv.ParseInt(request.QueryParameter("n"), 10, 32)
	if err != nil {
		log.Printf("Error with index1 alias: %v", err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var orsa []OrganicResultsSet

	im := indexMap() // index maps  (TODO later build from a config file, so not hard coded)
	for name := range im {
		index, err := getMultiIndexAlias([]string{name}, im) // we have our index
		if err != nil {
			response.WriteErrorString(422, "Error getting a index to search on.  (getMultiIndexalias)")
			return
		}
		ora := getResultSet(index, phrase, numToReturn, startPoint)
		score, err := maxFloat(ora) // set to highest score in ora..  deal with future error func return
		if err == nil {
			ors := OrganicResultsSet{OR: ora, HighScore: score, Index: name}
			orsa = append(orsa, ors)
		}
	}

	// sort array putting them in order of top score...
	sort.Sort(byScore(orsa))

	log.WithFields(log.Fields{
		"orsa":   orsa,
		"phrase": phrase,
		"start":  startPoint,
		"number": numToReturn,
	}).Info("An organic text call in P418 services SearchSetCall")

	response.WriteEntity(orsa)
}

func openIndex(indexPath string) (bleve.Index, error) {
	var bleveIdx bleve.Index

	if bleveIdx == nil {
		var err error
		bleveIdx, err = bleve.OpenUsing(indexPath, map[string]interface{}{"read_only": true})
		if err != nil {
			log.Printf("Error with an index opening: %v", err)
			return nil, err
		}
	}

	return bleveIdx, nil
}

func indexMap() map[string]string {
	ic, err := os.Open("./indexcatalog.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ic.Close()

	pa := []Provider{}
	jsonParser := json.NewDecoder(ic)
	jsonParser.Decode(&pa)

	im := make(map[string]string)

	for e := range pa {
		im[pa[e].IndexName] = pa[e].IndexLocation
	}

	// OLD static pattern..  remove this..
	// im["bcodmo"] = "indexes/bcodmo.bleve"
	// im["ocd"] = "indexes/ocd.bleve"
	// im["linkedearth"] = "indexes/linkedearth.bleve"
	// im["rwg"] = "indexes/rwg.bleve"
	// im["ieda"] = "indexes/ieda.bleve"
	// im["csdco"] = "indexes/csdco.bleve"
	// im["neotoma"] = "indexes/neotoma.bleve"

	log.Println(im)

	return im
}

// ref: http://www.blevesearch.com/docs/IndexAlias/
func getMultiIndexAlias(searchIndex []string, im map[string]string) (bleve.IndexAlias, error) {
	ia := make([]bleve.Index, 0)
	var err error

	for i := range searchIndex {
		index, err := openIndex(im[searchIndex[i]])
		if err != nil {
			log.Printf("Error with an index opening: %v", err) // really logged in openIndex
			return nil, err
		}
		ia = append(ia, index)
	}
	index := bleve.NewIndexAlias(ia...) // use variadic call

	return index, err
}

// func getIndexAlias(searchIndex string, im map[string]string) bleve.IndexAlias {
// 	var index1, index2, index3, index4, index5, index6 bleve.Index
// 	var err error
// 	var index bleve.IndexAlias
// 	if searchIndex == "" {
// 		index1, err = openIndex("indexes/bcodmo.bleve")
// 		index2, err = openIndex("indexes/ocd.bleve")
// 		index3, err = openIndex("indexes/linkedearth.bleve")
// 		index4, err = openIndex("indexes/rwg.bleve")
// 		index5, err = openIndex("indexes/ieda.bleve")
// 		index6, err = openIndex("indexes/csdco.bleve")
// 		if err != nil {
// 			log.Printf("Error with an index opening: %v", err) // really logged in openIndex
// 		}
// 		index = bleve.NewIndexAlias(index1, index2, index3, index4, index5, index6)
// 		log.Println("All indexes active")
// 	} else {
// 		index1, err := openIndex(im[searchIndex])
// 		if err != nil {
// 			log.Printf("Error with an index opening: %v", err)
// 		}
// 		index = bleve.NewIndexAlias(index1)
// 		log.Printf("Active index: %s", im[searchIndex])
// 	}
// 	return index
// }

func getResultSet(index bleve.IndexAlias, phrase string, numToReturn, startPoint int64) []OrganicResults {

	// Set up query and search.   OLD:  query := bleve.NewMatchQuery(phrase)
	query := bleve.NewQueryStringQuery(phrase)
	search := bleve.NewSearchRequestOptions(query, int(numToReturn), int(startPoint), false) // no explanation
	// search.Highlight = bleve.NewHighlight()                      // need Stored and IncludeTermVectors in index ?
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
			ors := OrganicResults{Position: k, IndexPath: item.Index, Score: shortenFloat(item.Score, 2), ID: item.ID}
			ora = append(ora, ors)
			// fmt.Printf("\n%d: %s, %f, %s, %v\n", k, item.Index, item.Score, item.ID, item.Fragments)
			fmt.Printf("  THIS IS THE ITEM %v\n", item)

			// for key, frag := range item.Fragments {
			// 	fmt.Printf("%s   %s\n", key, frag)
			// }
		}
	}

	return ora
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func shortenFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func maxFloat(ora []OrganicResults) (float64, error) {
	m := -1.0
	var err error
	for _, e := range ora {
		if e.Score > m {
			m = e.Score
		}
	}

	if m == -1 {
		err = errors.New("No items to score in the array")
	}

	return m, err
}

// Len, Swap, Less: Some sort logic to return orsa in a sorted order
func (s byScore) Len() int {
	return len(s)
}

func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byScore) Less(i, j int) bool {
	return s[i].HighScore > s[j].HighScore
}

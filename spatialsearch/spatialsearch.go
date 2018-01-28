package spatialsearch

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/emicklei/go-restful"
	"github.com/garyburd/redigo/redis"
	"github.com/paulmach/go.geojson"
)

// LocType is there to do a first cut marshalling to just get the type before  next marshalling
type LocType struct {
	Type string `json:"type"`
}

// // LocationPoint is a simple type and cooridnates struct for schema.org spatial info
// type LocationPoint struct {
// 	Type        string    `json:"type"`
// 	Coordinates []float64 `json:"coordinates"`
// }

// // LocationPoly is a simple type and cooridnates struct for schema.org spatial info
// type LocationPoly struct {
// 	Type        string        `json:"type"`
// 	Coordinates [][][]float64 `json:"coordinates"`
// }

// // LocationFeature is a simple type and cooridnates struct for schema.org spatial info
// type LocationFeature struct {
// 	Type     string       `json:"type"`
// 	Geometry LocationPoly `json:"geometry"`
// }

// // GeoJOSN struct
// type GeoJSON struct {
// 	Type     string        `json:"type"`
// 	Features []GeoFeatures `json:"features"`
// }

// // GeoFeatures
// type GeoFeatures struct {
// 	Type       string        `json:"type"`
// 	Properties GeoProperties `json:"properties"`
// 	Geometry   GeoGeometry   `json:"geometry"`
// }

// type GeoProperties struct {
// 	URL string `json:URL`
// }

// type GeoGeometry struct {
// 	Type        string    `json:"type"`
// 	Coordinates []float64 `json:"coordinates"`
// }

// New builds out the services calls..
func New() *restful.WebService {
	service := new(restful.WebService)

	service.
		Path("/api/v1/spatial").
		Doc("Spatial services to P418 holdings").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/search/test1").To(SpatialCall).
		Doc("get expeditions from a spatial polygon defined by wkt").
		Param(service.QueryParameter("geowithin", "Polygon in WKT format within which to look for features.  Try `POLYGON((-72.2021484375 38.58896696823242,-59.1943359375 38.58896696823242,-59.1943359375 28.11801628757283,-72.2021484375 28.11801628757283,-72.2021484375 38.58896696823242))`").DataType("string")).
		ReturnsError(400, "Unable to handle request", nil).
		Operation("SpatialCall"))
	return service
}

// SpatialCall calls to tile38 data store
func SpatialCall(request *restful.Request, response *restful.Response) {

	wktstring := request.QueryParameter("geowithin")

	// c, err := redis.Dial("tcp", "tile38:9851")
	c, err := redis.Dial("tcp", "localhost:9851")
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
	}
	defer c.Close()

	var value1 int
	var value2 []interface{}
	reply, err := redis.Values(c.Do("INTERSECTS", "p418", "LIMIT", "50000", "OBJECT", wktstring))
	// reply, err := redis.Values(c.Do("SCAN", "p418"))
	if err != nil {
		fmt.Printf("Error in reply %v \n", err)
	}
	if _, err := redis.Scan(reply, &value1, &value2); err != nil {
		fmt.Printf("Error in scan %v \n", err)
	}

	fmt.Println(value1)

	results, _ := templateTest(value2)
	response.Write([]byte(results))
}

func templateTest(results []interface{}) (string, error) {

	fc := geojson.NewFeatureCollection()

	for _, item := range results {
		valcast := item.([]interface{})
		val0 := fmt.Sprintf("%s", valcast[0])
		val1 := fmt.Sprintf("%s", valcast[1])

		// log.Printf("%s %s \n", val0, val1)

		lt := &LocType{}
		err := json.Unmarshal([]byte(val1), lt)
		if err != nil {
			log.Print(err)
			return "", err
		}

		rawGeometryJSON := []byte(val1)

		if lt.Type == "Point" || lt.Type == "Poly" {
			g, err := geojson.UnmarshalGeometry(rawGeometryJSON)
			if err != nil {
				log.Printf("Unmarshal geom error for %s with %s\n", val0, err)
			}

			switch {
			case g.IsPoint():
				nf := geojson.NewFeature(g)
				nf.SetProperty("URL", val0)
				fc.AddFeature(nf)
			case g.IsPolygon():
				nf := geojson.NewFeature(g)
				nf.SetProperty("URL", val0)
				fc.AddFeature(nf)
			default:
				log.Println(g.Type)
			}
		}

		if lt.Type == "Feature" {
			f, err := geojson.UnmarshalFeature(rawGeometryJSON)
			if err != nil {
				log.Printf("Unmarshal feature error for %s with %s\n", val0, err)
			}
			f.SetProperty("URL", val0)
			fc.AddFeature(f)
		}

	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(rawJSON), nil
}

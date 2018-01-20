package spatialsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/garyburd/redigo/redis"
	gj "github.com/kpawlik/geojson"
)

// Location is a simple type and cooridnates struct for schema.org spatial info
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

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
	// Consumes("application/vnd.flyovercountry.v1+json")  // Is this a good approach?
	// “application/vnd.laccore.flyovercountry+json; version=1”
	// “application/json; profile=vnd.laccore.flyovercountry version=1”
	// "application/json;vnd.laccore.flyovercountry+v1"

	return service
}

// SpatialCall calls to tile38 data store
func SpatialCall(request *restful.Request, response *restful.Response) {

	wktstring := request.QueryParameter("geowithin")

	c, err := redis.Dial("tcp", "tile38:9851")
	// c, err := redis.Dial("tcp", "localhost:9851")
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
	}
	defer c.Close()

	log.Print("connected")

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

	results, _ := tile38RespAsGeoJSON(value2)
	response.Write([]byte(results))
}

func tile38RespAsGeoJSON(results []interface{}) (string, error) {

	// Build the geojson section
	var (
		// fc *gj.FeatureCollection
		f  *gj.Feature
		fa []*gj.Feature
	)

	for _, item := range results {
		valcast := item.([]interface{})
		val0 := fmt.Sprintf("%s", valcast[0])
		val1 := fmt.Sprintf("%s", valcast[1])

		loc := &Location{}
		err := json.Unmarshal([]byte(val1), loc)
		if err != nil {
			return "", err
		}

		cd := gj.Coordinate{gj.Coord(loc.Coordinates[0]), gj.Coord(loc.Coordinates[1])} // is this long lat..  vs lat long?

		props := map[string]interface{}{"URL": val0}

		newp := gj.NewPoint(cd)
		f = gj.NewFeature(newp, props, nil)
		fa = append(fa, f)
	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(fc)
	if err != nil {
		log.Println(err)
	}

	return gjstr, nil
}

// WKTPolygonToFloatArray return float64 array for WKT Poly string
func WKTPolygonToFloatArray(wkt string) ([][][]float64, error) {
	// Take WKT string and parse down
	parsedString := strings.TrimSuffix(strings.TrimPrefix(wkt, "POLYGON(("), "))")
	wktarray := strings.Split(parsedString, ",")

	f := [][][]float64{}
	c := [][]float64{}

	for _, item := range wktarray {
		coordSet := strings.Split(item, " ")
		// TODO..  catch these errors..  this is bad form!  The whole function needs an error
		x, err := strconv.ParseFloat(coordSet[0], 64)
		y, err := strconv.ParseFloat(coordSet[1], 64)
		cd := []float64{x, y}
		c = append(c, cd)
		if err != nil {
			log.Println(err)
			return f, errors.New("Error in the conversion of WKT to GeoJSON Polygon to support spatial search")
		}
	}

	f = append(f, c)

	fmt.Println(f)
	return f, nil
}

// WKTPolygontoGeoJSON convert WKT to GeoJSON for Polygons.
// Ended up not using this in ocdService since mgo driver needed a bson structure for query
func WKTPolygontoGeoJSON(wkt string) string {
	var (
		// fc *gj.FeatureCollection
		//f  *gj.Feature
		//fa []*gj.Feature
		newp *gj.Polygon
	)

	// Take WKT string and parse down
	parsedString := strings.TrimSuffix(strings.TrimPrefix(wkt, "POLYGON(("), "))")
	wktarray := strings.Split(parsedString, ",")

	c := gj.Coordinates{}
	for _, item := range wktarray {
		coordSet := strings.Split(item, " ")
		// TODO..  catch these errors..  this is bad form!  The whole function needs an error
		x, _ := strconv.ParseFloat(coordSet[0], 64)
		y, _ := strconv.ParseFloat(coordSet[1], 64)
		cd := gj.Coordinate{gj.Coord(x), gj.Coord(y)}
		c = append(c, cd)

	}
	newml := gj.MultiLine{c}
	newp = gj.NewPolygon(newml)
	//f = gj.NewFeature(newp, nil, nil)
	//fa = append(fa, f)

	//fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(newp)
	if err != nil {
		//panic(err)
		log.Printf("Error event: %v \n", err)
	}

	fmt.Println(gjstr)

	return gjstr
}

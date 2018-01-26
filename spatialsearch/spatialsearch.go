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

// LocType is there to do a first cut marshalling to just get the type before  next marshalling
type LocType struct {
	Type string `json:"type"`
}

// LocationPoint is a simple type and cooridnates struct for schema.org spatial info
type LocationPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// LocationPoly is a simple type and cooridnates struct for schema.org spatial info
type LocationPoly struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// LocationFeature is a simple type and cooridnates struct for schema.org spatial info
type LocationFeature struct {
	Type     string       `json:"type"`
	Geometry LocationPoly `json:"geometry"`
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

	log.Print("connected")

	var value1 int // we seem to be getting RESP back..
	var value2 []interface{}
	reply, err := redis.Values(c.Do("INTERSECTS", "p418", "LIMIT", "50000", "OBJECT", wktstring))
	// reply, err := redis.String(c.Do("INTERSECTS", "p418", "LIMIT", "50000", "OBJECT", wktstring))
	// reply, err := redis.Values(c.Do("SCAN", "p418"))
	if err != nil {
		fmt.Printf("Error in reply %v \n", err)
	}
	if _, err := redis.Scan(reply, &value1, &value2); err != nil {
		fmt.Printf("Error in scan %v \n", err)
	}

	fmt.Println(value1)


	// results, _ := tile38RespAsGeoJSON(value2)
	// response.Write([]byte(results))
	response.Write([]byte("Hi there"))
}

// this is a a lot of pointless work!
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

		log.Printf("%s  %s", val0, val1)

		// serialize val1 to a simple struct with just type
		// then move on to the next level of serilization seen below

		lt := &LocType{}
		err := json.Unmarshal([]byte(val1), lt)
		if err != nil {
			log.Print(err)
			return "", err
		}

		if lt.Type == "Point" {
			lp := &LocationPoint{}
			err := json.Unmarshal([]byte(val1), lp)
			if err != nil {
				log.Print(err)
				return "", err
			}
			if lp.Type == "Point" {
				log.Println("add point")
				cd := gj.Coordinate{gj.Coord(lp.Coordinates[0]), gj.Coord(lp.Coordinates[1])} // is this long lat..  vs lat long?
				props := map[string]interface{}{"URL": val0}
				newp := gj.NewPoint(cd)
				f = gj.NewFeature(newp, props, nil)
				fa = append(fa, f)
			}
		}

		// ugh..  must I deconstruct and restructure what I know if valid GeoJSON?
		// Can I always ensure my spatial index is made of geojson..  since we build it..
		if lt.Type == "Polygon" {
			log.Println("add poly")
			lp := &LocationPoly{}
			err := json.Unmarshal([]byte(val1), lp)
			if err != nil {
				log.Print(err)
				return "", err
			}
			log.Print(lp)
		}

		if lt.Type == "Feature" {
			log.Println("add feature")
			lf := &LocationFeature{}
			err := json.Unmarshal([]byte(val1), lf)
			if err != nil {
				log.Print(err)
				return "", err
			}
			log.Print(lf)
		}

	}

	fc := gj.FeatureCollection{Type: "FeatureCollection", Features: fa}
	gjstr, err := gj.Marshal(fc)
	if err != nil {
		log.Println(err)
	}

	// log.Println(gjstr)

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

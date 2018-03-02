from geojsonio import display
import requests

# URL = 'http://localhost:6789/api/v1/spatial/search/object'
URL = 'http://geodex.org/api/v1/spatial/search/object'

# http://get.iedadata.org/doi/315201  
# {"type":"Polygon","coordinates":[[[-16.91266,28.11008],[-16.91266,33.31179],[-9.24511,33.31179],[-9.24511,28.11008],[-16.91266,28.11008]]]}
# http://opencoredata.org/id/dataset/273af50d-bffe-4094-9c0a-cd2f900b5474  
# {"type":"Point","coordinates":[-21.09,11.82]}
#  https://www.bco-dmo.org/dataset/628710  
# {"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-76.3573,34.98],[-76.3573,34.98],[-76.3573,34.98],[-76.3573,34.98],[-76.3573,34.98]]]}}

data='''{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "type": "Polygon",
        "coordinates": [
          [
            [
              -95.97656249999999,
              23.241346102386135
            ],
            [
              -88.24218749999999,
              23.241346102386135
            ],
            [
              -88.24218749999999,
              29.53522956294847
            ],
            [
              -95.97656249999999,
              29.53522956294847
            ],
            [
              -95.97656249999999,
              23.241346102386135
            ]
          ]
        ]
      }
    }
  ]
}
'''

# PARAMS = {'geowithin':data}
PARAMS = {'geowithin':data, 'filter':"bco-dmo"}
r = requests.get(url = URL, params = PARAMS)
print(r.content)
display(r.content)  # calls to geojson.io and opens your browser to view it..  

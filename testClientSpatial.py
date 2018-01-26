from geojsonio import display
import requests

URL = 'http://localhost:6789/api/v1/spatial/search/test1'
# URL = 'http://geodex.org/api/v1/spatial/search/test1'

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
              -131.1328125,
              -4.214943141390639
            ],
            [
              126.91406249999999,
              -4.214943141390639
            ],
            [
              126.91406249999999,
              58.63121664342478
            ],
            [
              -131.1328125,
              58.63121664342478
            ],
            [
              -131.1328125,
              -4.214943141390639
            ]
          ]
        ]
      }
    }
  ]
}
'''

PARAMS = {'geowithin':data}
r = requests.get(url = URL, params = PARAMS)
print(r.content)
#display(r.content)  # calls to geojson.io and opens your browser to view it..  
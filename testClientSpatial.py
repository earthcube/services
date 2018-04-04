from geojsonio import display
import requests

URL = 'http://geodex.org/api/v1/spatial/search/object'
#URL = 'http://geodex.org/api/v1/spatial/search/object'

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
              -112.8515625,
              -29.535229562948444
            ],
            [
              85.4296875,
              -29.535229562948444
            ],
            [
              85.4296875,
              65.36683689226321
            ],
            [
              -112.8515625,
              65.36683689226321
            ],
            [
              -112.8515625,
              -29.535229562948444
            ]
          ]
        ]
      }
    }
  ]
}
'''

data2='''{
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
              2.4609375,
              -41.77131167976406
            ],
            [
              64.6875,
              -41.77131167976406
            ],
            [
              64.6875,
              -8.754794702435618
            ],
            [
              2.4609375,
              -8.754794702435618
            ],
            [
              2.4609375,
              -41.77131167976406
            ]
          ]
        ]
      }
    }
  ]
}
'''

# PARAMS = {'geowithin':data}
#PARAMS = {'geowithin':data2, 'filter':"opencore"}
PARAMS = {'geowithin':data2}

r = requests.get(url = URL, params = PARAMS)
print(r.content)
#display(r.content)  # calls to geojson.io and opens your browser to view it..  

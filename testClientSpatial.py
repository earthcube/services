from geojsonio import display
import requests

#URL = 'http://localhost:6789/api/v1/spatial/search/object'
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

data3= '''
{
	"type": "Feature",
	"id": 1,
	"properties": {
		"dummy": 0.0
	},
	"geometry": {
		"type": "Polygon",
		"coordinates": [
			[
				[-1.684286582529027, 59.758556912514905],
				[11.250000000675005, 65.090030158598523],
				[24.184286582393444, 59.75855691307941],
				[17.834431060727248, 52.583358975375795],
				[4.665568940078517, 52.583358975086604],
				[-1.684286582529027, 59.758556912514905]
			]
		]
	}
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
              -102.3046875,
              -27.371767300523032
            ],
            [
              72.421875,
              -27.371767300523032
            ],
            [
              72.421875,
              65.94647177615738
            ],
            [
              -102.3046875,
              65.94647177615738
            ],
            [
              -102.3046875,
              -27.371767300523032
            ]
          ]
        ]
      }
    }
  ]
}'''

# PARAMS = {'geowithin':data}
# PARAMS = {'geowithin':data2, 'filter':"opencore"}
PARAMS = {'geowithin':data3}

r = requests.get(url = URL, params = PARAMS)
print(r.content)
# display(r.content)  # calls to geojson.io and opens your browser to view it..  

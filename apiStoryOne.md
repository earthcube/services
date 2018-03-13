# API Story One : A simple step through of API connections

## About

This is a flow of searches through the various web services.  

## Flow

A typical flow might start with an organic search on a term like "oxygen isotope"

```
GET http://geodex.org/api/v1/textindex/searchset?q=oxygen&n=20&s=0
```

```
[
  {
    "or": [
      {
        "position": 0,
        "indexpath": "indexes/linkedearth.bleve",
        "score": 0.6151675310472035,
        "URL": "http://wiki.linked.earth/Arc-DevonIceCap.Fisher.1983"
      },
      {
        "position": 1,
        "indexpath": "indexes/linkedearth.bleve",
        "score": 0.5904801774831937,
        "URL": "http://wiki.linked.earth/Arc-GISP2.Grootes.1997"
      },
      {
        "position": 2,
        "indexpath": "indexes/linkedearth.bleve",
        "score": 0.5542932979263926,
        "URL": "http://wiki.linked.earth/DSDP552.Shackleton.1984"
      },
      {
        "position": 3,
        "indexpath": "indexes/linkedearth.bleve",
        "score": 0.486703710181719,
        "URL": "http://wiki.linked.earth/Arc-MackenzieDelta.Porter.2013"
      }
    ],
    "highscore": 0.6151675310472035,
    "index": "linkedearth"
  },
...
```

From this we might pull the URLs and make a "set" to search some of the other 
web services with. There are service which take a set of URIs via POST.  We have tried
to add "set" as a suffix to these calls.

One of these might be something like the spatial call.

```
POST http://geodex.org/api/v1/spatial/search/resourceset
content-type: application/x-www-form-urlencoded

body= ["<http://wiki.linked.earth/Arc-DevonIceCap.Fisher.1983>","<http://wiki.linked.earth/Arc-GISP2.Grootes.1997>","<http://wiki.linked.earth/DSDP552.Shackleton.1984>","<http://wiki.linked.earth/Arc-MackenzieDelta.Porter.2013>"]
```


You could also call for details on a service with:

```
GET http://geodex.org/api/dev/graph/details?r=http://wiki.linked.earth/Arc-DevonIceCap.Fisher.1983

```

```
{
  "S": "http://opencoredata.org/id/dataset/0007e994-ba7f-4c74-b954-c7c58998d9b9",
  "Aname": "",
  "Name": "",
  "URL": "http://opencoredata.org/id/dataset/0007e994-ba7f-4c74-b954-c7c58998d9b9",
  "Description": "Janus Core Summary for ocean drilling expedition 125 site 781 hole A",
  "Citation": "",
  "Datepublished": "",
  "Curl": "http://opencoredata.org/api/v1/documents/download/125_781A_JanusCoreSummary_oRmJmVxx.csv",
  "Keywords": "Leg Site Hole Core Core_type Top_depth_mbsf Length_cored Length_recovered Percent_recovered Curated_length Ship_date_time Core_comment Janus Core Summary DSDP, OPD, IODP, JanusCoreSummary",
  "License": ""
}
```

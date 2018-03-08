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

body= ["<http://opencoredata.org/id/dataset/0007e994-ba7f-4c74-b954-c7c58998d9b9>","<http://opencoredata.org/id/dataset/0018797c-ce5c-46f2-a4a9-3a082954b875>","<http://opencoredata.org/id/dataset/001f4c51-0d7c-4cfe-bdaf-acb80e368f79>","<http://opencoredata.org/id/dataset/00878afc-68d4-438c-bfe2-c5bfec0d6c0d>","<http://opencoredata.org/id/dataset/00d46a42-b646-48c1-b159-1b0c821ca96f>","<http://opencoredata.org/id/dataset/0128ac30-67f3-4d90-85f5-03b111eb3d13>","<http://opencoredata.org/id/dataset/013025ee-651c-47d2-83bd-e539690bb679>","<http://opencoredata.org/id/dataset/0166f8a2-cfca-4d53-a37d-df6c619c6467>","<http://opencoredata.org/id/dataset/0191c9ae-d1f1-443b-96a2-9e355471135c>","<http://opencoredata.org/id/dataset/019c7064-0bce-465f-b6a5-afda259be8d6>","<http://opencoredata.org/id/dataset/01a67338-2598-4f55-bedb-84e09a3ca8dc>","<http://opencoredata.org/id/dataset/020c342b-b7f2-41df-b084-e59e9b25d2c8>","<http://opencoredata.org/id/dataset/022fa324-7083-4c24-89e4-687dbfd6dc34>","<http://opencoredata.org/id/dataset/028d9133-db6f-49fd-b085-ae012a4e1f3e>","<http://opencoredata.org/id/dataset/03399c22-fbdc-4314-b65e-93bae2c63f89>","<http://opencoredata.org/id/dataset/0378f88d-2d40-4ba3-86b9-07b48cda3de4>","<http://opencoredata.org/id/dataset/037ccc68-b6bf-44fa-b41e-0649aab68ad1>","<http://opencoredata.org/id/dataset/03ac8ab1-5478-451a-9d70-e14b541f3f68>","<http://opencoredata.org/id/dataset/03bc7b5f-13af-4630-a053-4eb41db5f062>","<http://opencoredata.org/id/dataset/03becfcf-4d43-4aef-97b5-60ea7d4d4564>"]
```


You could also call for details on a service with:

```
GET http://geodex.org/api/dev/graph/details?r=http://opencoredata.org/id/dataset/0007e994-ba7f-4c74-b954-c7c58998d9b9

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

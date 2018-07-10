## SHACL call tests

### About
Some example calls for SHACL




### Example calls
```
http --form POST localhost:6789/api/beta/shacl/testpost  body@datasetDataEx1.json
http -f POST localhost:6789/api/beta/shacl/testpost  body=@datasetDataEx1.json
http -f POST localhost:6789/api/beta/shacl/testpost  body=@ example1.ttl
http -f POST localhost:6789/api/beta/shacl/testpost  body=@example1.ttl
http -f POST localhost:6789/api/beta/shacl/testpost  body=@datasetDataEx1.json
http -f POST localhost:6789/api/beta/shacl/testpost  body=@datasetDataEx2.json
http -f POST localhost:6789/api/beta/shacl/testpost  body=@example1.nq
http -f POST localhost:6789/api/beta/shacl/igsn  body=@example1.nq
http -f POST localhost:6789/api/beta/shacl/eval  body=@example1.nq
```



http -f POST localhost:6789/api/beta/shacl/eval  datagraph=@pg1DataGraph.ttl shapegraph=@pg1ShapeGraph.ttl
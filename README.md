### P418 services

#### About
The services repo holds the code that exposed the indexes to 3rd party clients
on the net.

It follows the basic RESTful patterns and exposed swagger based descriptor documents.

Docker ENV variables:
	"GEODEX_BASEURL","http://geodex.org" 
	"GEODEX_PORT","6789" 
	blasturl "http://blast:10002/rest/_search"
	SPARQL:
	"SPARQL_EC", "http://geodex.org/blazegraph/namespace/p418/sparql"
	("SPARQL_rwg2", "http://geodex.org/blazegraph/namespace/p418/sparql")
	("SERVICES_LOGFILE","./runtime/log/serviceslog.txt" )
	INDEXCATALOG

push container: 
```
# change version
make docker
docker login docker.io -u USER -p PASS
docker push geodex/p418services:2.0.6
```


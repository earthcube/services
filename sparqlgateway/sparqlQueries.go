package sparqlgateway

const queries = `
# Comments are ignored, except those tagging a query.

# tag: LogoCall
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?graph ?type ?resource ?logo
WHERE {
  VALUES ?resource {
    <{{.RESID}}>
  }
  VALUES ?type {
    schema:Organization
    schema:Dataset
  }
  GRAPH ?graph {
    ?resource rdf:type ?type .
    OPTIONAL { ?resource schema:logo [ schema:url ?logo ] }
    
  }
}
ORDER BY ?graph ?type ?resource ?logo


# tag: ResourceResults
prefix schema: <http://schema.org/>
SELECT ?val ?desc ?pubname ?puburl
WHERE
{
  BIND(<{{.RESID}}> AS ?ID)
  ?ID schema:publisher ?pub .
  ?pub schema:name ?pubname .
  ?pub schema:url ?puburl .
  ?ID schema:variableMeasured ?res  .
  ?res a schema:PropertyValue .
  ?res schema:value ?val   .
  ?res schema:description ?desc     
} 

# tag: ResourceSetResults
prefix schema: <http://schema.org/>
SELECT DISTINCT ?val ?desc ?pubname ?puburl
WHERE
{
VALUES ?ID
{  {{.}}
}
?ID schema:variableMeasured ?res .
OPTIONAL {
?res a schema:PropertyValue .
?res schema:value ?val .
?res schema:description ?desc
}
OPTIONAL {
?ID schema:publisher ?pub .
OPTIONAL { ?pub schema:name ?pubname }
OPTIONAL { ?pub schema:url ?puburl }
}
}

# tag: ResourceSetPeople
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?g ?person (IF(?role = schema:contributor, "Contributor", IF(?role = schema:creator, "Creator/Author", "Author")) as ?rolename) ?name ?url ?orcid
WHERE
{
  GRAPH ?g {
   VALUES ?dataset
   {  {{.}}
   }
   VALUES ?role
   {
    schema:author
    schema:creator
    schema:contributor
   }
   { ?dataset ?role ?person }
   OPTIONAL {
    ?person a schema:Person .
    OPTIONAL { ?person schema:name ?name }
    OPTIONAL { ?person schema:url ?url }
    OPTIONAL { 
      ?person schema:identifier ?id .
      ?id schema:propertyID "datacite:orcid" .
      ?id schema:value ?orcid
    }
   }
  }
}


`

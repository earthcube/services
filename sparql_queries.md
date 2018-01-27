<a id="top"></a>
## Table of Contents ##
* [Repository](#repository)
  * [Services](#repository-services)
    * [Target Properties](#repository-services-target)
    * [Query Inputs](#repository-services-inputs)
* [Dataset](#dataset)
  * [Variables](#dataset-variables)
  * [People](#dataset-people)
  * [Funder](#dataset-funder)
  * [Publisher/Provider](#dataset-publisher_provider)
 
<hr/>
<a id="repository"></a>
## Repository ##

<a id="repository-services"></a>
### Repository Services ###

1. ?action_target_type = "@type" means you've got to query down to check for schema:urlTemplate, schema:contentType, and schema:httpMethod. See [Repository Service Target Properties](#repository-services-target)
2. ?action_query_input_type = "@type" means you've got to query down for all the input parameters to a search service. See [Repository Service Query Inputs](#repository-services-inputs)

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?repo ?service ?type ?name ?desc ?service_url ?channel_url 
(IF(isLiteral(?action_target), ?action_target , "@type") as ?action_target_type) 
(IF(isLiteral(?action_query_input),"", "@type") as ?action_query_input_type)
WHERE
{
  VALUES ?type
  {
    "gdx:SearchService"
    "gdx:SubmissionService"
    "gdx:SyndicationService"
  }
  ?repo schema:additionalType "gdx:ResearchRepositoryService" .
  ?repo schema:availableChannel ?channel .
  ?channel schema:providesService ?service .
  OPTIONAL { ?channel schema:serviceUrl ?channel_url }
  ?service schema:additionalType ?type .
  OPTIONAL { ?service schema:name ?name }
  OPTIONAL { ?service schema:description ?desc }
  OPTIONAL { ?service schema:url ?service_url }
  OPTIONAL { 
    ?service schema:potentialAction ?action .
    OPTIONAL { ?action schema:target ?action_target }
    OPTIONAL { 
      ?action a schema:SearchAction .
      ?action schema:query-input ?action_query_input .
    }
  }  
}
ORDER BY ?repo ?type ?service
```

<a id="repository-services-target"></a>
#### Repository Services Target Properties
1. HTTP Method is multi-valued, comma-separated list

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?service ?url_template ?content_type (GROUP_CONCAT(DISTINCT ?http_method ; separator=",") as ?http_method_list)
WHERE
{
  VALUES ?service { <http://lod.bco-dmo.org/sparql> }
  ?service schema:potentialAction [ schema:target ?target ] .
  FILTER EXISTS { ?target rdf:type schema:EntryPoint }
  OPTIONAL { ?target schema:urlTemplate ?url_template }
  OPTIONAL { ?target schema:contentType ?content_type }
  OPTIONAL { ?target schema:httpMethod ?http_method }
 
}
GROUP BY ?service ?url_template ?content_type
ORDER BY ?service
```

<a id="repository-services-inputs"></a>
#### Repository Services Query Inputs ####

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?service ?name ?required ?default ?pattern ?multiple_values ?read_only ?max_value ?min_value ?max_length ?min_length
WHERE
{
  VALUES ?service { <http://lod.bco-dmo.org/sparql> }
  ?service schema:potentialAction ?action .
  ?action a schema:SearchAction .
  ?action schema:query-input ?input .
  ?input a schema:PropertyValueSpecification .
  ?input schema:valueName ?name
  OPTIONAL { ?input schema:valueRequired ?required }
  OPTIONAL { ?input schema:defaultValue ?default }
  OPTIONAL { ?input schema:valuePattern ?pattern }
  OPTIONAL { ?input schema:multipleValues ?multiple_values }
  OPTIONAL { ?input schema:readonlyValue ?read_only }
  OPTIONAL { ?input schema:maxValue ?max_value }
  OPTIONAL { ?input schema:minValue ?min_value }
  OPTIONAL { ?input schema:valueMaxLength ?max_length }
  OPTIONAL { ?input schema:valueMinLength ?min_length }
  
}
ORDER BY ?service ?name
```

<hr/>
<a id="dataset"></a>
## Dataset ##

<a id="dataset-variables"></a>
### Dataset Variables ###

1. I'm not grouping by variable name becuase we can show the need for understanding semantic sameness of variables
```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?g ?value ?name ?desc ?url ?units ?var
WHERE
{
  GRAPH ?g {
  VALUES ?dataset
   {
    <https://www.bco-dmo.org/dataset/3300>
    <http://opencoredata.org/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816>
   }
   ?dataset schema:variableMeasured ?var  .
   OPTIONAL {
    ?var a schema:PropertyValue .
    OPTIONAL { ?var schema:value ?value }
    OPTIONAL { ?var schema:description ?desc }
    OPTIONAL { ?var schema:name ?name }
    OPTIONAL { ?var schema:url ?url }
    OPTIONAL { ?var schema:unitText ?units }
   }
  }
}
ORDER BY lcase(?value) lcase(?name) ?g
```

<a id="dataset-people"></a>
### Dataset People ###

1. Purposefuly, I do not group by ORCiD to demonstrate the NEED for persistent IDs.
2. I do not order the results because there is no guarantee that all person names will get broken out into familyName and givenName, etc. We can improve this query to try to construct name from all fields, for comparison if we have time.

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?g ?person (IF(?role = schema:contributor, "Contributor", IF(?role = schema:creator, "Creator/Author", "Author")) as ?rolename) ?name ?url ?orcid
WHERE
{
  GRAPH ?g {
   VALUES ?dataset
   {
    <https://www.bco-dmo.org/dataset/472032>
    <http://opencoredata.org/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816>
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
```
<a id="dataset-funder"></a>
### Dataset Funder ###

1. Needs testing

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?g ?funder ?legal_name ?name ?url ?award ?award_name ?award_url
WHERE
{
  GRAPH ?g {
   VALUES ?dataset
   {
    <https://www.bco-dmo.org/dataset/472032>
    <http://opencoredata.org/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816>
   }
   ?dataset schema:funder ?funder
   OPTIONAL {
    ?funder a schema:Organization .
    OPTIONAL { ?funder schema:legalName ?legal_name }
    OPTIONAL { ?funder schema:name ?name }
    OPTIONAL { ?funder schema:url ?url }
    OPTIONAL { 
      ?funder schema:makesOffer ?award .
      ?award schema:additionalType "geolink:Award" .
      ?award schema:name ?award_name .
      OPTIONAL { ?award schema:url ?award_url }
    }
   }
  }
}
```

<a id="dataset-publisher_provider"></a>
### Dataset Publisher/Provider ###

```
PREFIX schema: <http://schema.org/>
SELECT DISTINCT ?g ?org ?type (IF(?role = schema:publisher, "Publisher", "Provider") as ?rolename) ?legal_name ?name ?url
WHERE
{
  GRAPH ?g {
   VALUES ?dataset
   {
    <https://www.bco-dmo.org/dataset/472032>
    <http://opencoredata.org/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816>
   }
   VALUES ?role { schema:publisher  schema:provider }
   VALUES ?type { schema:Organization schema:Person }
   ?dataset ?role ?org .
   ?org a ?type .
   OPTIONAL { ?org schema:legalName ?legal_name }
   OPTIONAL { ?org schema:name ?name }
   OPTIONAL { ?org schema:url ?url }
  }
}
ORDER BY ?role ?legal_name ?name

```


Back to [Top](#top)

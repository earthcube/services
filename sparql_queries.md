<a id="top"></a>
* [Dataset](#dataset)
  * [Variables](#dataset-variables)
  
<a id="dataset"></a>
# Dataset #

<a id="dataset-variables"></a>
## Dataset Variables ##
```
PREFIX schema: <http://schema.org/>
SELECT ?g ?value ?name ?desc ?url ?units ?var
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
[Back to Top](#top)

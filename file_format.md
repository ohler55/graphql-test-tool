# Use Case File Format

Each use case file should contain only on JSON object. That object can have a "comment" and a "steps" element.

 - **comment** is an optional description of the use case.
 - **steps** is an array of steps that define the use case.

_Note: String fields such as comments and content can also be an array of strings. The array of strings is joined with newlines to for a string. The intent is to make it easier to enter multi-line comments more easily._

An example is:

```json
{
  "comment": [
    "This example demonstrates the use of multi-line comments (this comment), a GET",
    "query, sortBy, use of exact matching (artists.0.name), and the use of a regular",
    "expression match (artists.1.name)."
  ],
  "steps": [
    {
      "label": "List the artists",
      "path": "?query={artists{name}}",
      "sortBy": {
	"artists": "name"
      },
      "expect": {
	"data": {
	  "artists": [
	    {
              "name": "Fazerdaze"
	    },
	    {
              "name": "/Boys/"
	    }
	  ]
	}
      }
    }
  ]
}
```

The "steps" array contains Step objects that can include the following fields:

 - **label** is an optional label that will be displayed during a run
   if comments are being dislayed.

 - **comment** is an optional description of the step.

 - **path** is the path part of the URL that will be used in the
   request to the GraphQL server. The path is assumed to be relative
   to the base URL unless it starts with a '/' character.

 - **content** is the content when using a POST request. If the
   Content is not empty then a POST is used. If it is empty then a GET
   request is made to the server. The string should be either GraphQL
   or the HTTP GET format.

 - **json** if true and a POST request then the JSON format will be
   used with a COntent-Type of "applicaiton/json" otherwise the
   Content will be sent as GraphQL with a Content-Type of
   "application/graphql".

 - **remember** in a series of steps it is often helpful to remember
   values from earlier steps. The Remember map describes what to
   remember and what key to store that value in. In the map, the keys
   are the keys for the memory cache while the Remember map values are
   the path to the value to remember. The path is a simple dot
   delimited path. It is not a full JSON path (maybe in the future).

 - **op** is the operation to include in either the URL query or as a
   value for the 'operationName' if using JSON in the Content.

 - **vars** are the variables to be passed along with a GraphQL
   request. They are appended to the URL or added to the JSON
   'variables' element if the POST contents is JSON. The values in the
   Vars map can be either a value or a string that begines with a '$'
   which indicates a rmembered value should be used instead.

 - **sortBy** are the sort keys for the result. Depending on the
   implementation of the GraphQL server, the order of returned objects
   may not be consistent. To make it easier to use in testing the
   SortBy keys are the paths to arrays while the value for the keys
   are the attributes to sort on.

 - **expect** is the expected contents. Like the Content it can be a
   nil, string, array, or map. A nil value indicates no checking of
   the response is needed. Strings are used mostly for fetched HTML
   while the normal JSON responses are converted to a genereric
   map[string]interface{} and compared to the Expect value.

   The rules for comparison are:

    1) Any element in the Expect value must be present in the response.

    2) Elements in the response not in the Expect value are ignored unless a
       "*": null key value pair is present in which case any key not specified
       must be null or no present.

    3) If the Expect value is a string that starts and ends with a '/'
       character the string is assumed to be a regular expression is
       matched against the string in the reponse.

    4) Maps and arrays are followed recursively.

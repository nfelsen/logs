package main

import (
	"fmt"
  "github.com/mattbaird/elastigo/api"
  "github.com/mattbaird/elastigo/core"
)


api.Domain = "localhost"


type Tweet struct {
  User     string    `json:"user"`
  Message  string    `json:"message"`
}

// Set the Elasticsearch Host to Connect to
api.Domain = "localhost"
// api.Port = "9300"

// add single go struct entity
response, _ := core.Index("twitter", "tweet", "1", nil, Tweet{"kimchy", "Search is cool"})

// you have bytes
tw := Tweet{"kimchy", "Search is cool part 2"}
bytesLine, err := json.Marshal(tw)
response, _ := core.Index("twitter", "tweet", "2", nil, bytesLine)

// Bulk Indexing
t := time.Now()
core.IndexBulk("twitter", "tweet", "3", &t, Tweet{"kimchy", "Search is now cooler"})

// Search Using Raw json String
searchJson := `{
    "query" : {
        "term" : { "user" : "kimchy" }
    }
}`
out, err := core.SearchRequest(true, "twitter", "tweet", searchJson, "")
if len(out.Hits.Hits) == 1 {
  fmt.Println(string(out.Hits.Hits[0].Source))
}

// func main() {
//   var (
//     es_host        = config.String("elasticsearch.host", "undefined")
//     es_port        = config.Int("elasticsearch.port", 9200)
//     es_max_pending = config.Int("elasticsearch.max_pending", 1000000)
//     in_port        = config.Int("in.port", 2003)
//   )
// }


// go fun() {
//   searchJson := `{
//       "query": { 
//         "match_all": { } 
//       },
//       "facets": {
//         "tags": {
//           "terms": {
//             "field": "tags",
//             "all_terms": true
//           }
//         }
//       } 
//   }`
//   out, err := core.SearchRequest(true, "twitter", "tweet", searchJson, "")
//   if len(out.Hits.Hits) == 1 {
//     fmt.Println(string(out.Hits.Hits[0].Source))
//   }
// }


// func main() {
// 	fmt.Println("Hello")
// }

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"os"
)

var (
	host *string = flag.String("host", "localhost", "Elasticsearch Host")
)

func main() {
	c := elastigo.NewConn()
	var results = make(map[string]bool)
	log.SetFlags(log.LstdFlags)
	flag.Parse()

	// fmt.Println("host = ", *host)
	// Set the Elasticsearch Host to Connect to
	c.Domain = *host

	// Search Using Raw json String
	searchJson := `{
      "query": {
        "match_all": { }
      },
      "facets": {
        "tags": {
          "terms": {
            "field": "tags",
            "all_terms": true
          }
        }
      }
  }`
	out, err := c.Search("logstash-2014.09.03", "logs", nil, searchJson)
	if len(out.Hits.Hits) != 0 {
		for i := 0; i < len(out.Hits.Hits); i++ {
			var fields = out.Hits.Hits[i].Source
			c := make(map[string]interface{})
			err := json.Unmarshal(*fields, &c)
			exitIfErr(err)

			// copy c's keys into k
			for s, _ := range c {
				results[s] = true
			}
		}
	}
	for key, _ := range results {
		fmt.Println(key)
	}
	exitIfErr(err)
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

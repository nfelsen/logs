package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"os"
)

type Configuration struct {
	Host string `json:"host"`
	Type string `json:"type"`
}

var (
	ConfigFile    *string = flag.String("config", os.Getenv("HOME")+"/.logs_config.json", "Logs configuration")
	configuration         = Configuration{}
)

func LoadConfig() {
	file, _ := os.Open(*ConfigFile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	LoadConfig()
	bolB, _ := json.Marshal(configuration)
	fmt.Println(string(bolB))
	c := elastigo.NewConn()
	var results = make(map[string]bool)
	log.SetFlags(log.LstdFlags)
	flag.Parse()

	// fmt.Println("host = ", configuration)
	// Set the Elasticsearch Host to Connect to
	c.Domain = configuration.Host

	// Search Using Raw json String
	// searchJson := `{
	//      "query": {
	//        "match_all": { }
	//      },
	//      "facets": {
	//        "tags": {
	//          "terms": {
	//            "field": "tags",
	//            "all_terms": true
	//          }
	//        }
	//      }
	//  }`
	searchJson := `{
	  "query": {
	    "match_all": {}
	  },
	  "size": 1,
	  "sort": [
	    {
	      "_timestamp": {
	        "order": "desc"
	      }
	    }
	  ]
	}`
	out, err := c.Search("logstash-2015.08.01", "logs", nil, searchJson)
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

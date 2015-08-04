package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v2"
	"os"
	"regexp"
	s "strings"
	"time"
)

type Configuration struct {
	Host      string `json:"host"`
	ServerURL string `json:"server_url"`
	Type      string `json:"type"`
	Service   map[string]string
}

var (
	ConfigFile    *string = flag.String("config", os.Getenv("HOME")+"/.logs_config.json", "Logs configuration")
	configuration         = Configuration{}
	arg1                  = "sepia"
)

func LoadConfig() {
	file, _ := os.Open(*ConfigFile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

// TODO fix bug on index names like Neo4j-YYYY.MM.DD
func ListIndexes(client *elastic.Client, t time.Time) map[string]struct{} {
	set := make(map[string]struct{})
	indexes, err := client.IndexNames()
	if err != nil {
		panic(err)
	}
	for _, element := range indexes {
		r, _ := regexp.Compile("(.*-)[0-9]{4}(\\..*)")
		if r != nil {
			element = r.ReplaceAllString(element, "${1}2006${2}")
		}
		r, _ = regexp.Compile("(.*-2006.)[0-9]{2}(\\.?.*)")
		if r != nil {
			element = r.ReplaceAllString(element, "${1}01${2}")
		}
		r, _ = regexp.Compile("(.*-2006.01.)[0-9]{2}(\\.?.*)")
		if r != nil {
			element = r.ReplaceAllString(element, "${1}02${2}")
		}
		r, _ = regexp.Compile("(.*-2006.01.02.)[0-9]{2}(\\.?.*)")
		if r != nil {
			element = r.ReplaceAllString(element, "${1}15${2}")
		}
		index := fmt.Sprintf("%s", t.Format(element))
		set[index] = struct{}{}
	}
	return set
}

func ReplaceTime(service *string, t time.Time) {
	year := fmt.Sprintf("%04d", t.UTC().Year())
	month := fmt.Sprintf("%02d", t.UTC().Month())
	day := fmt.Sprintf("%02d", t.UTC().Day())
	hour := fmt.Sprintf("%02d", t.UTC().Hour())
	min := fmt.Sprintf("%02d", t.UTC().Minute())
	*service = s.Replace(*service, "YYYY", year, -1)
	*service = s.Replace(*service, "MM", month, -1)
	*service = s.Replace(*service, "DD", day, -1)
	*service = s.Replace(*service, "HH", hour, -1)
	*service = s.Replace(*service, "mm", min, -1)
}

func TailLog(client *elastic.Client, index string) {
	_, err := client.IndexExists(index).Do()
	if err != nil {
		panic(err)
	}
	// termQuery := elastic.NewTermQuery("service", "medium")
	// Query(&termQuery). // specify the query
	// Must(elastic.From("@timestamp").
	searchResult, err := client.Search().
		Index(index).
		Sort("@timestamp", false).
		From(0).Size(100). // take documents 0-9
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	if searchResult.Hits != nil {
		fmt.Printf("Found a total of %d logs\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var t map[string]interface{}
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				panic(err)
			}
			for key, element := range t {
				fmt.Printf("%s: %s, ", key, element)
			}
			println("\n")
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
}

func main() {
	LoadConfig()
	bolB, _ := json.Marshal(configuration)
	fmt.Println(string(bolB))
	client, err := elastic.NewClient(elastic.SetURL(configuration.ServerURL))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping().Do()
	if err != nil {
		panic(err)
	}
	t := time.Now()
	fmt.Printf("%s\n", t.Format("logstash-2006.01.02.03"))
	service := configuration.Service[arg1]
	ReplaceTime(&service, t)
	fmt.Println("service : ", service)
	ListIndexes(client, t)
	TailLog(client, service)
	os.Exit(0)
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}

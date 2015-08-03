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
	arg1                  = "medium"
)

func LoadConfig() {
	file, _ := os.Open(*ConfigFile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func ListIndexes(client *elastic.Client, t time.Time) {
	set := make(map[string]struct{})
	indexes, err := client.IndexNames()
	if err != nil {
		panic(err)
	}
	for _, element := range indexes {
		res := element
		match, _ := regexp.MatchString("[0-9]{4}.[0-9]{2}.[0-9]{2}.[0-9]{2}", element)
		if match {
			r, _ := regexp.Compile("[0-9]{4}.[0-9]{2}.[0-9]{2}.[0-9]{2}")
			res = r.ReplaceAllString(element, "2006.01.02.15")
		} else {
			match, _ := regexp.MatchString("[0-9]{4}.[0-9]{2}.[0-9]{2}", element)
			if match {
				r, _ := regexp.Compile("[0-9]{4}.[0-9]{2}.[0-9]{2}")
				res = r.ReplaceAllString(element, "2006.01.02")
			}
		}
		index := fmt.Sprintf("%s", t.Format(res))
		set[index] = struct{}{}

	}
	// newSlice := []string{}
	// for key := range set {
	// 	newSlice := append(newSlice, key)
	// }
	fmt.Println(set)
}

func ReplaceTime(service *string, t time.Time) {
	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())
	hour := fmt.Sprintf("%02d", t.Hour())
	min := fmt.Sprintf("%02d", t.Minute())
	*service = s.Replace(*service, "YYYY", year, -1)
	*service = s.Replace(*service, "MM", month, -1)
	*service = s.Replace(*service, "DD", day, -1)
	*service = s.Replace(*service, "HH", hour, -1)
	*service = s.Replace(*service, "mm", min, -1)
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
	os.Exit(0)
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}

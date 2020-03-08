package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/TiersLieuxEdu/discourse-locations-extractor/pkg/forum"
	gj "github.com/kpawlik/geojson"
)

// MarshalIndent object to json string representation
func MarshalIndent(object interface{}, indent string) (data string, err error) {
	if t, err := json.MarshalIndent(object, "", indent); err != nil {
		data = ""
	} else {
		data = string(t)
	}
	return
}

func main() {
	newlineRegexp := regexp.MustCompile(`\r?\n`)

	topics := forum.GetTopics()
	sort.SliceStable(topics, func(i, j int) bool {
		return topics[i].Id < topics[j].Id
	})
	fc := gj.NewFeatureCollection([]*gj.Feature{})
	for _, value := range topics {
		log.Printf("%s...\n", value.Title)
		info := forum.GetInformations(value)

		if info.Lat == 0 || info.Long == 0 {
			log.Println("Coordinates not set. Skipping.")
			continue
		}
		p := gj.NewPoint(gj.Coordinate{gj.CoordType(info.Long), gj.CoordType(info.Lat)})

		theHTMLAdress := newlineRegexp.ReplaceAllString(info.Adresse, ", ")
		props := map[string]interface{}{
			"name":     info.Name,
			"forum":    info.Forum,
			"site":     info.WebSite,
			"tags":     info.Tags,
			"machines": info.Machines,
			"adresse":  theHTMLAdress,
		}
		for _, aTag := range info.Tags {
			props[strings.ToLower(aTag)] = true
		}
		f2 := gj.NewFeature(p, props, nil)
		fc.AddFeatures(f2)
		time.Sleep(1 * time.Second)
	}

	if gjstr, err := MarshalIndent(fc, "  "); err != nil {
		panic(err)
	} else {
		fmt.Println(gjstr)
	}
}

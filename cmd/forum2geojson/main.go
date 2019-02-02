package main

import (
  "fmt"
  "log"
  gj "github.com/kpawlik/geojson"
  "github.com/TiersLieuxEdu/discourse-locations-extractor/pkg/forum"
  "strconv"
)

func main() {

  topics := forum.GetTopics()
  fc := gj.NewFeatureCollection([]*gj.Feature {})
  for _, value := range topics {
    log.Printf("%s...\n", value.Title)
    info := forum.GetInformations(value)

    lat, errLat := strconv.ParseFloat(info.Latitude, 64)
    if errLat != nil {
      log.Printf("Cannot convert latitude '%s' to float\n", info.Latitude)
      continue
    }
    long, errLong := strconv.ParseFloat(info.Longitude, 64)
    if errLong != nil {
      log.Printf("Cannot convert Longitude '%s' to float\n", info.Longitude)
      continue
    }
    if lat == 0 || long == 0 {
      log.Println("Coordinates not set. Skipping.")
      continue
    }
    p := gj.NewPoint(gj.Coordinate{gj.CoordType(long), gj.CoordType(lat)})
    props := map[string]interface{}{"name": info.Name, "forum": info.Forum, "site": info.WebSite}
    f2 := gj.NewFeature(p, props, nil)
    fc.AddFeatures(f2)
  }

  if gjstr, err := gj.Marshal(fc); err != nil {
      panic(err)
  } else {
      fmt.Println(gjstr)
  }
}

package main

import (
  "encoding/csv"
  "log"
  "os"
  "github.com/TiersLieuxEdu/discourse-locations-extractor/tierslieuxedu/forum"
)

func main() {

  topics := forum.GetTopics()
  csvOutput := csv.NewWriter(os.Stdout)
  for _, value := range topics {
    log.Printf("%s...\n", value.Title)
    info := forum.GetInformations(value)
    //fmt.Printf("%s, %s, %s, %s, %s\n", info.Name, info.Latitude, info.Longitude, info.WebSite, info.Forum)
    if err := csvOutput.Write(info.AsSlice()); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
  }

  // Write any buffered data to the underlying writer (standard output).
	csvOutput.Flush()
	if err := csvOutput.Error(); err != nil {
		log.Fatal(err)
	}

}

package main

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "log"
  "net/http"
)

type LieuInfo struct {
  Name string
  Latitude string
  Longitude string
  //Tags []string
}


type Post struct {
  Cooked  string
  Wiki bool
}

type PostStream struct {
  Posts []Post
}

type Topic struct {
		Id    int
		Title string
    Unpinned bool
    PostStream *PostStream `json:"post_stream"`
}

func getTopics() []Topic {
  resp, err := http.Get("https://forum.tierslieuxedu.org/c/lieux.json")
  if err != nil {
  	log.Fatal(err)
  }
  defer resp.Body.Close()

  body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

  //json.NewDecoder(res.Body)

  var jsonRoot map[string]*json.RawMessage
  jsonRootErr := json.Unmarshal(body, &jsonRoot)
  if jsonRootErr != nil {
		log.Fatal(jsonRootErr)
	}

  var jsonTopicList map[string]*json.RawMessage
  jsonTopicListErr := json.Unmarshal(*jsonRoot["topic_list"], &jsonTopicList)
  if jsonTopicListErr != nil {
		log.Fatal(jsonTopicListErr)
	}

  var topics []Topic
  jsonTopicListTopicsErr := json.Unmarshal(*jsonTopicList["topics"], &topics)
  if jsonTopicListTopicsErr != nil {
    log.Fatal(jsonTopicListTopicsErr)
  }

  return topics
}

func extractInfo(html string, info *LieuInfo) {
  
}

func getInformations(topic Topic) LieuInfo {
  url := fmt.Sprintf("https://forum.tierslieuxedu.org/t/%d.json", topic.Id)
  fmt.Printf("%s\n", url)
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

  var topicEnhance Topic
  jsonRootErr := json.Unmarshal(body, &topicEnhance)
  if jsonRootErr != nil {
		log.Fatal(jsonRootErr)
	}

  stream := topicEnhance.PostStream
  if (stream == nil) {
    log.Fatal("There is a null")
  }
  posts := stream.Posts

  var info LieuInfo
  info.Name = topic.Title
  info.Latitude = "0"
  info.Longitude = "0"

  for _, p := range posts {
    fmt.Printf("*");
    if (p.Wiki) {
      fmt.Printf("%s", p.Cooked)
      extractInfo(p.Cooked, &info)
    }
  }

  return info
}

func main() {

  topics := getTopics()

  for _, value := range topics {
    fmt.Printf("%s...\n", value.Title)
    info := getInformations(value)
    fmt.Printf("(%s, %s)\n", info.Latitude, info.Longitude)
  }

}

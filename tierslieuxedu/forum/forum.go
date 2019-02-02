package forum;

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "golang.org/x/net/html"
  "log"
  "net/http"
  "strings"
  "github.com/TiersLieuxEdu/discourse-locations-extractor/tierslieuxedu/lieux"
)

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

func GetTopics() []Topic {
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

func extractInfo(htmlSrc string, info *lieux.Info) {
  htmlReader := strings.NewReader(htmlSrc)
  foundDefinitions := make(map[string]string)
  z := html.NewTokenizer(htmlReader)
  inDL := false
  inDT := false
  inDD := false
  var lastKey string
  for {
      tt := z.Next()

      switch {
      case tt == html.ErrorToken:
        log.Printf("data: %s\n", foundDefinitions)
      	// End of the document, we're done
        if val, ok := foundDefinitions["Latitude"]; ok {
          info.Latitude = val
        }
        if val, ok := foundDefinitions["Longitude"]; ok {
          info.Longitude = val
        }
        if val, ok := foundDefinitions["Site"]; ok {
          info.WebSite = val
        }
        return
      case tt == html.StartTagToken:
        t :=  strings.TrimSpace(z.Token().Data)
        isDefinitionList := t == "dl"
        if isDefinitionList {
            inDL = true
        }
        isDefinitionTerm := inDL && t == "dt"
        if isDefinitionTerm {
          inDT = true
        }
        isDefinitionData := inDL && t == "dd"
        if isDefinitionData {
          inDD = true
        }
        break
      case tt == html.EndTagToken:
        t :=  strings.TrimSpace(z.Token().Data)
        isDefinitionList := t == "dl"
        if isDefinitionList {
            inDL = false
        }
        isDefinitionTerm := inDL && t == "dt"
        if isDefinitionTerm {
          inDT = false
        }
        isDefinitionData := inDL && t == "dd"
        if isDefinitionData {
          inDD = false
        }
        break
      case tt == html.TextToken:
        t :=  strings.TrimSpace(z.Token().Data)
        if inDT {
          lastKey = t
        } else if inDD {
          foundDefinitions[lastKey] = t
        }
        break
      }
  }
}

func GetInformations(topic Topic) lieux.Info {
  forumUrl := fmt.Sprintf("https://forum.tierslieuxedu.org/t/%d", topic.Id)
  url := fmt.Sprintf("%s.json", forumUrl)
  log.Printf("%s\n", url)
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

  var info lieux.Info
  info.Name = topic.Title
  info.Forum = forumUrl
  info.Latitude = "0"
  info.Longitude = "0"

  for _, p := range posts {
    log.Printf("Posts %v\n", p.Wiki);
    if (p.Wiki) {
      //fmt.Printf("%s\n", p.Cooked)
      extractInfo(p.Cooked, &info)
    }
  }

  return info
}

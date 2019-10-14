package forum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/TiersLieuxEdu/discourse-locations-extractor/pkg/lieux"
	"golang.org/x/net/html"
)

type Post struct {
	Cooked string
	Wiki   bool
}

type PostStream struct {
	Posts []Post
}

type Topic struct {
	Id         int
	Title      string
	Unpinned   bool
	PostStream *PostStream `json:"post_stream"`
}

func GetTopicsForPage(page int) []Topic {
	log.Printf("Getting Page %d...", page)
	url := fmt.Sprintf("https://forum.tierslieuxedu.org/c/lieux.json?page=%d", page)
	resp, err := http.Get(url)
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

func GetTopics() []Topic {
	pageIndex := 0
	var topics []Topic
	for {
		topicsToAppend := GetTopicsForPage(pageIndex)
		if len(topicsToAppend) == 0 {
			return topics
		}
		pageIndex += 1
		topics = append(topics, topicsToAppend...)
	}
}

func ConvertPositionning(pos string) (float64, error) {
	if strings.IndexRune(pos, '°') != -1 {
		f := func(c rune) bool {
			return c == '°' || c == '\'' || c == '"'
		}
		splitted := strings.FieldsFunc(pos, f)
		var result float64
		if len(splitted) >= 3 {
			for i, v := range splitted {
				converted, err := strconv.ParseFloat(v, 64)
				if err != nil && i != 3 {
					break
				}
				switch i {
				case 0:
					result = converted
					break
				case 1:
					result += converted / 60.0
					break
				case 2:
					result += converted / (60.0 * 60.0)
					break
				case 3:
					if v == "S" || v == "W" || v == "O" {
						result = -result
					}
				}
			}
			return result, nil
		}
	}
	return strconv.ParseFloat(pos, 64)
}

func ConvertTags(rawList string, lieu *lieux.Info) error {
	if lieu.Tags == nil {
		lieu.Tags = make([]string, 0, 5)
	}
	f := func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	}
	splitted := strings.FieldsFunc(rawList, f)
	for _, v := range splitted {
		if strings.HasPrefix(v, "#") {
			v = v[1:]
		}
		v = strings.TrimSpace(v)
		if len(v) != 0 {
			lieu.Tags = append(lieu.Tags, v)
		}
	}
	return nil
}

func ConvertMachines(rawList string, lieu *lieux.Info) error {
	if lieu.Machines == nil {
		lieu.Machines = make([]string, 0, 5)
	}
	f := func(c rune) bool {
		return c == ',' || c == '/'
	}
	splitted := strings.FieldsFunc(rawList, f)
	for _, v := range splitted {
		v = strings.TrimSpace(v)
		if len(v) != 0 {
			lieu.Machines = append(lieu.Machines, v)
		}
	}
	return nil
}

func ExtractInfo(htmlSrc string, info *lieux.Info) {
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
				info.Latitude = strings.Replace(val, ",", ".", -1)
				if lat, err := ConvertPositionning(info.Latitude); err == nil {
					info.Lat = lat
				} else {
					log.Printf("Cannot convert Latitude '%s' to float\n", info.Latitude)
					info.Lat = 0
				}
			}

			if val, ok := foundDefinitions["Longitude"]; ok {
				info.Longitude = strings.Replace(val, ",", ".", -1)
				if long, err := ConvertPositionning(info.Longitude); err == nil {
					info.Long = long
				} else {
					log.Printf("Cannot convert Longitude '%s' to float\n", info.Longitude)
					info.Long = 0
				}
			}
			if val, ok := foundDefinitions["Site"]; ok {
				info.WebSite = val
			}

			if val, ok := foundDefinitions["Tags"]; ok {
				ConvertTags(val, info)
			}

			if val, ok := foundDefinitions["Machines"]; ok {
				ConvertMachines(val, info)
			}

			if val, ok := foundDefinitions["Adresse"]; ok {
				info.Adresse = val
			}

			return
		case tt == html.StartTagToken:
			t := strings.TrimSpace(z.Token().Data)
			//log.Printf("start tok: %s\n", t)
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
			t := strings.TrimSpace(z.Token().Data)
			//log.Printf("end tok: %s\n", t)
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
		case tt == html.SelfClosingTagToken:
			t := strings.TrimSpace(z.Token().Data)
			//log.Printf("self closing tok: %s\n", t)
			if t == "br" {
				foundDefinitions[lastKey] += "\n"
				break
			}
			break
		case tt == html.TextToken:
			t := strings.TrimSpace(z.Token().Data)
			if inDT {
				lastKey = strings.Title(strings.ToLower(t))
			} else if inDD {
				foundDefinitions[lastKey] += t
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
	if stream == nil {
		log.Fatal("There is a null")
	}
	posts := stream.Posts

	var info lieux.Info
	info.Name = topic.Title
	info.Forum = forumUrl
	info.Latitude = "0"
	info.Longitude = "0"

	for _, p := range posts {
		log.Printf("Posts %v\n", p.Wiki)
		if p.Wiki {
			//fmt.Printf("%s\n", p.Cooked)
			ExtractInfo(p.Cooked, &info)
			//log.Printf("Found %v\n", topic.Id)
			//break;
		}
	}

	return info
}

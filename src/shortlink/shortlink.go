package shortlink

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"github.com/marksalpeter/token"
	"log"
	"regexp"
	"github.com/FedericoPonzi/distributed-systems-bot/src/repository"
	"net/http"
	"io/ioutil"
	"fmt"
)

type ShortlinkService struct {
	repo *repository.MysqlRepository
	url  string
}


func NewShortLinkService(repo *repository.MysqlRepository) (*ShortlinkService) {

	return &ShortlinkService{repo: repo, url: "https://ds.fponzi.me"}
}

func (service *ShortlinkService) GenerateShortlinkFromLink(link string) (string) {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	r := regexp.MustCompile(`(?i)<\s*title\s*>\s*(.+)\s*<\s*/title\s*>`)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error on reding response body.")
	}
	re := r.FindSubmatch(body)
	if re == nil {
		log.Fatal("Impossibile to find title in the web page. Please supply a custom one.")
	}
	fmt.Println("Found: '" + string(re[1]) + "'")
	id := service.generateId(string(re[1]))
	service.repo.AddShortlink(id, link)
	return service.url + "/" + id
}
func (service *ShortlinkService) generateId(title string) (string) {
	return randString(6) + "-" + formatTitleForUrl(title, 43) //total: 50
}
func (service *ShortlinkService) GenerateShortlinkFromFeedItem(item gofeed.Item) (string){
	id := service.generateId(item.Title)
	log.Println("Generated id: "+ id)
	service.repo.AddShortlink(id, item.Link)
	return service.url + "/" +id
}

func formatTitleForUrl(title string, maxlength int) (url  string){
	url = strings.ToLower(title)
	if len(url) > maxlength {
		url = url[:maxlength-1]
	}
	reg := regexp.MustCompile("[^a-zA-Z ]*")
	url = reg.ReplaceAllString(url, "")
	url = strings.Replace(url, " ", "-", -1)
	return url
}

func randString(maxChars int) string {
	return token.New(maxChars).Encode()
}

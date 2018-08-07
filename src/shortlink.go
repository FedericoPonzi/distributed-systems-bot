package main

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"github.com/marksalpeter/token"
	"log"
	"regexp"
)

type ShortlinkService struct {
	repo *MysqlRepository
	url  string
}


func NewShortLinkService(repo *MysqlRepository) (*ShortlinkService) {

	return &ShortlinkService{repo: repo, url: "https://ds.fponzi.me"}
}


func (service *ShortlinkService) generateShortlink(item gofeed.Item) (string){
	t := item.Title
	id := randString(6) + "-" + formatTitleForUrl(t, 43) //total: 50
	log.Println("Generated id: "+ id)
	service.repo.addShortlink(id, item.Link)

	return service.url + "/" +id
}

func formatTitleForUrl(title string, maxlength int) (url  string){
	url = strings.ToLower(title)
	if len(url) > maxlength {
		url = url[:maxlength-1]
	}
	reg := regexp.MustCompile("[^a-zA-Z ]+")
	url = reg.ReplaceAllString(url, "")
	url = strings.Replace(url, " ", "-", -1)
	return url
}



func randString(maxChars int) string {
	return token.New(maxChars).Encode()
}

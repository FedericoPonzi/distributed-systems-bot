package main

import (
	"github.com/mmcdole/gofeed"
	"strings"
)

type ShortlinkService struct {
	repo *MysqlRepository
}


func NewShortLinkService(repo *MysqlRepository) (*ShortlinkService) {

	return &ShortlinkService{repo: repo}
}


func (service *ShortlinkService) generateShortlink(item gofeed.Item) (string){
	t := item.Title

	id := "first-" + formatTitleForUrl(t)
	service.repo.addShortlink(id, item.Link)

	return id
}

func formatTitleForUrl(title string) (url  string){
	url = strings.ToLower(title)
	if len(url) > 45 {
		url = url[:45]
	}
	url = strings.Replace(url, " ", "-", -1)
	return url
}

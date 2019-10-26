package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/marksalpeter/token"
)

type ShortlinkService struct {
	repo *MysqlRepository
	url  string
}

func NewShortLinkService(repo *MysqlRepository) *ShortlinkService {

	return &ShortlinkService{repo: repo, url: "https://ds.fponzi.me"}
}
func (service *ShortlinkService) FetchTitle(link string) (title string, err error) {
	res, err := http.Get(link)
	if err != nil {
		log.Println("Impossible to fetch the webpage.")
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("Wrong status code")
		return "", errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	r := regexp.MustCompile(`(?i)<\s*title\s*>\s*(.+)\s*<\s*/title\s*>`)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error on reading response body.")
		return "", err
	}
	re := r.FindSubmatch(body)
	if re == nil {
		log.Fatal("Impossible to find title in the web page. Please supply a custom one.")
		return "", err
	}
	return string(re[1]), err
}
func (service *ShortlinkService) GenerateShortlink(link string) (shortlink string, title string, err error) {
	title, err = service.FetchTitle(link)
	if err != nil {
		return "", "", err
	}
	id := service.generateId(title)
	err = service.repo.AddShortlink(id, link)
	if err != nil {
		return "", "", err
	}
	return service.url + "/" + id, title, nil
}

func (service *ShortlinkService) generateId(title string) string {
	return randString(6) + "-" + formatTitleForUrl(title, 43) //total: 50
}
func (service *ShortlinkService) GenerateShortlinkWithTitle(link string, title string) string {
	id := service.generateId(title)
	log.Println("Generated id: " + id)
	service.repo.AddShortlink(id, link)
	return service.url + "/" + id
}

func formatTitleForUrl(title string, maxlength int) (url string) {
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

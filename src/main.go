package main

import (
	"fmt"
	"log"
)
import (
	_ "github.com/go-sql-driver/mysql"
)
var config *Config

func main() {
	fmt.Println("Welcome to DistributedSystems bot!")
	var err  error
	config, err = loadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config", err)
	}
	repo := NewMysqlRepository()
	defer repo.close()
	defer func() {	log.Println("End of the execution. Thanks for playing :)")
	}()
	twitterHandler := NewTwitterHandler(config.Twitter)

	feedHandler := NewFeedHandler(repo, twitterHandler)
	feedHandler.main()
	//fmt.Println("Going to publish:", feedHandler.getUpdatedItems(feedHandler.fetchAllRss()))
	//
	//go twitterHandler.runStreaming();
	//twitterHandler.runStreaming()

	//telegramBot := NewTelegramHandler(config.Telegram, twitterHandler)
	//telegramBot.run()
	//https://stackoverflow.com/questions/38386762/running-code-at-noon-in-golang

	/*url := "https://blog.acolyer.org/feed/"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	fmt.Println(feed.Updated)
	*/


}




package main

import (
	"fmt"
	"log"
)
import (
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"github.com/FedericoPonzi/distributed-systems-bot/src/config"
	"github.com/FedericoPonzi/distributed-systems-bot/src/repository"
	"github.com/FedericoPonzi/distributed-systems-bot/src/twitter_handler"
	"github.com/FedericoPonzi/distributed-systems-bot/src/feed_rss"
	"github.com/FedericoPonzi/distributed-systems-bot/src/shortlink"
)



var configPath string
var fetchRssRun bool
var shortLink string
func parseArgs() {
	flag.StringVar(&configPath, "config", "config.yaml", "Complete path to the config yaml file.")
	flag.BoolVar(&fetchRssRun, "fetch-rss", false, "Fetch feed rss")
	flag.StringVar(&shortLink, "shortlink", "", "Generate a shortlink.")
	flag.Parse()
}

func main() {
	fmt.Println("Welcome to DistributedSystems bot!")
	parseArgs()
	fmt.Println(configPath)

	config, err := config.LoadConfig(configPath)

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	repo := repository.NewMysqlRepository(config)

	defer func() {
		repo.Close()
		log.Println("End of the execution. Thanks for playing :)")

	}()

	twitterHandler := twitter_handler.NewTwitterHandler(config.Twitter)

	if fetchRssRun {
		feedHandler := feed_rss.NewFeedHandler(repo, twitterHandler)
		feedHandler.Main()
	} else if len(shortLink) > 0{
		fmt.Println("I'm going to generate a shortlink for: " + shortLink + " just a sec...")
		shortLinkService := shortlink.NewShortLinkService(repo)
		generated := shortLinkService.GenerateShortlinkFromLink(shortLink)
		fmt.Println("Generated shortlink: " + generated)
	}
	//fmt.Println("Going to publish:", feedHandler.getUpdatedItems(feedHandler.fetchAllRss()))
	//
	//go twitterHandler.runStreaming();
	//twitterHandler.runStreaming()/home/bots/distributed-systems-bot/src/distributed-systems-bot --fetch-rss --config /home/bots/distributed-systems-bot/src/config.yaml 2>&1 > /home/bots/distributed-systems-bot/output

	//telegramBot := NewTelegramHandler(config.Telegram, twitterHandler)
	//telegramBot.run()
	//https://stackoverflow.com/questions/38386762/running-code-at-noon-in-golang

}

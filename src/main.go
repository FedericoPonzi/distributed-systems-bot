package main

import (
	"fmt"
	"log"
)
import (
	_ "github.com/go-sql-driver/mysql"
	"flag"
)



var configPath string
var fetchRssRun bool
func parseArgs() {
	flag.StringVar(&configPath, "config", "config.yaml", "Complete path to the config yaml file.")
	flag.BoolVar(&fetchRssRun, "fetch-rss", false, "Fetch feed rss")
	flag.Parse()
}

func main() {
	fmt.Println("Welcome to DistributedSystems bot!")
	parseArgs()
	fmt.Println(configPath)

	config, err := loadConfig(configPath)

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	repo := NewMysqlRepository(config)

	defer func() {
		repo.close()
		log.Println("End of the execution. Thanks for playing :)")

	}()

	twitterHandler := NewTwitterHandler(config.Twitter)

	if fetchRssRun {
		feedHandler := NewFeedHandler(repo, twitterHandler)
		feedHandler.main()
	}
	//fmt.Println("Going to publish:", feedHandler.getUpdatedItems(feedHandler.fetchAllRss()))
	//
	//go twitterHandler.runStreaming();
	//twitterHandler.runStreaming()/home/bots/distributed-systems-bot/src/distributed-systems-bot --fetch-rss --config /home/bots/distributed-systems-bot/src/config.yaml 2>&1 > /home/bots/distributed-systems-bot/output

	//telegramBot := NewTelegramHandler(config.Telegram, twitterHandler)
	//telegramBot.run()
	//https://stackoverflow.com/questions/38386762/running-code-at-noon-in-golang

}

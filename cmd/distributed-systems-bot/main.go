package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/FedericoPonzi/distributed-systems-bot/pkg/handlers"
	app "github.com/FedericoPonzi/distributed-systems-bot/pkg/main"

	// importing for side effects
	_ "github.com/go-sql-driver/mysql"
)

var configPath string
var fetchRssRun bool
var shortLink string
var telegramBot bool
var dryRun bool

func parseArgs() {
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run: runs every function, but dosen't publish the tweet.")
	flag.StringVar(&configPath, "config", "config.yaml", "Complete path to the config yaml file.")
	flag.BoolVar(&fetchRssRun, "fetch-rss", false, "Fetch feed rss")
	flag.StringVar(&shortLink, "shortlink", "", "Generate a shortlink.")
	flag.BoolVar(&telegramBot, "telegram-bot", false, "Run the telegram bot")
	flag.Parse()
}

func main() {
	fmt.Println("Welcome to DistributedSystems bot!")
	parseArgs()
	fmt.Println(configPath)

	config, err := app.LoadConfig(configPath)

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	config.Twitter.DryRun = dryRun

	repo := app.NewMysqlRepository(config)

	defer func() {
		repo.Close()
		log.Println("End of the execution. Thanks for playing :)")
	}()

	twitterHandler := handlers.NewTwitterHandler(repo, config.Twitter)
	fmt.Println("Fetch-rss:", strconv.FormatBool(fetchRssRun), " , telegram: ", strconv.FormatBool(telegramBot))
	if fetchRssRun {
		fmt.Println("Running fetch rss command.")
		feedHandler := handlers.NewFeedHandler(repo, twitterHandler)
		feedHandler.Main()
		fmt.Println("Done fetching feeds.")
		return
	}

	if len(shortLink) > 0 {
		fmt.Println("I'm going to generate a shortlink for: " + shortLink + " just a sec...")
		shortLinkService := app.NewShortLinkService(repo)
		shortlink, title, _ := shortLinkService.GenerateShortlink(shortLink)
		fmt.Println("Generated shortlink: " + shortlink + ", parsed title:" + title)
		return
	}

	if telegramBot {
		fmt.Println("Running telegram bot.")
		handler := handlers.NewTelegramHandler(config.Telegram, twitterHandler)
		handler.Run()
		return
	}

	//fmt.Println("Going to publish:", feedHandler.getUpdatedItems(feedHandler.fetchAllRss()))
	//
	//go twitterHandler.runStreaming();
	//twitterHandler.runStreaming()/home/bots/distributed-systems-bot/src/distributed-systems-bot --fetch-rss --config /home/bots/distributed-systems-bot/src/config.yaml 2>&1 > /home/bots/distributed-systems-bot/output

	//telegramBot := NewTelegramHandler(config.Telegram, twitterHandler)
	//telegramBot.run()
	//https://stackoverflow.com/questions/38386762/running-code-at-noon-in-golang

}

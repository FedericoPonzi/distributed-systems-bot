package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramHandler struct {
	botApi *tgbotapi.BotAPI
	userId int
	twitterHandler *TwitterHandler
}


func (handler *TelegramHandler) postUpdate(text string) string {
	if !handler.containsLink(text) {
		return "Sorry at this time I can only share links for you!"
	}
	err := handler.twitterHandler.PublishLink(text)
	/*tweet, resp, err := handler.twitterHandler.bot.Statuses.Update(text, nil)
	if err != nil {
		errStr := fmt.Sprint("Error twitting message:", tweet, resp)
		log.Println(errStr)
		return errStr
	}
	*/

	if err != nil {
		return err.Error()
	}
	return "Update sent! Good Job!"

}


func NewTelegramHandler(config TelegramConfig, twitterHandler *TwitterHandler) (toRet TelegramHandler) {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey);
	if err != nil {
		log.Panic(err)
	}
	toRet.botApi = bot
	toRet.userId = config.UserId
	toRet.twitterHandler = twitterHandler
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return toRet
}



func (handler *TelegramHandler) run(){
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := handler.botApi.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("TelegramHandler: impossible to get updates.")
	}
	for update := range updates {
		if update.Message == nil ||  handler.defaultMessage(update) {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		resp := tgbotapi.NewMessage(update.Message.Chat.ID, "Ok")
		resp.ReplyToMessageID = update.Message.MessageID
		respText := ""
		if update.Message.Text == "/start" {
			respText = "Hello man! What's up? I will help you send updates to Twitter. I can only handle links updates tho.\nLet's get started, send me your first link!"
		} else {
			respText = handler.postUpdate(update.Message.Text)
		}
		resp.Text = respText
		handler.botApi.Send(resp)
	}
}

func (handler *TelegramHandler) defaultMessage(update tgbotapi.Update) (bool) {
	if update.Message.From.ID != handler.userId {
		resp := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello there! " +
			"\nSorry, but at the moment I can't do anything useful for you." +
			"\nI'll go open source in the future, but I'm just not ready yet :) " +
			"\nIn the meanwhile, if you're interested in Distributed Systems, " +
			"feel free to follow me on https://twitter.com/DistribSystems")
		resp.ReplyToMessageID = update.Message.MessageID
		handler.botApi.Send(resp)
		return true
	}
	return false
}

func (handler *TelegramHandler) containsLink(s string) bool {
	return true
}


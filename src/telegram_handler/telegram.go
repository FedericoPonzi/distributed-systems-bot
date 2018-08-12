package telegram_handler

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramHandler struct {
	botApi *tgbotapi.BotAPI
	userId int
	twitterHandler TwitterHandler
}


func (handler *TelegramHandler) postUpdate(text string) string {

	/*tweet, resp, err := handler.twitterHandler.bot.Statuses.Update(text, nil)
	if err != nil {
		errStr := fmt.Sprint("Error twitting message:", tweet, resp)
		log.Println(errStr)
		return errStr
	}
	*/
	return "Update sent! Good Job!"

}


func NewTelegramHandler(config TelegramConfig, twitterHandler TwitterHandler) (toRet TelegramHandler) {
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
		if update.Message == nil || update.Message.From.ID == handler.userId {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = handler.postUpdate(update.Message.Text)

		handler.botApi.Send(msg)
	}
}


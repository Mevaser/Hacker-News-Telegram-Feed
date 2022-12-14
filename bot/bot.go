package bot

import (
	"hn_feed/bot/handlers"
	"hn_feed/bot/utils"
	"hn_feed/config"
	"log"

	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func Init() *telebot.Bot {
	log.Printf("[*] Creating the bot.")
	bot := utils.CreateBot(config.Get().BotToken)

	log.Printf("[*] Creating handlers for the bot.")
	initHandlers(bot)
	return bot
}

func Run(bot *telebot.Bot) {
	log.Printf("[*] Starting the bot.")
	bot.Start()
}

func initHandlers(bot *telebot.Bot) {
	bot.Handle("/start", handlers.HandleStart)
	bot.Handle(telebot.OnChannelPost, handlers.ChannelCommandsHandler(
		map[string]telebot.HandlerFunc{
			"/help":      handlers.OnChannelHelp,
			"/info":      handlers.OnChannelInfo,
			"/feed":      handlers.OnChannelConfigureFeedType,
			"/count":     handlers.OnChannelConfigureCount,
			"/score":     handlers.OnChannelConfigureScore,
			"/whitelist": handlers.OnChannelConfigureWhitelist,
			"/blacklist": handlers.OnChannelConfigureBlacklist,
			"/register":  handlers.OnChannelRegister,
		},
	))

	admin := bot.Group()
	admin.Use(middleware.Whitelist(config.Get().AdminIds...))
	admin.Handle(&handlers.AdminBtnList, handlers.HandleAdminListChannels)
}

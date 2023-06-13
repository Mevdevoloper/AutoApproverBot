//   Approver Bot
//   Copyright (C) 2021 Reeshuxd (@reeshuxd)

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	bot, err := gotgbot.NewBot(
		os.Getenv("TOKEN"),
		&gotgbot.BotOpts{
			APIURL:      "",
			Client:      http.Client{},
			GetTimeout:  gotgbot.DefaultGetTimeout,
			PostTimeout: gotgbot.DefaultPostTimeout,
		},
	)
	if err != nil {
		fmt.Println("Failed to create bot:", err.Error())
	}

	updater := ext.NewUpdater(
		&ext.UpdaterOpts{
			ErrorLog: nil,
			DispatcherOpts: ext.DispatcherOpts{
				Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
					fmt.Println("an error occurred while handling update:", err.Error())
					return ext.DispatcherActionNoop
				},
				Panic:       nil,
				ErrorLog:    nil,
				MaxRoutines: 0,
			},
		})
	dp := updater.Dispatcher

	// Commands
	dp.AddHandler(handlers.NewCommand("start", Start))
	dp.AddHandler(handlers.NewChatJoinRequest(nil, Approve))
	dp.AddHandler(handlers.NewCommand("status", Status))

	// Start Polling()
	poll := updater.StartPolling(bot, &ext.PollingOpts{DropPendingUpdates: true})
	if poll != nil {
		fmt.Println("Failed to start bot:", poll.Error())
	}

	fmt.Printf("@%s has been successfully started\n💝Made by @CodeMasterTG\n", bot.Username)
	updater.Idle()
}

func Start(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "private" {
		return nil
	}

	user := ctx.EffectiveSender.User
	text := `
<b>Hello <a href="tg://user?id=%v">%v</a></b> ❤️
I am a bot made for accepting newly coming join requests at the time they come.

Bot made with 💝 by <a href="t.me/CodeMasterTG">Code Master Bots</a> for you!
<b>Support Chat:</b> <a href="t.me/+4KDIm0IQ_NQ0NDdl">Support Chat</a>
	`
	ctx.EffectiveMessage.Reply(
		bot,
		fmt.Sprintf(text, user.Id, user.FirstName),
		&gotgbot.SendMessageOpts{
			ParseMode:             "html",
			DisableWebPagePreview: true,
		},
	)
	return nil
}

func Approve(bot *gotgbot.Bot, ctx *ext.Context) error {
	_, err := bot.ApproveChatJoinRequest(ctx.EffectiveChat.Id, ctx.EffectiveSender.User.Id)
	if err != nil {
		fmt.Println("Error while approving:", err.Error())
	}
	return nil
}

func Status(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Get the number of users in the chat
	count, err := bot.GetChatMembersCount(ctx.EffectiveChat.Id)
	if err != nil {
		fmt.Println("Error while getting chat members count:", err.Error())
		return nil
	}

	response := fmt.Sprintf("Number of users in the chat: %d",

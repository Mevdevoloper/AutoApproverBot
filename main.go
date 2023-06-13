//   Approver Bot
//   Copyright (C) 2021 Reeshuxd (@reeshuxd)

package main

import (
	"fmt"
	"net/http"

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
	dp.AddHandler(handlers.NewCommand("broadcast", Broadcast))
	dp.AddHandler(handlers.NewCommand("force_subscribe", ForceSubscribe))

	// Start Polling()
	poll := updater.StartPolling(bot, &ext.PollingOpts{DropPendingUpdates: true})
	if poll != nil {
		fmt.Println("Failed to start bot:", poll.Error())
	}

	fmt.Printf("@%s has been successfully started\n💝Made by @CodeMasterTG\n", bot.Username)
	updater.Idle()
}

func ForceSubscribe(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Check if the user is already a member of the chat
	isMember, err := bot.IsMemberOfChat(ctx.EffectiveChat.Id, ctx.EffectiveSender.User.Id)
	if err != nil {
		fmt.Println("Error while checking membership:", err.Error())
		return nil
	}

	if !isMember {
		// User is not a member, send a message asking to join
		ctx.EffectiveMessage.Reply(
			bot,
			"Please join the chat to use this bot's features.",
			nil,
		)
	} else {
		// User is already a member, no action required
		ctx.EffectiveMessage.Reply(
			bot,
			"You are already a member of the chat.",
			nil,
		)
	}

	return nil
}

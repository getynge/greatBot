package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/events"
	"os"
	"os/signal"
	"syscall"
)

func getBotString() string {
	key := os.Getenv("DISCORD_BOT_KEY")
	return fmt.Sprintf("Bot %s", key)
}

func main() {
	// if this ever hits prod, move key to environment variable
	session, err := discordgo.New(getBotString())

	if err != nil {
		fmt.Printf("Could not initialize bot with error: %+v\n", err)
		return
	}

	user, err := session.User("@me")

	if err != nil {
		fmt.Printf("Could not get bot user with error: %+v\n", err)
		return
	}

	dispatcher := events.NewDispatcher(user.ID, "$")
	session.AddHandler(dispatcher.CommandHandler)
	session.AddHandler(dispatcher.ReadyHandler)

	err = session.Open()
	if err != nil {
		fmt.Printf("failed to open connection to discord %+v", err)
		return
	}

	signalCatcher := make(chan os.Signal, 1)
	signal.Notify(signalCatcher, syscall.SIGINT, syscall.SIGKILL, os.Interrupt, os.Kill)
	<-signalCatcher

	session.Close()
}

package main

import (
	"context"
	"log"

	bot "github.com/serjyuriev/dong-size-bot/internal/app/dong_size_bot"
)

func main() {
	b, err := bot.NewDongSizeBot(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}

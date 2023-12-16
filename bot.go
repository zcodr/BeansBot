package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// Uncomment below line if you are going to use uptimerobot to ping
	//"net/http"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Uncomment this code block if you are going to use uptimerobot to ping
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, World!")
	// })

	go http.ListenAndServe(":8080", nil)

	// Load bot token
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bottoken := os.Getenv("TOKEN")

	// Create a new Discord session using the bot token from .env
	bot, err := discordgo.New("Bot " + bottoken)
	if err != nil {
		panic(err)
	}

	// Register events
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	// Start sesson
	err = bot.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "BEANS")
	fmt.Println("logged in as user " + string(s.State.User.ID))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	tokens := strings.Split(m.Content, " ")

	if tokens[0] == "/beans" {
		filename := "Words/common.txt"
		rareness := rand.Intn(10000) + 1
		if rareness <= 500 {
			filename = "Words/uncommon.txt"
		}
		if rareness <= 100 {
			filename = "Words/rare.txt"
		}
		if rareness <= 25 {
			filename = "Words/ultimate.txt"
		}

		file, _ := os.Open(filename)
		filescanner := bufio.NewScanner(file)
		filescanner.Split(bufio.ScanLines)
		var filelines []string
		for filescanner.Scan() {
			filelines = append(filelines, filescanner.Text())
		}
		s.ChannelMessageSend(m.ChannelID, filelines[rand.Intn(399)+1]+" beans")
		file.Close()
	}
}

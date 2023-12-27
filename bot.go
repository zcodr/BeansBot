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
		congrats_message := ""
		filename := "Words/common.txt"
		rareness := rand.Intn(10000) + 1
		fmt.Println(rareness)
		if rareness <= 100 {
			congrats_message = "<@" + m.Author.ID + "> :bell: :partying_face: YOU GOT AN **UNCOMMON** BEAN :partying_face: :bell:   "
			filename = "Words/uncommon.txt"
		}
		if rareness <= 10 {
			congrats_message = "<@" + m.Author.ID + "> :bell: :bell: :partying_face: :partying_face: :tada: :tada: YOU GOT A **RARE** BEAN :tada: :tada: :partying_face: :partying_face: :bell: :bell:   "
			filename = "Words/rare.txt"
		}
		if rareness <= 1 {
			congrats_message = "<@" + m.Author.ID + "> :bell: :bell: :partying_face: :partying_face: :tada: :tada: :regional_indicator_y: :regional_indicator_o: :regional_indicator_u:    :regional_indicator_g: :regional_indicator_o: :regional_indicator_t:    :regional_indicator_a: :regional_indicator_n:    :regional_indicator_u: :regional_indicator_l: :regional_indicator_t: :regional_indicator_i: :regional_indicator_m: :regional_indicator_a: :regional_indicator_t: :regional_indicator_e:    :regional_indicator_b: :regional_indicator_e: :regional_indicator_a: :regional_indicator_n: :tada: :tada: :partying_face: :partying_face: :bell: :bell:   "
			filename = "Words/ultimate.txt"
		}

		file, _ := os.Open(filename)
		filescanner := bufio.NewScanner(file)
		filescanner.Split(bufio.ScanLines)
		var filelines []string
		for filescanner.Scan() {
			filelines = append(filelines, filescanner.Text())
		}

		s.ChannelMessageSendReply(m.ChannelID, congrats_message+filelines[rand.Intn(len(filelines))]+" beans\n", m.Reference())
		file.Close()
	}
}

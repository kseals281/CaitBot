package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const ErrOpenJournalPrompt = JournalErr("unable to open journal prompts")

type JournalErr string

func (e JournalErr) Error() string {
	return string(e)
}

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()

// Read in all options from environment variables and command line arguments.
func init() {
	rand.Seed(time.Now().Unix())

	// Discord Authentication Token
	Session.Token = os.Getenv("CAITBOT")
	fmt.Println()
	if Session.Token == "" {
		// Pointer, flag, default, description
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {

	// Declare any variables needed later.
	var err error

	// Setup interrupt
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Verify the Token is valid and grab user information
	Session.State.User, err = Session.User("@me")
	errCheck("error retrieving account", err)

	Session.AddHandler(CommandHandler)
	Session.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A friendly Carrot Cait bot!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("CaitBot has started on %d servers\n", len(servers))
	})

	// Open a websocket connection to Discord
	err = Session.Open()
	defer Session.Close()
	errCheck("Error opening connection to Discord", err)

	<-interrupt
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	botID := s.State.User.ID
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	if m.Content == "!journal" {
		promptJournal(s, m.ChannelID)
	}
}

func promptJournal(s *discordgo.Session, chID string) error {
	f, err := os.Open("journal_prompts.txt")
	if err != nil {
		s.ChannelMessageSend(chID, "Unable to open journal prompts.")
		return ErrOpenJournalPrompt
	}
	s.ChannelMessageSend(chID, prompt(f))
	return nil
}

func prompt(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	var prompts []string
	for scanner.Scan() {
		prompts = append(prompts, scanner.Text())
	}
	if len(prompts) == 0 {
		return ""
	}
	i := rand.Intn(len(prompts))
	return prompts[i]
}

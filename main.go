package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strings"
	"syscall"

	"github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/javascript"
	"github.com/aichaos/rivescript-go/sessions/memory"

	"github.com/bwmarrin/discordgo"
)

var (
	bot *rivescript.RiveScript = rivescript.New(&rivescript.Config{
		UTF8:           true,
		Debug:          false,
		Strict:         true,
		Depth:          50,
		CaseSensitive:  false,
		SessionManager: memory.New(),
		Seed:           int64(rand.Uint64()),
	})
	dg       *discordgo.Session
	guilds   []string
	channels []string
	admins   []string
)

func init() {
	bot.SetUnicodePunctuation(`[^a-zA-Z0-9 ÄÖÜäöüßÁČĎÉĚÍŇÓŘŠŤÚŮÝŽáčďéěíňóřšťúůýž]`)
	bot.SetHandler("javascript", javascript.New(bot))

	err := bot.LoadDirectory("brain", ".rive")
	if err != nil {
		log.Fatal(err)
	}
	err = bot.SortReplies()
	if err != nil {
		log.Fatal(err)
	}

	token, err := bot.GetGlobal("token")
	if err != nil {
		log.Fatal(err)
	}

	g, err := bot.GetGlobal("guilds")
	if err != nil {
		log.Fatal(err)
	}
	guilds = strings.Split(g, ",")

	c, err := bot.GetGlobal("channels")
	if err != nil {
		log.Fatal(err)
	}
	channels = strings.Split(c, ",")

	a, err := bot.GetGlobal("admins")
	if err != nil {
		log.Fatal(err)
	}
	admins = strings.Split(a, ",")

	bot.SetGlobal("token", "undefined")
	bot.SetGlobal("admins", "undefined")
	bot.SetGlobal("channels", "undefined")
	bot.SetGlobal("guilds", "undefined")

	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageUpdate) {
		if m.Author.ID == s.State.User.ID || !slices.Contains(guilds, m.GuildID) || !slices.Contains(channels, m.ChannelID) {
			return
		}
		var msg string = CleanMessage(m.Content)
		if reply, err := bot.Reply(m.Author.Username, msg); err != nil {
			log.Println(err, msg)
		} else {
			if len([]rune(reply)) > 0 {
				_, err = s.ChannelMessageSendReply(m.ChannelID, reply, &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
					GuildID:   m.GuildID,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID || !slices.Contains(guilds, m.GuildID) || !slices.Contains(channels, m.ChannelID) {
			return
		}

		if m.Content == "!reload" && slices.Contains(admins, m.Author.ID) {
			bot = rivescript.New(&rivescript.Config{
				UTF8:           true,
				Debug:          false,
				Strict:         true,
				Depth:          50,
				CaseSensitive:  false,
				SessionManager: memory.New(),
				Seed:           int64(rand.Uint64()),
			})

			bot.SetUnicodePunctuation(`[^a-zA-Z0-9 äöüß]`)
			bot.SetHandler("javascript", javascript.New(bot))

			err := bot.LoadDirectory("brain", ".rive")
			if err != nil {
				log.Fatal(err)
			}
			err = bot.SortReplies()
			if err != nil {
				log.Fatal(err)
			}

			g, err := bot.GetGlobal("guilds")
			if err != nil {
				log.Fatal(err)
			}
			guilds = strings.Split(g, ",")

			c, err := bot.GetGlobal("channels")
			if err != nil {
				log.Fatal(err)
			}
			channels = strings.Split(c, ",")

			a, err := bot.GetGlobal("admins")
			if err != nil {
				log.Fatal(err)
			}
			admins = strings.Split(a, ",")

			bot.SetGlobal("token", "undefined")
			bot.SetGlobal("admins", "undefined")
			bot.SetGlobal("channels", "undefined")
			bot.SetGlobal("guilds", "undefined")

			if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
				log.Println(err)
			}
			return
		}
		var msg string = CleanMessage(m.Content)
		if reply, err := bot.Reply(m.Author.Username, msg); err != nil {
			log.Println(err, msg)
		} else {
			if len([]rune(reply)) > 0 {
				_, err = s.ChannelMessageSendReply(m.ChannelID, reply, &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
					GuildID:   m.GuildID,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	})

	dg.Identify.Intents = discordgo.IntentsAll
	err := dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

var spaces *regexp.Regexp = regexp.MustCompile(`\s{1,}`)

func CleanMessage(in string) string {
	in = spaces.ReplaceAllString(in, " ")
	return strings.TrimSpace(in)
}

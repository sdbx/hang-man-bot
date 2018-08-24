package mole

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sdbx/hang-man-bot/config"
	"github.com/sdbx/hang-man-bot/display"
	"github.com/sdbx/hang-man-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/game"
)

var ErrAlreadyExist = errors.New("sdbx/hang-man-bot/mole you already registed")

type GameMeta struct {
	Creator  string
	Hint     string
	Korean   bool
	Solution []rune
}

type moleGame struct {
	hint  bool
	manID int
}

type Mole struct {
	mu sync.RWMutex

	sess  *discordgo.Session
	msgID string
	msgs  []string

	requests map[string]GameMeta

	cGame *game.Game
	cMole *moleGame
}

func New(sess *discordgo.Session) *Mole {
	return &Mole{
		sess:     sess,
		requests: make(map[string]GameMeta),
		cGame:    nil,
		cMole:    nil,
	}
}

func (m *Mole) Add(gm GameMeta) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.requests[gm.Creator]; ok {
		return ErrAlreadyExist
	}
	m.requests[gm.Creator] = gameMeta
	return nil
}

func (m *Mole) Start() {
	msg, err := m.sess.ChannelMessageSendComplex(config.Conf.ChannelID,
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Description: "​",
			},
		})
	if err != nil {
		panic(err)
	}

	m.msgID = msg.ID

	for {
		var g *game.Game
		if r, ok := m.PickRequest(); ok {
			g = m.createGame(r)
		} else {
			g = m.createGame(m.createMetaByPreset())
		}
		m.mu.Lock()
		m.cGame = g
		m.mu.Unlock()

		m.Display()
	}
}

func (m *Mole) Game() *game.Game {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cGame
}

func (m *Mole) MessageHandler(s *discordgo.Session, me *discordgo.MessageCreate) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.msgs = append(m.msgs, me.ID)

	if len(m.msgs) > 5 {
		for i := 0; i < len(m.msgs)-5; i++ {
			id := m.msgs[i]
			go func() {
				err := m.sess.ChannelMessageDelete(config.Conf.ChannelID, id)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
		m.msgs = m.msgs[len(m.msgs)-5:]
	}
}

func (m *Mole) deleteMessages() error {
	messages, _ := m.sess.ChannelMessages(config.Conf.ChannelID, 100, "", "", "")
	messageIDs := make([]string, 0, len(messages))
	for _, msg := range messages {
		if m.msgID == msg.ID {
			continue
		}
		messageIDs = append(messageIDs, msg.ID)
	}
	return m.sess.ChannelMessagesBulkDelete(config.Conf.ChannelID, messageIDs)
}

func (m *Mole) Display() {
	err := m.deleteMessages()
	if err != nil {
		fmt.Println(err)
	}

	stop, c := m.cGame.Listen()
	defer stop()

	m.cMole = &moleGame{
		manID: display.PickID(config.Conf.MaxHp),
	}

L:
	for {
		select {
		case r := <-c:
			switch e := r.(type) {
			case *game.EndEvent:
				if e.Win {

				} else {

				}
				break L
			case *game.TurnEvent:
				m.RefreshMsg()
			case *game.RevealEvent:
				fmt.Println(e)
			case *game.HintEvent:
				fmt.Println(e)
			}
		}
	}
}

func (m *Mole) RefreshMsg() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	str := display.DisplayDefault(
		config.Conf.MaxHp, m.cGame.Hp(), m.cGame.Log(), m.cGame.Current())

	me := &discordgo.MessageEdit{
		ID:      m.msgID,
		Channel: config.Conf.ChannelID,
		Embed: &discordgo.MessageEmbed{
			Description: str,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: utils.CreateURL(config.Conf.URL, config.Conf.MaxHp, m.cMole.manID, m.cGame.Hp()),
			},
		},
	}

	m.sess.ChannelMessageEditComplex(me)
}

func (m *Mole) PickRequest() (GameMeta, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	reqs := m.requestsToSlice()
	if len(reqs) == 0 {
		return GameMeta{}, false
	}

	req := reqs[rand.Intn(len(reqs))]
	delete(m.requests, req.Creator)

	return req, true
}

func (m *Mole) createMetaByPreset() GameMeta {
	p := getPreset()
	return GameMeta{
		Creator:  "Bot",
		Korean:   p.kor,
		Solution: []rune(p.str),
		Hint:     p.hint,
	}
}

func (m *Mole) requestsToSlice() []GameMeta {
	out := make([]GameMeta, 0, len(m.requests))
	for _, r := range m.requests {
		out = append(out, r)
	}
	return out
}

func (m *Mole) createGame(gm GameMeta) *game.Game {
	return game.New(
		gm.Creator,
		gm.Hint,
		gm.Korean,
		time.Duration(config.Conf.Cool)*time.Second,
		gm.Solution,
		config.Conf.MaxHp,
	)
}

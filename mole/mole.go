package mole

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sdbx/hang-man-bot/config"
	"github.com/sdbx/hang-man-bot/display"
	"github.com/sdbx/hang-man-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/sdbx/hang-man-bot/game"
)

const (
	MaxMessages = 5
	GameDelay   = 3
)

var ErrAlreadyExist = errors.New("sdbx/hang-man-bot/mole you already registed")

type moleGame struct {
	meta    game.GameMeta
	hint    bool
	start   time.Time
	recents []string
	msgs    []string
	manID   int
}

type Mole struct {
	mu           sync.RWMutex
	msgID        string
	channelID    string
	logChannelID string

	sess     *discordgo.Session
	requests map[string]game.GameMeta

	cGame *game.Game
	cMole *moleGame
}

func New(sess *discordgo.Session, channelID string, logChannelID string) *Mole {
	return &Mole{
		sess:         sess,
		requests:     make(map[string]game.GameMeta),
		channelID:    channelID,
		logChannelID: logChannelID,
		cGame:        nil,
		cMole:        nil,
	}
}

func (m *Mole) Start() {
	msg, err := m.sess.ChannelMessageSendComplex(m.channelID,
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
		var g game.GameMeta
		if r, ok := m.PickRequest(); ok {
			g = r
		} else {
			g = m.createMetaByPreset()
		}

		m.mu.Lock()
		m.cMole = &moleGame{
			meta:    g,
			hint:    false,
			start:   time.Now(),
			msgs:    []string{},
			recents: []string{},
			manID:   display.PickManID(config.Conf.MaxHp),
		}
		m.cGame = m.createGame(g)
		m.mu.Unlock()

		m.HandleGame()
	}
}

func (m *Mole) Game() *game.Game {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cGame
}

func (m *Mole) Add(gm game.GameMeta) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.requests[gm.Creator]; ok {
		return ErrAlreadyExist
	}
	m.requests[gm.Creator] = gm
	return nil
}

func (m *Mole) PickRequest() (game.GameMeta, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	reqs := m.requestsToSlice()
	if len(reqs) == 0 {
		return game.GameMeta{}, false
	}

	req := reqs[rand.Intn(len(reqs))]
	delete(m.requests, req.Creator)

	return req, true
}

func (m *Mole) createMetaByPreset() game.GameMeta {
	p := getPreset()
	return game.GameMeta{
		Creator:  "Bot",
		Korean:   p.kor,
		Solution: []rune(p.str),
		Hint:     p.hint,
	}
}

func (m *Mole) requestsToSlice() []game.GameMeta {
	out := make([]game.GameMeta, 0, len(m.requests))
	for _, r := range m.requests {
		out = append(out, r)
	}
	return out
}

func (m *Mole) createGame(gm game.GameMeta) *game.Game {
	return game.New(
		gm.Creator,
		gm.Hint,
		gm.Korean,
		utils.IntToSeconds(config.Conf.Cool),
		gm.Solution,
		config.Conf.MaxHp,
	)
}

func (m *Mole) MessageHandler(s *discordgo.Session, me *discordgo.MessageCreate) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cMole == nil {
		return
	}

	m.logMessages(me)
	m.recentMessages(me)
}

func (m *Mole) getUserName(userID string) (string, error) {
	if userID == "Bot" {
		return "Bot", nil
	}

	ch, err := m.sess.Channel(m.channelID)
	if err != nil {
		return "", err
	}

	user, err := m.sess.GuildMember(ch.GuildID, userID)
	if err != nil {
		return "", err
	}

	username := user.User.Username
	if user.Nick != "" {
		username = user.Nick
	}

	return username, nil
}

func (m *Mole) logMessages(me *discordgo.MessageCreate) {
	username, err := m.getUserName(me.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}
	m.logMessage(username + ":" + me.Content)
}

func (m *Mole) logMessage(str string) {
	m.cMole.msgs = append(m.cMole.msgs, str)
}

func (m *Mole) recentMessages(me *discordgo.MessageCreate) {
	if me.ID == m.msgID {
		return
	}

	m.cMole.recents = append(m.cMole.recents, me.ID)
	if len(m.cMole.recents) > MaxMessages {
		for i := 0; i < len(m.cMole.recents)-MaxMessages; i++ {
			id := m.cMole.recents[i]
			go func() {
				err := m.sess.ChannelMessageDelete(m.channelID, id)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		m.cMole.recents = m.cMole.recents[len(m.cMole.recents)-MaxMessages:]
	}
}

func (m *Mole) HandleGame() {
	m.sess.ChannelMessageSend(m.channelID, fmt.Sprintf("%d초후 새로운 게임이 시작됩니다.", GameDelay))

	time.Sleep(GameDelay * time.Second)

	err := m.deleteMessages()
	if err != nil {
		log.Println(err)
	}

	m.mu.RLock()

	stop, c := m.cGame.Listen()
	defer stop()

	meta := m.cMole.meta
	username, err := m.getUserName(meta.Creator)
	if err != nil {
		log.Println(err)
	}
	m.logMessage(fmt.Sprintf("게임시작! 정답: %s 출제자: %s 힌트: %s", string(meta.Solution), username, meta.Hint))

	err = m.viewStateMsg(config.Conf.MaxHp, m.cGame.Current(), []rune{})
	if err != nil {
		log.Println(err)
	}

	m.mu.RUnlock()

	for {
		select {
		case e := <-c:
			end := m.handleEvent(e)
			if end {
				return
			}
		}
	}
}

func (m *Mole) deleteMessages() error {
	msgs, _ := m.sess.ChannelMessages(m.channelID, 100, "", "", "")

	msgids := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		if m.msgID == msg.ID {
			continue
		}
		msgids = append(msgids, msg.ID)
	}

	return m.sess.ChannelMessagesBulkDelete(m.channelID, msgids)
}

func (m *Mole) handleEvent(r game.Event) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch e := r.(type) {
	case *game.EndEvent:
		if e.Win {
			m.logMessage("승리!")
			err := m.viewWonMsg()
			if err != nil {
				log.Println(err)
			}
		} else {
			m.logMessage("패배")
			err := m.viewDiedMsg()
			if err != nil {
				log.Println(err)
			}
		}
		d := time.Now().Sub(m.cMole.start)
		m.logMessage(fmt.Sprintf("총 %f분 경과", d.Minutes()))
		m.sendLogMessages()
		return true
	case *game.TurnEvent:
		err := m.viewStateMsg(e.Hp, e.Current, e.Log)
		if err != nil {
			log.Println(err)
		}

		username, err := m.getUserName(e.Player)
		if err != nil {
			log.Println(err)
			break
		}

		if e.Right {
			m.logMessage(fmt.Sprintf("%s님이 맞추셨습니다. (%c)", username, e.Char))
		} else {
			m.logMessage(fmt.Sprintf("%s님이 틀리셧습니다. (%c)", username, e.Char))
		}
		m.logMessage(fmt.Sprintf("%s 체력: %d/%d", string(e.Current), e.Hp, config.Conf.MaxHp))
	case *game.RevealEvent:
		fmt.Println(e)
	case *game.HintEvent:
		m.cMole.hint = true
	}
	return false
}

func (m *Mole) viewStateMsg(hp int, current []rune, log []rune) error {
	username, err := m.getUserName(m.cMole.meta.Creator)
	if err != nil {
		return err
	}

	hint := m.cMole.meta.Hint
	if !m.cMole.hint {
		hint = ""
	}

	em := display.DisplayDefault(m.cMole.manID, username, hint, hp, current, log)
	me := &discordgo.MessageEdit{
		ID:      m.msgID,
		Channel: m.channelID,
		Embed:   em,
	}

	_, err = m.sess.ChannelMessageEditComplex(me)
	return err
}

func (m *Mole) viewWonMsg() error {
	username, err := m.getUserName(m.cMole.meta.Creator)
	if err != nil {
		return err
	}

	hint := m.cMole.meta.Hint

	em := display.DisplayWin(m.cMole.meta.Solution, username, hint, m.cMole.manID)
	me := &discordgo.MessageEdit{
		ID:      m.msgID,
		Channel: m.channelID,
		Embed:   em,
	}

	_, err = m.sess.ChannelMessageEditComplex(me)
	return err
}

func (m *Mole) viewDiedMsg() error {
	username, err := m.getUserName(m.cMole.meta.Creator)
	if err != nil {
		return err
	}

	hint := m.cMole.meta.Hint

	em := display.DisplayDied(m.cMole.meta.Solution, hint, m.cMole.manID, username)
	me := &discordgo.MessageEdit{
		ID:      m.msgID,
		Channel: m.channelID,
		Embed:   em,
	}

	_, err = m.sess.ChannelMessageEditComplex(me)
	return err
}

func (m *Mole) sendLogMessages() {
	msgs := m.cMole.msgs
	size := 0
	content := "```"
	for i := 0; i < len(msgs); i++ {
		content += msgs[i] + "\n"
		size += len(msgs[i])
		if size >= 500 {
			m.sess.ChannelMessageSend(m.logChannelID, content+"```")
			content = "```"
			size = 0
		}
	}
	if size > 0 {
		m.sess.ChannelMessageSend(m.logChannelID, content+"```")
	}
}

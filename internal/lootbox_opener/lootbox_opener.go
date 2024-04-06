package lootbox_opener

import (
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jagerente/discord-lootbox-opener/pkg/discord_sdk"
	"sync"
	"time"
)

var (
	ErrAlreadyRunning = fmt.Errorf("already running")
	ErrAlreadyStopped = errors.New("already stopped")
)

type Logger interface {
	Log(content string)
}

type LootboxOpener struct {
	*sync.RWMutex
	stopCh          chan struct{}
	isLoopRunning   bool
	logger          Logger
	openLootboxChan chan *discord_sdk.OpenLootboxResponse
}

func NewLootboxOpener(logger Logger) *LootboxOpener {
	return &LootboxOpener{
		RWMutex:         &sync.RWMutex{},
		stopCh:          make(chan struct{}),
		logger:          logger,
		openLootboxChan: make(chan *discord_sdk.OpenLootboxResponse),
	}
}

func (l *LootboxOpener) Run(token string, delay time.Duration) error {
	l.Lock()
	defer l.Unlock()

	l.logger.Log("Starting...")

	if l.isLoopRunning {
		return ErrAlreadyRunning
	}

	discordSDK := discord_sdk.New(&discord_sdk.Config{
		UserAgent: gofakeit.UserAgent(),
		Token:     token,
	})

	go func() {
		timer := time.NewTimer(0)
		defer timer.Stop()

		for {
			select {
			case <-l.stopCh:
				return
			case <-timer.C:
				response, err := discordSDK.OpenLootbox()
				if err != nil {
					l.logger.Log(fmt.Sprintf("Error opening lootbox: %s", err))
					l.isLoopRunning = false
					return
				}

				l.logger.Log(fmt.Sprintf("Successfully opened lootbox: %s", response.GetOpenedItemName()))

				l.openLootboxChan <- response

				timer.Reset(delay)
			}
		}
	}()
	l.isLoopRunning = true
	return nil
}

func (l *LootboxOpener) Stop() error {
	l.Lock()
	defer l.Unlock()

	l.logger.Log("Stopping...")

	if !l.isLoopRunning {
		return ErrAlreadyStopped
	}

	l.stopCh <- struct{}{}
	l.isLoopRunning = false
	return nil
}

func (l *LootboxOpener) OpenedLootboxCh() chan *discord_sdk.OpenLootboxResponse {
	return l.openLootboxChan
}

package main

import (
	"fmt"
	"github.com/jagerente/discord-lootbox-opener/internal/analytics"
	"github.com/jagerente/discord-lootbox-opener/internal/gui"
	"github.com/jagerente/discord-lootbox-opener/internal/gui/logger"
	"github.com/jagerente/discord-lootbox-opener/internal/lootbox_opener"
	"time"
)

func main() {
	const delay = 5 * time.Second
	g := gui.New(logger.New(64))

	lb := lootbox_opener.NewLootboxOpener(g)

	g.
		RegisterOnStopHandler(func() {
			if err := lb.Stop(); err != nil {
				g.Log(fmt.Sprintf("Failed to stop lootbox opener: %s", err))
				return
			}

			g.Log("Lootbox opener stopped")
		}).
		RegisterOnStartHandler(func(token string) {
			if err := lb.Run(token, delay); err != nil {
				g.Log(fmt.Sprintf("Failed to start lootbox opener: %s", err))
				return
			}
		})

	a := analytics.New(lb.OpenedLootboxCh())
	a.RegisterOnInventoryUpdateHandler(func(inventory map[string]int) {
		g.UpdateStats(inventory)
	})

	go func() {
		if err := a.Run(); err != nil {
			g.Log(fmt.Sprintf("Failed to start analytics: %s", err))
		}
	}()

	if err := g.Draw(); err != nil {
		panic(err)
	}
}

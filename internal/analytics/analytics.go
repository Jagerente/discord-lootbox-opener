package analytics

import (
	"github.com/jagerente/discord-lootbox-opener/pkg/discord_sdk"
)

type Analytics struct {
	openedLootboxCh   chan *discord_sdk.OpenLootboxResponse
	onInventoryUpdate func(inventory map[string]int)
}

func New(
	openedLootboxCh chan *discord_sdk.OpenLootboxResponse,
) *Analytics {
	return &Analytics{
		openedLootboxCh: openedLootboxCh,
	}
}

func (a *Analytics) Run() error {
	a.startOpenedLootboxConsuming()
	return nil
}

func (a *Analytics) RegisterOnInventoryUpdateHandler(handler func(map[string]int)) *Analytics {
	a.onInventoryUpdate = handler
	return a
}

func (a *Analytics) startOpenedLootboxConsuming() {
	for msg := range a.openedLootboxCh {
		if a.onInventoryUpdate != nil {
			a.onInventoryUpdate(msg.UserLootboxData.OpenedItems.TryGetNamed())
		}
	}
}

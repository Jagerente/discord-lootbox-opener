package discord_sdk

var ItemNamesByID = map[string]string{
	"1214340999644446726": "Quack!!",
	"1214340999644446724": "⮕⬆⬇⮕⬆⬇",
	"1214340999644446722": "Wump Shell",
	"1214340999644446720": "Buster Blade",
	"1214340999644446725": "Power Helmet",
	"1214340999644446723": "Speed Boost",
	"1214340999644446721": "Cute Plushie",
	"1214340999644446728": "Dream Hammer",
	"1214340999644446727": "OHHHHH BANANA",
}

type Inventory map[string]int

func (i Inventory) TryGetNamed() map[string]int {
	inventory := make(map[string]int)
	for id, count := range i {
		itemName, ok := ItemNamesByID[id]
		if ok {
			inventory[itemName] = count
		} else {
			inventory[id] = count
		}
	}

	return inventory
}

type OpenLootboxResponse struct {
	UserLootboxData struct {
		UserId        string    `json:"user_id"`
		OpenedItems   Inventory `json:"opened_items"`
		RedeemedPrize bool      `json:"redeemed_prize"`
	} `json:"user_lootbox_data"`
	OpenedItem string `json:"opened_item"`
}

func (r *OpenLootboxResponse) GetOpenedItemName() string {
	item, ok := ItemNamesByID[r.OpenedItem]
	if !ok {
		return "unknown item - " + r.OpenedItem
	}

	return item
}

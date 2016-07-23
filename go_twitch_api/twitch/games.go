package twitch

type TwitchGame struct {
	Box         TwitchGameImageBox `json:"box"`
	GiantbombId int64              `json:"giantbomb_id"`
	Name        string             `json:"name"`
	Links       TwitchGameLinkList `json:"_links"`
	Logo        TwitchGameImageBox `json:"logo"`
	Id          int64              `json:"_id"`
}
type TwitchGameLinkList struct {
	Self string `json:"self"`
	Next string `json:"next"`
}
type TwitchTopGame struct {
	Channels int64      `json:"channels"`
	Game     TwitchGame `json:"game"`
	Viewers  int64      `json:"viewers"`
}
type TwitchTopGameList struct {
	Total int64              `json:"_total"`
	Top   []TwitchTopGame    `json:"top"`
	Links TwitchGameLinkList `json:"_links"`
}
type TwitchGameImageBox struct {
	Large    string `json:"large"`
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Template string `json:"template"`
}

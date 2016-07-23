package twitch

type TwitchStreamList struct {
	Total   int64                `json:"_total"`
	Links   TwitchStreamLinkList `json:"_links"`
	Streams []TwitchStream       `json:"streams"`
}
type TwitchPreview struct {
	Large    string `json:"large"`
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Template string `json:"template"`
}
type TwitchChannel struct {
	Partner                      bool                 `json:"partner"`
	Status                       string               `json:"status"`
	DisplayName                  string               `json:"display_name"`
	Name                         string               `json:"name"`
	Language                     string               `json:"language"`
	Views                        int64                `json:"views"`
	Url                          string               `json:"url"`
	CreatedAt                    string               `json:"created_at"`
	UpdatedAt                    string               `json:"updated_at"`
	Delay                        int64                `json:"delay"`
	ProfileBanner                string               `json:"profile_banner"`
	Game                         string               `json:"game"`
	Links                        TwitchStreamLinkList `json:"_links"`
	Mature                       bool                 `json:"mature"`
	Logo                         string               `json:"logo"`
	ProfileBannerBackgroundColor string               `json:"profile_banner_background_color"`
	Followers                    int64                `json:"followers"`
	Id                           int64                `json:"_id"`
	BroadcasterLanguage          string               `json:"broadcaster_language"`
	VideoBanner                  string               `json:"video_banner"`
}
type TwitchStream struct {
	Preview     TwitchPreview        `json:"preview"`
	CreatedAt   string               `json:"created_at"`
	VideoHeight int64                `json:"video_height"`
	Game        string               `json:"game"`
	Links       TwitchStreamLinkList `json:"_links"`
	Viewers     int64                `json:"viewers"`
	AverageFps  float32              `json:"average_fps"`
	Id          int64                `json:"_id"`
	Channel     TwitchChannel        `json:"channel"`
}
type TwitchStreamLinkList struct {
	Videos        string `json:"videos"`
	Features      string `json:"features"`
	StreamKey     string `json:"stream_key"`
	Subscriptions string `json:"subscriptions"`
	Follows       string `json:"follows"`
	Self          string `json:"self"`
	Commercial    string `json:"commercial"`
	Editors       string `json:"editors"`
	Teams         string `json:"teams"`
	Chat          string `json:"chat"`
}

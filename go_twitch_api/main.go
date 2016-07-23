package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/TrungACZNE/go_twitch_api/twitch"
	"github.com/TrungACZNE/go_util/console"
	http_helpers "github.com/TrungACZNE/go_util/http"
)

const TWITCH_API = "https://api.twitch.tv/kraken"

func GetEndpoint(endpoint string, params map[string]string) string {
	return TWITCH_API + endpoint + http_helpers.QueryString(params)
}

func GetTwitchStreams(game string, max int) (*twitch.TwitchStreamList, error) {
	params := map[string]string{
		"game":  game,
		"limit": strconv.Itoa(max),
	}
	streams := &twitch.TwitchStreamList{}
	err := http_helpers.GetAndUnmarshal(GetEndpoint("/streams", params), streams)
	return streams, err
}

func GetTwitchTopGames(max int) (*twitch.TwitchTopGameList, error) {
	params := map[string]string{
		"limit": strconv.Itoa(max),
	}
	games := &twitch.TwitchTopGameList{}
	err := http_helpers.GetAndUnmarshal(GetEndpoint("/games/top", params), games)
	return games, err
}

func main() {
	games, err := GetTwitchTopGames(10)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	} else {
		choices := make([]string, 0, len(games.Top))
		for _, topGame := range games.Top {
			choices = append(choices, topGame.Game.Name)
		}
		gameName := console_helpers.GetUserChoiceFromMenu("Select one of the top games: ", choices)
		streams, err := GetTwitchStreams(gameName, 25)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		} else {
			fmt.Println(streams.Streams[0].Channel.Name)
		}

	}
}

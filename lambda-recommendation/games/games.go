package games

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GamesJSON struct {
	Response struct {
		UserId    string `json:"user_id"`
		GameCount int    `json:"game_count"`
		Games     []struct {
			Appid           int `json:"appid"`
			PlaytimeForever int `json:"playtime_forever"`
			Playtime2Weeks  int `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}

type Friendslist struct {
	Friends []struct {
		Steamid      string `json:"steamid"`
		Relationship string `json:"relationship"`
		FriendSince  int    `json:"friend_since"`
	} `json:"friends"`
}

type Game struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string `json:"type"`
		Name                string `json:"name"`
		SteamAppid          int    `json:"steam_appid"`
		RequiredAge         int    `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		Dlc                 []int  `json:"dlc"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		SupportedLanguages  string `json:"supported_languages"`
		Reviews             string `json:"reviews"`
		HeaderImage         string `json:"header_image"`
		Website             string `json:"website"`
		View                int    `json:"view"`
		PcRequirements      struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"`
		MacRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"mac_requirements"`
		LinuxRequirements []interface{} `json:"linux_requirements"`
		Developers        []string      `json:"developers"`
		Publishers        []string      `json:"publishers"`
		Demos             []struct {
			Appid       int    `json:"appid"`
			Description string `json:"description"`
		} `json:"demos"`
		PriceOverview struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             int    `json:"display_type"`
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs                    []struct {
				Packageid                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        string `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int    `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		ContentDescriptors struct {
			Ids   []interface{} `json:"ids"`
			Notes interface{}   `json:"notes"`
		} `json:"content_descriptors"`
	} `json:"data"`
}

func SliceUniqMap(s []int) []int {
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}
func Get() {
	var ApiKey = "0DF75B7B51BAF08BB992D5327748BFB8"

	var session, err = mgo.Dial("localhost")
	var conn = session.DB("recommendation").C("games")
	var connDetail = session.DB("recommendation").C("detail")
	var connFriends = session.DB("recommendation").C("friends")

	conn.DropCollection()
	connDetail.DropCollection()
	var games []int

	if err != nil {
		fmt.Println("Connection problem")
		os.Exit(1)
	}

	var list []Friendslist
	connFriends.Find(nil).All(&list)

	friends := list[0].Friends

	for _, v := range friends {

		var url = strings.Join([]string{"http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=", ApiKey, "&steamid=", v.Steamid, "&format=json"}, "")

		var target GamesJSON
		r, err := http.Get(url)
		if err != nil {
			fmt.Printf("Server not responding %s", err.Error())
			os.Exit(1)
		}

		t, _ := ioutil.ReadAll(r.Body)

		json.Unmarshal(t, &target)
		target.Response.UserId = v.Steamid
		for _, k := range target.Response.Games {
			games = append(games, k.Appid)
		}

		err = conn.Insert(target.Response)

		if err != nil {
			fmt.Printf("Server not responding %s", err.Error())
			os.Exit(1)

		}
		defer r.Body.Close()
	}
	defer func() {
		games = SliceUniqMap(games)
		fmt.Println("Friends information completed")
		fmt.Println("Games count:", len(games))
		for _, t := range games {
			time.Sleep(2000 * time.Millisecond)
			gameDetailUrl := strings.Join([]string{"http://store.steampowered.com/api/appdetails?appids=", strconv.Itoa(t), "&cc=us&l=en"}, "")
			var target map[string]Game
			r, err := http.Get(gameDetailUrl)
			if err != nil {
				fmt.Printf("Server not responding %s", err.Error())
				os.Exit(1)
			}

			t, _ := ioutil.ReadAll(r.Body)

			json.Unmarshal(t, &target)

			for _, v := range target {
				if v.Data.Type == "game" {
					err = connDetail.Insert(v.Data)
				}
			}

			if err != nil {
				fmt.Printf("Server not responding %s", err.Error())
				os.Exit(1)
			}

			defer r.Body.Close()
		}
	}()
}
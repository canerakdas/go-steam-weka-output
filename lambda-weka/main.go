package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"text/template"
	"os"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

type Response struct {
	UserId    string `json:"user_id"`
	GameCount int    `json:"game_count"`
	Games     []struct {
		Appid           int `json:"appid"`
		PlaytimeForever int `json:"playtime_forever"`
		Playtime2Weeks  int `json:"playtime_2weeks,omitempty"`
	} `json:"games"`
}

type Friendslist struct {
	Friends []struct {
		Steamid      string `json:"steamid"`
		Relationship string `json:"relationship"`
		FriendSince  int    `json:"friend_since"`
	} `json:"friends"`
}


	type Data    struct {
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
	}

	type ResponseData struct {
	UserId    string `json:"user_id"`
	GameCount int    `json:"game_count"`
	Games     []struct {
		UserId    string `json:"user_id"`
		Appid           int `json:"appid"`
		GameDetail			Data
		PlaytimeForever int `json:"playtime_forever"`
		Playtime2Weeks  int `json:"playtime_2weeks,omitempty"`
	} `json:"games"`
}

func main() {
	var session, err = mgo.Dial("localhost")
	var conn = session.DB("recommendation").C("games")
	var connDetail = session.DB("recommendation").C("detail")

	var games []ResponseData
	var game Data
	
	err = conn.Find(nil).All(&games)
	for j,_ := range games{
	for i,_ := range games[j].Games{
		id := games[j].Games[i].Appid

		err = connDetail.Find(bson.M{"steamappid": id}).One(&game)
		game.ReleaseDate.Date = strings.Replace(game.ReleaseDate.Date,",","-",-1)
		for i,_ := range game.Publishers{
			game.Publishers[i] = strings.Replace(game.Publishers[i],",","",-1)
			game.Publishers[i] = strings.Replace(game.Publishers[i],"'","",-1)
			game.Publishers[i] = strings.Replace(game.Publishers[i],"\"","",-1)
			game.Publishers[i] = strings.Replace(game.Publishers[i],"	","",-1)
		}
		for i,_ := range game.Developers{
			game.Developers[i] = strings.Replace(game.Developers[i],",","",-1)
			game.Developers[i] = strings.Replace(game.Developers[i],"'","",-1)
			game.Developers[i] = strings.Replace(game.Developers[i],"\"","",-1)
			game.Developers[i] = strings.Replace(game.Developers[i],"	","",-1)
		}
		games[j].Games[i].GameDetail = game
		games[j].Games[i].UserId = games[j].UserId
	}
	}

	tmpl := template.Must(template.ParseFiles("layout.csv"))
	err = tmpl.Execute(os.Stdout, games)
	
	if err != nil {
		fmt.Println(err,"Whoops")
	}
}

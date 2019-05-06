package friends

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FriendJSON struct {
	Friendslist struct {
		Friends []struct {
			Steamid      string `json:"steamid"`
			Relationship string `json:"relationship"`
			FriendSince  int    `json:"friend_since"`
		} `json:"friends"`
	} `json:"friendslist"`
}

func Get() {

	var ApiKey = "0DF75B7B51BAF08BB992D5327748BFB8"
	var UserId = "76561198072752123"

	var session, err = mgo.Dial("localhost")
	var conn = session.DB("recommendation").C("friends")
	var url = strings.Join([]string{"http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=", ApiKey, "&steamid=", UserId, "&relationship=friend"}, "")
	if err != nil {
		fmt.Println("Connection problem")
		os.Exit(1)
	}

	var target FriendJSON
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)
	}

	t, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(t, &target)

	err = conn.Insert(target.Friendslist)

	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)

	}

	defer r.Body.Close()
}

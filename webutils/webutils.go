package webutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gopherworks/bawt"
	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
)

type Utils struct {
	bot *bawt.Bot
}

func init() {
	bawt.RegisterPlugin(&Utils{})
}

func (utils *Utils) InitWebPlugin(bot *bawt.Bot, privRouter *mux.Router, pubRouter *mux.Router) {
	utils.bot = bot
	privRouter.HandleFunc("/slack/channels", utils.handleGetChannels)
	privRouter.HandleFunc("/slack/users", utils.handleGetUsers)
}

func (utils *Utils) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	out := struct {
		Users map[string]slack.User `json:"users"`
	}{
		Users: utils.bot.Users,
	}

	err := enc.Encode(out)
	if err != nil {
		webReportError(w, "Error encoding JSON", err)
		return
	}
	return
}

func (utils *Utils) handleGetChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	out := struct {
		Channels map[string]bawt.Channel `json:"channels"`
	}{
		Channels: utils.bot.Channels,
	}

	err := enc.Encode(out)
	if err != nil {
		webReportError(w, "Error encoding JSON", err)
		return
	}
	return
}

func webReportError(w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("%s\n\n%s\n", msg, err)))
}

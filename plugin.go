package bawt

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/slack-go/slack"
)

//
// Bot plugins
//

// Plugin describes the generic bot plugin
type Plugin interface{}

// PluginInitializer describes the interface is used to check which plugins
// can be initialized during plugin initalization initChatPlugins
type PluginInitializer interface {
	InitPlugin(*Bot)
}

// WebServer describes the interface for webserver plugins
type WebServer interface {
	// Used internally by the `bawt` library.
	InitWebServer(*Bot, []string)
	RunServer()

	// Used by an Auth provider.
	SetAuthMiddleware(func(http.Handler) http.Handler)
	SetAuthenticatedUserFunc(func(req *http.Request) (*slack.User, error))

	// Can be called by any plugins.
	PrivateRouter() *mux.Router
	PublicRouter() *mux.Router
	GetSession(*http.Request) *sessions.Session
	AuthenticatedUser(*http.Request) (*slack.User, error)
}

// WebPlugin initializes plugins with a `Bot` instance, a `privateRouter` and a `publicRouter`. All URLs handled by the `publicRouter` must start with `/public/`.
type WebPlugin interface {
	InitWebPlugin(bot *Bot, private *mux.Router, public *mux.Router)
}

// WebServerAuth returns a middleware warpping the passed on `http.Handler`. Only one auth handler can be added.
type WebServerAuth interface {
	InitWebServerAuth(bot *Bot, webserver WebServer)
}

var registeredPlugins = make([]Plugin, 0)

// RegisterPlugin adds the provided Plugin to the list of registered plugins
func RegisterPlugin(plugin Plugin) {
	registeredPlugins = append(registeredPlugins, plugin)
}

// RegisteredPlugins returns the list of registered plugins
func RegisteredPlugins() []Plugin {
	return registeredPlugins
}

func initChatPlugins(bot *Bot) {
	for _, plugin := range registeredPlugins {
		chatPlugin, ok := plugin.(PluginInitializer)
		if ok {
			chatPlugin.InitPlugin(bot)
		}
	}
}

func initWebServer(bot *Bot, enabledPlugins []string) {
	for _, plugin := range registeredPlugins {
		webServer, ok := plugin.(WebServer)
		if ok {
			webServer.InitWebServer(bot, enabledPlugins)
			bot.WebServer = webServer
			return
		}
	}
}

func initWebPlugins(bot *Bot) {
	if bot.WebServer == nil {
		return
	}

	for _, plugin := range registeredPlugins {
		if webPlugin, ok := plugin.(WebPlugin); ok {
			webPlugin.InitWebPlugin(bot, bot.WebServer.PrivateRouter(), bot.WebServer.PublicRouter())
		}

		count := 0
		if webServerAuth, ok := plugin.(WebServerAuth); ok {
			count++

			if count > 1 {
				log.Fatalln("Can not load two WebServerAuth plugins. Already loaded one.")
			}
			webServerAuth.InitWebServerAuth(bot, bot.WebServer)
		}

	}
}

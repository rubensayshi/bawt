package totw

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gopherworks/bawt"
)

type Totw struct {
	bot *bawt.Bot
}

func init() {
	bawt.RegisterPlugin(&Totw{})
}

func (totw *Totw) InitPlugin(bot *bawt.Bot) {
	bawt.RegisterStringList("useless techs", []string{
		"http://i.minus.com/ib2bUNs2W1CI1V.gif",
		"http://media.giphy.com/media/anl0wydLNhKus/giphy.gif",
		"http://www.ptc.dcs.edu/Moody/comphistory/cavemanwriting.gif",
		"http://i.imgur.com/VbzhAbd.gif",
		"http://www.patrickcarrion.com/wp-content/uploads/2014/05/mowingdressgif.gif",
		"http://cdn.shopify.com/s/files/1/0243/7593/products/MKSB023_UselessMachine_Animation_grande.gif",
		"http://i.imgur.com/CRuLGek.gif",
		"http://i.imgur.com/EteBF9K.gif",
		"http://www.ohmagif.com/wp-content/uploads/2011/12/useless-invention.gif",
		"http://i3.kym-cdn.com/photos/images/original/000/495/044/9b8.gif",
		"http://uproxx.files.wordpress.com/2012/09/iron.gif",
	})
	bawt.RegisterStringList("tech adept", []string{
		"you're a real tech adept",
		"what an investigator",
		"such deep search!",
		"a real innovator you are",
		"way to go, I'm impressed",
		"hope it's better than my own code",
		"noted, but are you sure it's good ?",
		"I'll take a look into this one",
		"you're generous!",
		"hurray!",
	})

	totw.bot = bot

	go totw.ScheduleAlerts(bot.Config.GeneralChannel, time.Thursday, 16, 0)

	bot.Listen(&bawt.Listener{
		MessageHandlerFunc: totw.ChatHandler,
	})
}

func (totw *Totw) ChatHandler(conv *bawt.Listener, msg *bawt.Message) {
	if strings.HasPrefix(msg.Text, "!totw") || strings.HasPrefix(msg.Text, "!techoftheweek") {
		msg.ReplyMention(bawt.RandomString("tech adept"))
	}
}

func (totw *Totw) ScheduleAlerts(channel string, w time.Weekday, hour, min int) {
	for {
		next, when := bawt.NextWeekdayTime(time.Now(), w, hour, min)
		log.Println("TOTW: Next occurrence: ", next)

		<-time.After(when)

		totw.bot.SendToChannel(channel, bawt.RandomString("useless techs"))
		totw.bot.SendToChannel(channel, `Time for some Tech of the Week! What's your pick ?  Start your line with "!techoftheweek"`)
	}
}

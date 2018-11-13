package recognition

import (
	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/gopherworks/bawt"
)

type Plugin struct {
	bot    *bawt.Bot
	config Config
	store  Store
}

func init() {
	bawt.RegisterPlugin(&Plugin{})
}

func (p *Plugin) InitPlugin(bot *bawt.Bot) {
	p.bot = bot

	err := bot.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		log.Fatalln("Couldn't create the `recognition` bucket")
	}

	var conf struct {
		Recognition Config
	}
	bot.LoadConfig(&conf)
	p.config = conf.Recognition

	p.store = &boltStore{db: bot.DB}

	p.listenRecognize()
	p.listenUpvotes()
}

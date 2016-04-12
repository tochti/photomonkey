package photomonkey

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/tochti/chief"

	"bitbucket.org/mrd0ll4r/tbotapi"
)

var (
	AppName = "photomonkey"
)

type (
	Specs struct {
		Token string `required:"true" envconfig:"TOKEN"`
	}
)

func Start(token string) {

	bot, err := tbotapi.New(token)
	if err != nil {
		log.Fatal(err)
	}

	//router := new(photomonkey.Router)
	//router.HandleFunc().IsPhoto()

	c := chief.New(5, decoder(messageHandler))
	c.Start()

	for {
		select {
		case update := <-bot.Updates:
			if update.Error() != nil {
				log.Println(update.Error())
				continue
			}
			c.Jobs <- chief.Job{Order: update.Update()}
		}
	}

}

func ReadSpecs() (*Specs, error) {
	s := &Specs{}
	err := envconfig.Process(AppName, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func messageHandler(update tbotapi.Update) {
	msg := update.Message
	log.Println(*msg.Text)
}

func decoder(fn func(tbotapi.Update)) chief.HandleFunc {
	return func(j chief.Job) {
		update, ok := j.Order.(tbotapi.Update)
		if !ok {
			log.Println("Error in decoder func")
			return
		}
		fn(update)
	}
}

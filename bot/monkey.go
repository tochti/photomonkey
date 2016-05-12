package bot

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/kelseyhightower/envconfig"
	"github.com/tochti/chief"

	"bitbucket.org/mrd0ll4r/tbotapi"
)

var (
	AppName            = "photomonkey"
	ErrMissingToken    = errors.New("Missing TOKEN env")
	ErrMissingImageDir = errors.New("Missing IMAGE_DIR env")
)

type (
	Specs struct {
		Token    string `required:"true" envconfig:"TOKEN"`
		ImageDir string `required:"true" envconfig:"IMAGE_DIR"`
	}

	handler interface {
		HandleUpdate(tbotapi.Update)
	}

	photoHandler struct {
		Bot        *tbotapi.TelegramBotAPI
		HTTPClient *http.Client
		Token      string
		ImageDir   string
	}
)

func Start(token string, imageDir string) {

	log.Println("Monkey is running....")

	bot, err := tbotapi.New(token)
	if err != nil {
		log.Fatal(err)
	}

	photoHandler := &photoHandler{
		Bot:        bot,
		HTTPClient: http.DefaultClient,
		Token:      token,
		ImageDir:   imageDir,
	}

	c := chief.New(5, decodeJob(photoHandler))
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

	if s.Token == "" {
		return nil, ErrMissingToken
	} else if s.ImageDir == "" {
		return nil, ErrMissingImageDir
	}

	return s, nil
}

func decodeJob(h handler) chief.HandleFunc {
	return func(j chief.Job) {
		update, ok := j.Order.(tbotapi.Update)
		if !ok {
			log.Println("Error in decoder func")
			return
		}
		h.HandleUpdate(update)
	}
}

func (h *photoHandler) HandleUpdate(update tbotapi.Update) {
	msg := update.Message
	if msg.Type() != tbotapi.PhotoMessage {
		return
	}

	err := h.HandlePhoto(msg)
	if err != nil {
		log.Println(err)
	}
}

func (h *photoHandler) HandlePhoto(message *tbotapi.Message) error {
	photo := findBiggestPhoto(message.Photo)
	fileID := photo.ID
	log.Printf("Receive image with id=%v\n", fileID)

	botResp, err := h.Bot.GetFile(fileID)
	if err != nil {
		return err
	}

	filePath := botResp.File.Path
	photoURL := fmt.Sprintf("https://api.telegram.org/file/bot%v/%v", h.Token, filePath)

	log.Println(photoURL)

	resp, err := h.HTTPClient.Get(photoURL)
	if err != nil {
		return err
	}

	ext := path.Ext(filePath)
	photoName := fileID + ext
	photoPath := path.Join(h.ImageDir, photoName)

	if _, err := os.Stat(photoPath); os.IsExist(err) {
		return fmt.Errorf("%v already exists\n", photoPath)
	}

	fh, err := os.Create(photoPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(fh, resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Store photo in %v\n", photoPath)
	saveCaption(message, photoName)

	return nil
}

func findBiggestPhoto(tmp *[]tbotapi.PhotoSize) tbotapi.PhotoSize {
	photos := *tmp
	biggest := photos[0].Width * photos[0].Height
	biggestPhoto := 0
	for i, photo := range photos[1:] {
		size := photo.Width * photo.Height
		if size >= biggest {
			biggest = size
			biggestPhoto = i
		}
	}

	return photos[biggestPhoto]
}

func saveCaption(message *tbotapi.Message, photoName string) {
}

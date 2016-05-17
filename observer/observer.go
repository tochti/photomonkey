package observer

import "github.com/tochti/photomonkey/database"

type (
	PhotoObservers struct {
		observers []chan database.Photo
	}
)

func (o *PhotoObservers) Add(c chan database.Photo) {
	o.observers = append(o.observers, c)
}

func (o *PhotoObservers) Broadcast(photo database.Photo) {
	for _, observer := range o.observers {
		go func(c chan database.Photo) {
			c <- photo
		}(observer)
	}
}

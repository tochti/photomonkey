package observer

import (
	"sync"

	"github.com/tochti/photomonkey/database"
)

type (
	PhotoObservers struct {
		mutex     *sync.Mutex
		observers map[int64]chan database.Photo
	}
)

func New() *PhotoObservers {
	return &PhotoObservers{
		mutex:     new(sync.Mutex),
		observers: map[int64]chan database.Photo{},
	}
}

func (o *PhotoObservers) Add(c chan database.Photo) int64 {
	defer o.mutex.Unlock()
	o.mutex.Lock()
	next := int64(len(o.observers) + 1)
	o.observers[next] = c
	return next
}

func (o *PhotoObservers) Remove(id int64) {
	defer o.mutex.Unlock()
	o.mutex.Lock()
	close(o.observers[id])
	delete(o.observers, id)
}
func (o *PhotoObservers) Broadcast(photo database.Photo) {
	for _, observer := range o.observers {
		go func(c chan database.Photo) {
			c <- photo
		}(observer)
	}
}

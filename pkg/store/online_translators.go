package store

import (
	"sync"
	"time"
)

type OnlineTranslatorsStore struct {
	translator map[string]int64
	l          sync.Mutex
}

func NewOnlineTranslatorsStore() *OnlineTranslatorsStore {
	return &OnlineTranslatorsStore{
		translator: make(map[string]int64),
	}
}

func (otStore *OnlineTranslatorsStore) CleanOnlineTranslators() {
	otStore.l.Lock()
	for pseudo, deletedAt := range otStore.translator {
		if deletedAt < time.Now().Unix() {
			delete(otStore.translator, pseudo)
		}
	}
	otStore.l.Unlock()
}

func (otStore *OnlineTranslatorsStore) Put(pseudo string) {
	otStore.l.Lock()
	otStore.translator[pseudo] = time.Now().Add(time.Duration(17 * time.Minute)).Unix()
	otStore.l.Unlock()
}

func (otStore *OnlineTranslatorsStore) TotalOnlineTranslators() int {
	otStore.l.Lock()
	totalOnlineTranslators := len(otStore.translator)
	otStore.l.Unlock()
	return totalOnlineTranslators
}

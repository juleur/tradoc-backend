package store

import (
	"fmt"
	"sync"
	"time"
)

type OnGoingTranslation struct {
	translatorID int
	texteID      int
	deleteAt     int64
}

type OnGoingTranslationsStore struct {
	m map[string][]OnGoingTranslation
	l sync.Mutex
}

func NewOnGoingTranslationsStore() *OnGoingTranslationsStore {
	onGoingTranslationsStore := &OnGoingTranslationsStore{
		m: make(map[string][]OnGoingTranslation),
	}
	return onGoingTranslationsStore
}

func (ogtStore *OnGoingTranslationsStore) CleanOnGoingTranslations() <-chan string {
	translationDeletedCh := make(chan string, 10)
	go func() {
		ogtStore.l.Lock()
		for dialect, onGoingTranslations := range ogtStore.m {
			if len(onGoingTranslations) == 0 {
				continue
			}
			newOnGoingTranslations := []OnGoingTranslation{}
			for _, onGoingTranslation := range onGoingTranslations {
				if onGoingTranslation.deleteAt < time.Now().Unix() {
					translationDeletedCh <- fmt.Sprintf("%s:%d", dialect, onGoingTranslation.texteID)
					continue
				}
				newOnGoingTranslations = append(newOnGoingTranslations, onGoingTranslation)
			}
			ogtStore.m[dialect] = make([]OnGoingTranslation, len(newOnGoingTranslations))
			ogtStore.m[dialect] = newOnGoingTranslations
		}
		ogtStore.l.Unlock()
		close(translationDeletedCh)
	}()
	return translationDeletedCh
}

func (ogtStore *OnGoingTranslationsStore) LimitRate(dialect string, translatorID int) bool {
	ogtStore.l.Lock()
	onGoingTranslations, ok := ogtStore.m[dialect]
	ogtStore.l.Unlock()
	if !ok {
		return false
	}
	totalOnGoingTranslations := countHowManyOnGoingTranslations(onGoingTranslations, translatorID)
	return totalOnGoingTranslations > 300
}

func (ogtStore *OnGoingTranslationsStore) Delete(dialect string, texteID int) {
	ogtStore.l.Lock()
	onGoingTranslations, ok := ogtStore.m[dialect]
	if !ok {
		return
	}
	newOnGoingTranslations := []OnGoingTranslation{}
	for _, onGoingTranslation := range onGoingTranslations {
		if onGoingTranslation.texteID == texteID {
			continue
		}
		newOnGoingTranslations = append(newOnGoingTranslations, onGoingTranslation)
	}
	ogtStore.m[dialect] = make([]OnGoingTranslation, len(newOnGoingTranslations))
	ogtStore.m[dialect] = newOnGoingTranslations
	ogtStore.l.Unlock()
}

func (ogtStore *OnGoingTranslationsStore) Put(dialect string, translatorID int, texteIDs []int) {
	ogtStore.l.Lock()
	onGoingTranslations := ogtStore.m[dialect]

	for _, texteID := range texteIDs {
		if isItExist(onGoingTranslations, texteID) {
			return
		}

		newOnGoingTranslation := OnGoingTranslation{
			translatorID: translatorID,
			texteID:      texteID,
			deleteAt:     time.Now().Add(time.Duration(2 * time.Hour)).Unix(),
		}

		onGoingTranslations = append(onGoingTranslations, newOnGoingTranslation)
		ogtStore.m[dialect] = onGoingTranslations
	}

	ogtStore.l.Unlock()
}

func (ogtStore *OnGoingTranslationsStore) GetTexteIDs(dialect string) []int {
	ogtStore.l.Lock()
	onGoingTranslations, ok := ogtStore.m[dialect]
	ogtStore.l.Unlock()
	if !ok {
		return []int{}
	}
	return extractTexteIDs(onGoingTranslations)
}

func isItExist(onGoingTranslations []OnGoingTranslation, texteID int) bool {
	for _, ogt := range onGoingTranslations {
		if ogt.texteID == texteID {
			return true
		}
	}
	return false
}

func extractTexteIDs(onGoingTranslations []OnGoingTranslation) []int {
	var IDs []int
	for _, onGoingTranslation := range onGoingTranslations {
		IDs = append(IDs, onGoingTranslation.texteID)
	}
	return IDs
}

func countHowManyOnGoingTranslations(onGoingTranslations []OnGoingTranslation, translatorID int) int {
	var counter int
	for _, onGoingTranslation := range onGoingTranslations {
		if onGoingTranslation.translatorID == translatorID {
			counter += 1
		}
	}
	return counter
}

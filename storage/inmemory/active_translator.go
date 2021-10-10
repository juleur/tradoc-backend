package inmemory

import (
	"sync"
	"time"
)

const DELAY time.Duration = 10 * time.Minute

type ActiveTranslators struct {
	translatorLastSeen map[string]time.Time
	sync.RWMutex
}

func NewActiveTranslatorsTracker() *ActiveTranslators {
	activeTranslators := ActiveTranslators{
		translatorLastSeen: make(map[string]time.Time),
	}

	go activeTranslators.removeInactiveTranslator()

	return &activeTranslators
}

// removeInactiveTranslator removes translator who is inactive every 10 minutes
func (at *ActiveTranslators) removeInactiveTranslator() {
	for {
		<-time.After(DELAY)

		for translatorID, lastSeen := range at.translatorLastSeen {
			now := time.Now()
			// remove translatorID from the map after 10 minutes of inactivity
			if now.After(lastSeen.Add(DELAY)) {
				at.Delete(translatorID)
			}
		}
	}
}

// AddOrKeepActive adds new translator to the list of active translators or keep him/she active
func (at *ActiveTranslators) AddOrKeepActive(translatorID string) {
	at.RWMutex.RLock()
	at.translatorLastSeen[translatorID] = time.Now()
	at.RWMutex.RUnlock()
}

// Delete deletes a given translator from the list of active translators
func (at *ActiveTranslators) Delete(translatorID string) {
	at.RWMutex.RLock()
	delete(at.translatorLastSeen, translatorID)
	at.RWMutex.RUnlock()
}

// Total returns how many translators are active
func (at *ActiveTranslators) Total() int {
	at.RWMutex.RLock()
	defer at.RWMutex.RUnlock()

	return len(at.translatorLastSeen)
}

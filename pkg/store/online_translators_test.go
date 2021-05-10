package store

import (
	"testing"
	"time"
)

func TestCleanOnlineTranslators(t *testing.T) {
	onlineTranslatorsStore := NewOnlineTranslatorsStore()

	onlineTranslatorsStore.Put("Test1")
	onlineTranslatorsStore.Put("teSt2")

	t.Log(onlineTranslatorsStore.TotalOnlineTranslators())
	time.Sleep(30 * time.Second)

	onlineTranslatorsStore.Put("tEST3")
	onlineTranslatorsStore.Put("TeSt4")

	t.Log(onlineTranslatorsStore.TotalOnlineTranslators())
	time.Sleep(32 * time.Second)

	onlineTranslatorsStore.CleanOnlineTranslators()

	t.Log(onlineTranslatorsStore.TotalOnlineTranslators())
}

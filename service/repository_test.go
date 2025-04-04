package service

import "testing"

func TestStoreSet(t *testing.T) {
	store := NewStore()
	chatID := int64(1)
	round := Round{
		SpyID:   0,
		Place:   Place{"school"},
		Members: 10,
		Roles:   nil,
	}
	store.Set(chatID, round)
	result, ok := store.Get(chatID)
	if !ok {
		t.Errorf("store.Get(%d): expected true, got false", chatID)
	}
	expected := round
	if !result.Equal(expected) {
		t.Errorf("Expected %s; actual %s", result, expected)
	}
}

package main

import "testing"

func TestNewLogger(t *testing.T) {
	log, err := NewLogger("access.log", L_DEBUG)
	if err != nil {
		t.Error(err)
		return
	}
	log.Debug("a", "123", "d", "123")
	log.Info("a", "123", "d", "123")
	log.Warning("a", "123", "d", "123")
	log.Error("a", "123", "d", "123")
	log.Close()
}

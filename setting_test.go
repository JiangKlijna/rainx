package main

import "testing"

func TestNewSetting(t *testing.T) {
	setting, err := NewSetting("setting.json")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(setting)
	err = setting.IsValid()
	if err != nil {
		t.Error(err)
		return
	}
}

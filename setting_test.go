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
	for _, s := range setting.Servers() {
		t.Log(s.Listen())
		for _, l := range s.Locations() {
			t.Log(l)
		}
	}
}

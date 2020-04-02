package main

import (
	"bufio"
	"github.com/bwmarrin/discordgo"
	"reflect"
	"strings"
	"testing"
)

func Test_prompt(t *testing.T) {
	tests := []struct {
		name string
		f    *strings.Reader
	}{
		{
			"One line file",
			strings.NewReader("Hello World!"),
		}, {
			"Multi line file",
			strings.NewReader("Hello\nWorld!"),
		}, {
			"Empty file",
			strings.NewReader(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prompt(tt.f); reflect.DeepEqual(got, "") {
				scanner := bufio.NewScanner(tt.f)
				if scanner.Scan() {
					t.Errorf("prompt() = %v, want \"\"", got)
				}
			}
		})
	}
}

func Test_promptJournal(t *testing.T) {
	type args struct {
		s    *discordgo.Session
		chID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

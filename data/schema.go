package data

import "gorm.io/gorm"

type Schema struct {
	gorm.Model
	TeamName      string
	TeamNumber    int
	MatchNumber   int
	AutoAmps      int
	AutoSpeaker   int
	AutoLeave     bool
	AutoMiddle    bool
	TeleopAmps    int
	TeleopSpeaker int
	Chain         bool
	Harmony       bool
	Trap          bool
	Park          bool
	Ground        bool
	Feeder        bool
	LLVm          string
	Defense       bool
	Notes         string
}

package data

import "gorm.io/gorm"

type Schema struct {
	gorm.Model           //00
	TeamName      string //01
	TeamNumber    int    //02
	MatchNumber   int    //03
	AutoAmps      int    //04
	AutoSpeaker   int    //05
	AutoLeave     bool   //06
	AutoMiddle    bool   //07
	TeleopAmps    int    //08
	TeleopSpeaker int    //09
	Chain         bool   //10
	Harmony       bool   //11
	Trap          bool   //12
	Park          bool   //13
	Ground        bool   //14
	Feeder        bool   //15
	Mobility      bool   //16
	Penalties     int    //17
	TechPenalties int    //18
	GroundPickup  bool   //19
	StartingPos   int    //20
	Defense       bool   //21
	CenterRing    bool   //22
	Notes         string //23

}

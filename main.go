package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main/data"
	"reflect"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	//parse sqldb
	if err != nil {
		panic("Could not connect to DB")
	}

	allData := []data.Schema{}
	db.Select("*").Find(&allData)
	app := app.New()
	current := app.NewWindow("TKO Crescendo Tracker (patented)")
	current.Resize(fyne.NewSize(800, 600))
	//lvm
	llvm := widget.NewTable(
		func() (int, int) {
			val := allData[0]
			//assuming exsits and not acesing memaddr with nothing allloca egfault
			v := reflect.ValueOf(val)

			return len(allData), v.NumField()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("test")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			v := reflect.ValueOf(allData[i.Row])
			//expensive
			fields := make([]interface{}, reflect.ValueOf(v).NumField()) //types
			for i := 0; i < v.NumField(); i++ {
				fields[i] = v.Field(i).Interface() //reflection my belobed
			}

			o.(*widget.Label).SetText(fields[i.Col].(string))
		},
	)

	mainContainer := container.NewVSplit(llvm, widget.NewButton("test", func() {}))

	//migrate data and term[late
	_ = db.AutoMigrate(&data.Schema{})
	//expressions ast
	//db.Raw("GET WHERE name = 1")
	db.Create(&data.Schema{TeamName: "test", TeamNumber: 1, MatchNumber: 1, AutoAmps: 1, AutoSpeaker: 1, AutoLeave: true, AutoMiddle: true, TeleopAmps: 1, TeleopSpeaker: 1, Chain: true, Harmony: true, Trap: true, Park: true, Ground: true, Feeder: true, LLVm: "test", Defense: true, Notes: "test"})
	current.SetContent(mainContainer)
	current.ShowAndRun() //defer?
	//llvm go sadck jwt auth
	//gorm orm sql wrapper

}

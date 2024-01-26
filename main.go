package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main/data"
	"reflect"
	"strconv"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	//parse sqldb
	if err != nil {
		panic("Could not connect to DB")
	}

	allData := []data.Schema{}
	tcp := [][]string{}
	tcp = append(tcp, []string{"TeamName", "TeamNumber", "MatchNumber", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"}) //hardcoded since migrated
	db.Select("*").Find(&allData)
	for _, val := range allData {
		v := reflect.ValueOf(val)
		//expensive
		fields := make([]interface{}, v.NumField()) //types
		fmt.Println("fields: ", fields)
		for j := 0; j < v.NumField(); j++ {
			fields[j] = v.Field(j).Interface() //reflection my belobed
		}
		//opengl reference
		//mmap syscall
		vala := []string{}
		for _, tcppacket := range fields {

			if reflect.TypeOf(tcppacket) == reflect.TypeOf(gorm.Model{}) {
				vala = append(vala, strconv.Itoa(int(tcppacket.(gorm.Model).ID)))
			} else {
				vala = append(vala, fmt.Sprintf("%v", tcppacket))
			}

		}
		tcp = append(tcp, vala)

	}

	fmt.Println("llvm", tcp)
	app := app.New()
	current := app.NewWindow("TKO Crescendo Tracker (patented)")
	current.Resize(fyne.NewSize(1200, 600))
	fmt.Println(allData)
	//lvm
	llvm := widget.NewTable(
		func() (int, int) {

			//assuming exsits and not acesing memaddr with nothing allloca egfault

			return len(tcp), len(tcp[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("test")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(tcp[i.Row][i.Col]) //nop
		},
	)

	//could change the code st
	for i := 0; i < len(tcp); i++ {
		llvm.SetColumnWidth(i, 100)
	}

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

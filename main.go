package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main/data"
	"reflect"
	"strconv"
	"strings"
)

var headers map[string]int = map[string]int{
	"TeamName":      1,
	"TeamNumber":    2,
	"MatchNumber":   3,
	"AutoAmps":      4,
	"AutoSpeaker":   5,
	"AutoLeave":     6,
	"AutoMiddle":    7,
	"TeleopAmps":    8,
	"TeleopSpeaker": 9,
	"Chain":         10,
	"Harmony":       11,
	"Trap":          12,
	"Park":          13,
	"Ground":        14,
	"Feeder":        15,
	"LLVm":          16,
	"Defense":       18,
	"Notes":         19,
}

func populate(db *gorm.DB, allData []data.Schema, tcp [][]string, fields []string) ([][]string, []data.Schema) {
	//reflection no tneeded
	//by reference?
	tcp = [][]string{}
	tcp = append(tcp, fields) //hardcoded since migrated
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
		//supposde to be emtpy fields to be filled populated
		for _, tcppacket := range fields {

			if reflect.TypeOf(tcppacket) == reflect.TypeOf(gorm.Model{}) {
				vala = append(vala, strconv.Itoa(int(tcppacket.(gorm.Model).ID)))
			} else {
				vala = append(vala, fmt.Sprintf("%v", tcppacket))
			}

		}
		tcp = append(tcp, vala)

	}
	return tcp, allData

}

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	//parse sqldb
	if err != nil {
		panic("Could not connect to DB")
	}

	allData := []data.Schema{}
	tcp := [][]string{}

	averageLabel := widget.NewLabel("Average: 0")
	teamChoose := widget.NewEntry()
	averageChoose := widget.NewSelect([]string{"AutoAmps", "AutoSpeaker", "TeleopAmps", "TeleopSpeaker"}, func(s string) {
		total := 0
		amnt := 0
		for _, val := range tcp {
			fmt.Println(val[1] == teamChoose.Text)
			if val[1] == teamChoose.Text {
				amnt++
				fmt.Println("val", val[headers[s]])
				vala, err := strconv.Atoi(val[headers[s]])
				fmt.Println(err)
				if err == nil {
					total += vala
				}
			}

		}
		if amnt == 0 {
			amnt = 1
		}

		averageLabel.SetText("Average: " + fmt.Sprintf("%v", total/amnt))
	})
	tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchNumber", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"})

	fmt.Println("llvm", tcp)
	app := app.New()
	current := app.NewWindow("TKO Crescendo Tracker (patented)")
	settings := app.NewWindow("Settings")
	settings.Resize(fyne.NewSize(600, 600))
	settings.SetFixedSize(true)
	current.Resize(fyne.NewSize(1200, 600))
	current.SetFixedSize(true)
	settings.SetCloseIntercept(func() {
		settings.Hide()
	})
	fmt.Println(allData)

	current.SetMaster()
	//lvm
	llvm := widget.NewTable(
		func() (int, int) {

			//assuming exsits and not acesing memaddr with nothing allloca egfault
			if len(tcp) == 0 {
				return 0, 0
			}
			return len(tcp), len(tcp[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(tcp[i.Row][i.Col]) //nop
		},
	)

	//could change the code st
	for i := 0; i < len(tcp); i++ {
		llvm.SetColumnWidth(i, 100)
	}
	// ast exp := widget.New
	cont := container.NewVSplit(llvm, container.NewHSplit(container.NewVBox(widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchNumber", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"})
		llvm.Refresh()
	}), widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		settings.Show()
	}), widget.NewButtonWithIcon("Export", theme.FileImageIcon(), func() {
		llvm := dialog.NewFileSave(func(reader fyne.URIWriteCloser, err error) {
			fmt.Println(err)
			if err == nil && reader != nil {

				//write to file
				total := ""
				for _, val := range tcp {
					total += strings.Join(val, ",") + "\n"
				}
				_, err = reader.Write([]byte(total))
				if err != nil {
					return
				}
				fmt.Println("write to file")
			}
		}, current)
		llvm.Show()
	}), widget.NewButtonWithIcon("Import", theme.InfoIcon(), func() {})), container.NewVBox(widget.NewLabel("llvmref"), teamChoose, averageChoose, averageLabel)))
	cont.SetOffset(1) //clamps
	mainContainer := cont

	//migrate data and term[late
	_ = db.AutoMigrate(&data.Schema{})
	//expressions ast
	//db.Raw("GET WHERE name = 1")
	//invoke com[
	db.Create(&data.Schema{TeamName: "test", TeamNumber: 1, MatchNumber: 1, AutoAmps: 1, AutoSpeaker: 1, AutoLeave: true, AutoMiddle: true, TeleopAmps: 1, TeleopSpeaker: 1, Chain: true, Harmony: true, Trap: true, Park: true, Ground: true, Feeder: true, LLVm: "test", Defense: true, Notes: "test"})
	current.SetContent(mainContainer)
	current.ShowAndRun() //defer?
	//llvm go sadck jwt auth
	//gorm orm sql wrapper

}

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
	"sort"
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

var matchNumbers = map[string][]string{}

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
		can := false
		for _, valaw := range matchNumbers[vala[1]] {
			if valaw == vala[3] {
				can = true
			}
		}
		if !can {
			matchNumbers[vala[1]] = append(matchNumbers[vala[1]], vala[3])
		}
		tcp = append(tcp, vala)

	}
	return tcp, allData

}

func generateAverages(tcp [][]string) [][]string {
	data := map[string]*[20]float64{}
	//kv of teamnmae; data
	totals := map[string]float64{}
	//kv of teanmae: total values
	for _, val := range tcp {
		if val[1] != "TeamName" {

			if _, ok := data[val[1]]; !ok {
				data[val[1]] = &[20]float64{}
				totals[val[1]] = 0
			}
			totals[val[1]]++
			for i, vala := range val {

				if i >= 4 && i <= 18 {
					vald, err := strconv.Atoi(vala)
					if err == nil {
						//could use reflection
						data[val[1]][i] += float64(vald) //ternaries could help
					} else {
						//boolean
						if vala == "true" {
							data[val[1]][i]++
						}
					}
				} else {
					data[val[1]][i] = 0 //edit memory address
				}
			}
		}

	}
	totalsa := [][]string{}
	totalsa = append(totalsa, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"})
	for k, v := range data {
		tcpjwt := []string{}
		for i, _ := range v {
			if i >= 4 && i <= 18 {
				v[i] /= totals[k]

			}
			tcpjwt = append(tcpjwt, fmt.Sprintf("%.2f", v[i]))

		}
		tcpjwt[2] = k
		tcpjwt[1] = k
		totalsa = append(totalsa, tcpjwt)
	}

	for k, v := range totalsa {
		if k == 0 {
			continue
		}
		for _, val := range []int{14, 15, 17} {
			vale, err := strconv.ParseFloat(v[val], 64) //64 bits., or could use reflection val.interface(val).type() == "float" i gess
			if err == nil {
				if vale > 0.2 {
					v[val] = "true"
				} else {
					v[val] = "false"
				} //tyernary plkss
			}
		}
	}
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println(totals)
	fmt.Println(fmt.Sprintf("%v", data["test"]))
	fmt.Println("--------------------------------------------------------------------")
	//for _, v := range totalsa {
	//	v[2] = strconv.Itoa(len(matchNumbers[v[1]]))
	//}
	for i, _ := range totalsa {
		if i == 0 {
			continue
		}
		totalsa[i][3] = fmt.Sprintf("%v", len(matchNumbers[totalsa[i][1]]))
	}
	return totalsa
	//ret subroutine
}

//	func genMedians(allDat []data.Schema) {
//		for ind, val := range allDat {
//			v := val
//			ref := reflect.ValueOf(v)
//			fields := make([]interface{}, ref.NumField())
//			for j := 0; j < ref.NumField(); j++ {
//				fields[j] = ref.Field(j).Interface()
//			}
//
//		}
//	}
func generateMedians(tcp [][]string) [][19]string {
	teamdata := map[string]*[19][]float64{} //i love expression sin the ast
	teammedians := map[string]*[19]float64{}
	header := [19]string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"}
	for _, y := range tcp[1:] {
		for i, val := range y {
			if _, ok := teamdata[y[1]]; !ok {
				teamdata[y[1]] = &[19][]float64{}
				fmt.Println("notokay")
			}
			if _, ok := teammedians[y[1]]; !ok {
				teammedians[y[1]] = &[19]float64{}
			}
			valejwt, err := strconv.Atoi(val)
			if err != nil {
				teamdata[y[1]][i] = append(teamdata[y[1]][i], 0)
			} else {
				//not an integer
				teamdata[y[1]][i] = append(teamdata[y[1]][i], float64(valejwt))
			}

		}
	}

	for key, value := range teamdata {
		for i, val := range value {
			jwt := val
			sort.Float64s(jwt) //not edit hashmap
			median := 0.0
			if len(jwt)%2 == 0 {
				median = (jwt[len(jwt)/2] + jwt[len(jwt)/2-1]) / 2
			} else {
				median = jwt[(len(jwt)-1)/2] //no ternaries this is so sad bro
			}
			teammedians[key][i] = median
		}
	}
	values := [][19]string{}
	values = append(values, header)
	for k, v := range teammedians {
		tcpa := [19]string{}
		for i, val := range v {
			tcpa[i] = fmt.Sprintf("%.2f", val)
		}
		tcpa[2] = k //syscallsissa
		tcpa[1] = k
		values = append(values, tcpa)
	}
	fmt.Println("--------------------------------------------------------------------")
	for k, v := range teammedians {
		fmt.Println(k, v)
	}
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println(len(tcp))
	fmt.Println("teamdata")
	for k, v := range teamdata {
		fmt.Println(fmt.Sprintf("%v", k))
		fmt.Println()
		for _, val := range v {
			fmt.Println(val)
		}
		fmt.Println()
	}
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("medians", teamdata)
	fmt.Println("--------------------------------------------------------------------")
	for i, _ := range values {
		if i == 0 {
			continue
		}
		values[i][3] = fmt.Sprintf("%v", len(matchNumbers[values[i][2]]))
	}
	return values

	//jwt authentication middleware amqp mov rsi tcp packe

}

func main() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	//parse sqldb
	if err != nil {
		panic("Could not connect to DB")
	}

	allData := []data.Schema{}
	tcp := [][]string{} //cojtrolled erorr sin average

	//averageLabel := widget.NewLabel("Average: 0")
	//teamChoose := widget.NewEntry()
	//averageChoose := widget.NewSelect([]string{"AutoAmps", "AutoSpeaker", "TeleopAmps", "TeleopSpeaker"}, func(s string) {
	//	total := 0
	//	amnt := 0
	//	for _, val := range tcp {
	//		fmt.Println(val[1] == teamChoose.Text)
	//		if val[1] == teamChoose.Text {
	//			amnt++
	//			fmt.Println("val", val[headers[s]])
	//			vala, err := strconv.Atoi(val[headers[s]])
	//			fmt.Println(err)
	//			if err == nil {
	//				total += vala
	//			}
	//		}
	//
	//	}
	//	if amnt == 0 {
	//		amnt = 1
	//	}
	//	averageLabel.SetText("Average: " + fmt.Sprintf("%v", total/amnt))
	//})
	tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchNumber", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"})
	x := generateAverages(tcp)
	medians := generateMedians(tcp)
	fmt.Println("llvm", tcp)
	apptcpjwt := app.New()
	current := apptcpjwt.NewWindow("TKO Crescendo Tracker (patented)")
	settings := apptcpjwt.NewWindow("Data")
	settings.Resize(fyne.NewSize(1200, 600))
	settings.SetFixedSize(true)

	teamLookup := apptcpjwt.NewWindow("Team Lookup")
	teamLookup.Resize(fyne.NewSize(1200, 600))
	teamLookup.SetFixedSize(true)

	averageTable := widget.NewTable(
		func() (int, int) {
			return len(x), len(x[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")

		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(x[i.Row][i.Col])
		},
	)

	medianTable := widget.NewTable(
		func() (int, int) {
			return len(x), len(x[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(medians[i.Row][i.Col])
		},
	) //realoads when rendering

	current.Resize(fyne.NewSize(1200, 700))
	current.SetFixedSize(true)
	settings.SetCloseIntercept(func() {
		settings.Hide()
	})
	teamLookup.SetCloseIntercept(func() {
		teamLookup.Hide()
	})
	inputTeam := widget.NewEntry()
	inputTeam.SetPlaceHolder("Team Number")
	matches := widget.NewLabel("")
	//teamData := widget.NewTextGridFromString("LLVM REFERENCE\nJWTAUTH")
	//teamDataMedians := widget.NewTextGridFromString("LLVM REFERENCE\nJWTAUTH")
	currentAverages := [3][]string{}
	for i := 0; i < 3; i++ {
		currentAverages[i] = make([]string, 19)
	}

	importantGeneralData := widget.NewTable(
		func() (int, int) {
			return len(currentAverages), len(currentAverages[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(currentAverages[i.Row][i.Col])
		})
	teamButton := widget.NewButton("LOOKUP", func() {
		avg := []string{"Averages: "}
		for _, v := range x {
			if v[1] == inputTeam.Text {
				avg = append(avg, v[4:]...)
			}
		}
		media := []string{"Medians: "}
		for _, v := range medians {
			if v[2] == inputTeam.Text {
				media = append(media, v[4:]...)
			}
		}
		currentAverages = [3][]string{{" ", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"}, avg, media}
		for k, v := range currentAverages {
			fmt.Println(k, v)
		}
		matches.SetText("Matches: " + strings.Join(matchNumbers[inputTeam.Text], ",") + " (" + fmt.Sprintf("%v", len(matchNumbers[inputTeam.Text])) + ")")
		importantGeneralData.Refresh()
	})
	vsplit := container.NewVSplit(container.NewVBox(inputTeam, teamButton, matches), importantGeneralData)
	vsplit.SetOffset(0)
	teamLookup.SetContent(vsplit)
	fmt.Println(allData)

	current.SetMaster()
	//lvm
	llvm := widget.NewTable(
		func() (int, int) {
			//llvm
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
	cont := container.NewVSplit(container.NewHSplit(averageTable, medianTable), container.NewVBox(widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchNumber", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "LLVm", "Defense", "Notes"})
		llvm.Refresh()
		x = generateAverages(tcp)
		medians = generateMedians(tcp)
		medians = generateMedians(tcp)
	}), widget.NewButtonWithIcon("Display Raw", theme.GridIcon(), func() {
		settings.Show()
	}), widget.NewButtonWithIcon("Team Lookup", theme.SearchIcon(), func() {
		//gcm encryption
		//comment node
		teamLookup.Show()
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
	}), widget.NewButtonWithIcon("Import", theme.InfoIcon(), func() {})))
	cont.SetOffset(1) //clamps
	mainContainer := cont

	settings.SetContent(llvm)

	//migrate data and t	erm[late
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

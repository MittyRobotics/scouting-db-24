package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wcharczuk/go-chart/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main/data"
	"os"
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
	"Mobility":      16,
	"Penalties":     17,
	"TechPenalties": 18,
	"GroundPickup":  19,
	"StartingPos":   20,
	"Defense":       21,
	"CenterRing":    22,
	"Notes":         23,
}

func trimTwo(val [][]string) [][]string {
	var jwt [][]string //semantics and usage changesd
	for _, vala := range val {
		jwt = append(jwt, vala[2:])
	}
	return jwt
}

//the key is the match number, the value is the team numbers in B{1..4} and R{1..4}
var matchIndex map[string][]string

func createNewElem(vala map[string]string) data.Schema {
	//lld reference
	//tcp packet stream

	//llvm reference
	schem := data.Schema{}
	schem.TeamName = vala["TEAMNUM"]
	schem.TeamNumber, _ = strconv.Atoi(vala["TEAMNUM"])
	schem.MatchNumber, _ = strconv.Atoi(vala["MATCHNUM"]) //value, reference
	schem.AutoAmps, _ = strconv.Atoi(vala["AUTONAMP"])
	schem.AutoSpeaker, _ = strconv.Atoi(vala["AUTONSPEAKER"])
	schem.AutoLeave, _ = strconv.ParseBool(vala["MOBILITY"])
	schem.AutoMiddle, _ = strconv.ParseBool(vala["CENTERRING"])
	schem.TeleopAmps, _ = strconv.Atoi(vala["TELEOPAMP"])
	schem.TeleopSpeaker, _ = strconv.Atoi(vala["TELEOPSPEAKER"])
	//schem.Chain, _ = strconv.ParseBool(vala["TEAMNUM"])
	//schem.Park, _ = strconv.ParseBool(vala["TEAMNUM"])
	switch vala["ENDGAME"] {
	case "PARK":
		schem.Park = true
		break
	case "CHAIN":
		schem.Chain = true
		break
	case "NONE":
		schem.Park = false
		schem.Chain = false
		break
	}
	schem.Harmony, _ = strconv.ParseBool(vala["HARMONY"])
	schem.Trap, _ = strconv.ParseBool(vala["TRAP"])

	schem.Ground, _ = strconv.ParseBool(vala["GROUNDPICKUP"])
	schem.Feeder, _ = strconv.ParseBool(vala["FEEDER"])
	schem.Mobility, _ = strconv.ParseBool(vala["MOBILITY"])
	schem.Penalties, _ = strconv.Atoi(vala["PENALTIES"])
	schem.TechPenalties, _ = strconv.Atoi(vala["TECHPENALTIES"])
	schem.GroundPickup, _ = strconv.ParseBool(vala["GROUNDPICKUP"])

	if vala["STARTINGPOS"] == "" {
		schem.StartingPos = 1
	} else {
		schem.StartingPos, _ = strconv.Atoi(string(rune((vala["STARTINGPOS"][0]))))
	}
	schem.Defense, _ = strconv.ParseBool(vala["DEFENDING"])
	schem.CenterRing, _ = strconv.ParseBool(vala["CENTERRING"])
	schem.Notes = vala["NOTES"]
	return schem

} //llvm go sdk

var matchNumbers = map[string][]string{}
var matchNumbersR = map[string][][]string{} //llvm reference comment node

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
		// llvm golang sdk reference sinc ei compile to ir
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
		if _, okay := matchNumbersR[vala[3]]; !okay {
			matchNumbersR[vala[3]] = [][]string{}
		}
		cana := false
		for _, valaw := range matchNumbersR[vala[3]] {
			if valaw[1] == vala[1] {
				cana = true
			}
		}
		//not needed since three distinct per match
		if !cana {
			matchNumbersR[vala[3]] = append(matchNumbersR[vala[3]], vala)
		}

		if !can {
			matchNumbers[vala[1]] = append(matchNumbers[vala[1]], vala[3])
		}
		tcp = append(tcp, vala)

	}
	return tcp, allData

}

func generateAverages(tcp [][]string) [][]string {
	data := map[string]*[24]float64{}
	//kv of teamnmae; data
	totals := map[string]float64{}
	//kv of teanmae: total values
	for _, val := range tcp {
		if val[1] != "TeamName" {

			if _, ok := data[val[1]]; !ok {
				data[val[1]] = &[24]float64{}
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
	totalsa = append(totalsa, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
	for k, v := range data {
		tcpjwt := []string{}
		for i, _ := range v {
			v[i] /= totals[k]
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
		for _, val := range []int{14, 15, 21, 16, 19, 22} {
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
func generateMedians(tcp [][]string) [][24]string {
	teamdata := map[string]*[24][]float64{} //i love expression sin the ast
	teammedians := map[string]*[24]float64{}
	header := [24]string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}
	for _, y := range tcp[1:] {
		for i, val := range y {
			if _, ok := teamdata[y[1]]; !ok {
				teamdata[y[1]] = &[24][]float64{}
				fmt.Println("notokay")
			}
			if _, ok := teammedians[y[1]]; !ok {
				teammedians[y[1]] = &[24]float64{}
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
	values := [][24]string{}
	values = append(values, header)
	for k, v := range teammedians {
		tcpa := [24]string{}
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
	//lld re
	graph := chart.
	BarChart{
		Title: "Autonomous Amps",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 40,
			},
		},
		YAxis: chart.YAxis{
			Name: "Amps Scored",
		},
		Height: 612,
		Width:  2048,
		Bars: []chart.Value{
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
			{Value: 5.0, Label: "llvm compiler"},
			{Value: 7.0, Label: "mmap llvm compiler"},
			{Value: 8.0, Label: "jwtllvm compiler"},
		},
	}
	bufToWrite, err := os.Create("graphs/graph.png")
	if err != nil {
		panic("Could not open file")
	}
	err = graph.Render(chart.PNG, bufToWrite)
	if err != nil {
		panic("could not render")
	}
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	err = db.AutoMigrate(&data.Schema{})
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
	os.Setenv("FYNE_THEME", "dark")
	os.Setenv("FYNE_FONT", "Ubuntu")
	tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
	x := generateAverages(tcp)
	medians := generateMedians(tcp)

	fmt.Println("llvm", tcp)
	apptcpjwt := app.New()

	//label := widget.NewLabel("Settings")
	label := widget.NewRichTextWithText("Control Panel")
	label.ParseMarkdown("# **Control Panel**")
	divider := widget.NewSeparator()

	current := apptcpjwt.NewWindow("TKO Crescendo Tracker (patented)")
	settings := apptcpjwt.NewWindow("Data")
	settings.Resize(fyne.NewSize(1200, 600))
	settings.SetFixedSize(true)
	//glsl reference

	teamLookup := apptcpjwt.NewWindow("Team Lookup")
	teamLookup.Resize(fyne.NewSize(1200, 600))
	teamLookup.SetFixedSize(true)

	matchScheduleShow := apptcpjwt.NewWindow("Match Schedule")
	matchScheduleShow.Resize(fyne.NewSize(1200, 600))
	matchScheduleShow.SetFixedSize(true)
	matchScheduleShow.SetCloseIntercept(func() {
		matchScheduleShow.Hide()
	})

	matchScheduleShowButton := widget.NewButtonWithIcon("Show Match Schedule", theme.RadioButtonCheckedIcon(), func() {
		matchScheduleShow.Show()
	})
	matchSchedule := [][]string{{"Number", "Blue 1", "Blue 2", "Blue 3", "Red 1", "Red 2", "Red 3"}}
	for i := 0; i < 77; i++ {
		matchSchedule = append(matchSchedule, make([]string, 7))
	}
	matchScheduleTable := widget.NewTable(
		func() (int, int) {
			return len(matchSchedule), len(matchSchedule[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(matchSchedule[i.Row][i.Col])
		},
	)
	matchScheduleImportButton := widget.NewButton("Import Match Schedule", func() {
		dia := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, matchScheduleShow)
				return
			}
			if reader.URI().Extension() != ".CSV" && reader.URI().Extension() != ".csv" {
				dialog.ShowError(errors.New(fmt.Sprintf("Invalid file type; expecting csv, not %v", reader.URI().Extension())), matchScheduleShow)
				return
			}
			defer reader.Close()
			buf := make([]byte, 1000000)
			_, err = reader.Read(buf)
			if err != nil {
				dialog.ShowError(err, matchScheduleShow)
				return
			}
			total := [][]string{}
			total = append(total, []string{"Number", "Blue 1", "Blue 2", "Blue 3", "Red 1", "Red 2", "Red 3"})
			lines := strings.Split(string(buf), "\n")
			fmt.Println("length: ", string(buf))
			for _, line := range lines {
				fmt.Println(line)
				values := strings.Split(line, ",")
				if len(values) < 7 {
					continue
				}
				matchIndex[values[1]] = values[2:]
				total = append(total, []string{values[1], values[2], values[3], values[4], values[5], values[6], values[7]})
			}
			matchSchedule = total
			matchScheduleTable.Refresh()

		}, matchScheduleShow)
		dia.Show()
	})

	vsplitMatchTables := container.NewVSplit(matchScheduleImportButton, matchScheduleTable)
	vsplitMatchTables.SetOffset(0)
	matchScheduleShow.SetContent(vsplitMatchTables)
	teamCharts := apptcpjwt.NewWindow("Team Charts")
	teamCharts.Resize(fyne.NewSize(1024, 1024))
	teamCharts.SetCloseIntercept(func() {
		teamCharts.Hide()
	})

	image := canvas.NewImageFromFile("graphs/graph.png")
	image.SetMinSize(fyne.NewSize(1024, 612))
	image.FillMode = canvas.ImageFillContain
	image.Resize(fyne.NewSize(1024, 612))
	imageSelect := widget.NewSelect([]string{"AutoAmps", "AutoSpeaker", "TeleopAmps", "TeleopSpeaker"}, func(s string) {
		image.File = "graphs/" + s + ".png"
		image.Refresh()
	})

	easterEgg := apptcpjwt.NewWindow("ITS A SECRET TO EVERYBODY")
	easterEgg.Resize(fyne.NewSize(513, 293))
	imagetwo := canvas.NewImageFromFile("0image.png")
	imagetwo.FillMode = canvas.ImageFillContain
	easterEgg.SetContent(imagetwo)
	easterEgg.SetCloseIntercept(func() {
		easterEgg.Hide()
	})
	button := widget.NewButton("", func() {
		easterEgg.Show()
	})
	vspli := container.NewVSplit(imageSelect, image)
	vspli.SetOffset(0)
	vsplitwo := container.NewVSplit(vspli, button)
	vsplitwo.SetOffset(1)
	teamCharts.SetContent(vsplitwo)

	matchLookup := apptcpjwt.NewWindow("Match Lookup")
	//lld
	matchLookup.Resize(fyne.NewSize(1200, 600))
	matchLookup.SetCloseIntercept(func() {
		matchLookup.Hide()
		//llvm reference
		//jwt reference
	})
	matchLookup.SetFixedSize(true)

	averageTable := widget.NewTable(
		func() (int, int) {
			return len(x), len(x[0]) - 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")

		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(x[i.Row][i.Col+2])
		},
	)

	medianTable := widget.NewTable(
		func() (int, int) {
			//commetnode reference llvm referecne
			return len(x), len(x[0]) - 2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(medians[i.Row][i.Col+2])
		},
	) //realoads when rendering

	avgssorted := widget.NewSelect([]string{"AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Penalties", "TechPenalties"}, func(s string) {
		index := headers[s]
		sort.Slice(x, func(illvm, jllvm int) bool {
			if x[illvm][1] == "TeamName" || x[jllvm][1] == "TeamName" {
				return false
			}
			val1, _ := strconv.ParseFloat(x[illvm][index], 64)
			val2, _ := strconv.ParseFloat(x[jllvm][index], 64)
			return val1 > val2
		})

		sort.Slice(medians, func(igcc, jgcc int) bool {
			if medians[igcc][1] == "TeamName" || medians[jgcc][1] == "TeamName" {
				return false
			}
			val1, _ := strconv.ParseFloat(medians[igcc][index], 64)
			val2, _ := strconv.ParseFloat(medians[jgcc][index], 64)
			fmt.Println(val2, " ", val1)
			return val1 > val2

		})
		//change valuie at mem addr

		averageTable.Refresh()

	})

	messorted := widget.NewSelect([]string{"AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Penalties", "TechPenalties"}, func(s string) {
		index := headers[s]
		sort.Slice(medians, func(igcc, jgcc int) bool {
			if medians[igcc][1] == "TeamName" || medians[jgcc][1] == "TeamName" {
				return false
			}
			val1, _ := strconv.ParseFloat(medians[igcc][index], 64)
			val2, _ := strconv.ParseFloat(medians[jgcc][index], 64)
			fmt.Println(val2, " ", val1)
			return val1 > val2

		})
		//change valuie at mem addr
		medianTable.Refresh()

	})

	averageWindow := apptcpjwt.NewWindow("Averages")
	medianWindow := apptcpjwt.NewWindow("Medians")
	averageWindow.Resize(fyne.NewSize(1200, 600))
	medianWindow.Resize(fyne.NewSize(1200, 600))
	averageVSplit := container.NewVSplit(avgssorted, averageTable)
	averageVSplit.SetOffset(0)
	medianVSplit := container.NewVSplit(messorted, medianTable)
	medianVSplit.SetOffset(0)
	averageWindow.SetContent(averageVSplit)
	medianWindow.SetContent(medianVSplit)
	averageWindow.SetCloseIntercept(func() {
		averageWindow.Hide()
	})
	medianWindow.SetCloseIntercept(func() {
		medianWindow.Hide()
	})

	current.Resize(fyne.NewSize(500, 300))
	current.SetFixedSize(true)
	settings.SetCloseIntercept(func() {
		settings.Hide()
	})
	teamLookup.SetCloseIntercept(func() {
		teamLookup.Hide()
	})
	inputTeam := widget.NewEntry()
	inputMatch := widget.NewEntry()
	inputTeam.SetPlaceHolder("Team Number")
	inputMatch.SetPlaceHolder("Match Number")
	matches := widget.NewLabel("")

	//teamData := widget.NewTextGridFromString("LLVM REFERENCE\nJWTAUTH")
	//teamDataMedians := widget.NewTextGridFromString("LLVM REFERENCE\nJWTAUTH")
	currentAverages := [3][]string{}
	matchDatas := [19][]string{}
	for i := 0; i < 19; i++ {
		matchDatas[i] = make([]string, 24)
	}
	for i := 0; i < 3; i++ {
		currentAverages[i] = make([]string, 24)
	}

	importantMatchData := widget.NewTable(
		func() (int, int) {
			return len(matchDatas), len(matchDatas[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(matchDatas[i.Row][i.Col])
		},
	)
	importantGeneralData := widget.NewTable(
		func() (int, int) {
			return len(currentAverages), len(currentAverages[0])
		},
		//lld reference
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(currentAverages[i.Row][i.Col])
		})

	teamButton := widget.NewButton("LOOKUP", func() {
		teamCharts.Hide()
		if inputTeam.Text == "" {
			return
		}
		avg := []string{"Averages: "}
		for _, v := range x {
			if v[1] == inputTeam.Text {
				avg = append(avg, v[4:]...)
			}
		}
		if len(avg) == 1 {
			dialog.ShowError(errors.New(fmt.Sprintf("Team '%v' not found", inputTeam.Text)), teamLookup)
			return
		}
		media := []string{"Medians: "}
		for _, v := range medians {
			if v[2] == inputTeam.Text {
				media = append(media, v[4:]...)
			}
		}
		datas := map[string][][]int{
			"AutoAmps":      {},
			"AutoSpeaker":   {},
			"TeleopAmps":    {},
			"TeleopSpeaker": {},
		}
		for _, v := range allData {
			if v.TeamName != inputTeam.Text || strconv.Itoa(v.TeamNumber) != inputTeam.Text {
				continue
			}
			datas["AutoAmps"] = append(datas["AutoAmps"], []int{v.AutoAmps, v.MatchNumber})
			datas["AutoSpeaker"] = append(datas["AutoSpeaker"], []int{v.AutoSpeaker, v.MatchNumber})
			datas["TeleopAmps"] = append(datas["TeleopAmps"], []int{v.TeleopAmps, v.MatchNumber})
			datas["TeleopSpeaker"] = append(datas["TeleopSpeaker"], []int{v.TeleopSpeaker, v.MatchNumber})
			//could use reflection
		}
		fmt.Println(datas)
		for k, v := range datas {
			sort.Slice(v, func(i, j int) bool {
				return v[i][1] < v[j][1]
			})
			values := []chart.Value{}
			for _, val := range v {
				values = append(values, chart.Value{Value: float64(val[0]), Label: fmt.Sprintf("Match %v (Score: %v)", val[1], val[0])})
			}

			graph := chart.
			BarChart{
				Title: k,
				Background: chart.Style{
					Padding: chart.Box{
						Top:    40,
						Bottom: 40,
					},
				},
				YAxis: chart.YAxis{
					Name:  k,
					Range: &chart.ContinuousRange{Min: 0, Max: 10},
				},
				Height: 612,
				Width:  1024,
				Bars:   values,
			}
			bufToWrite, _ := os.Create(fmt.Sprintf("graphs/%v.png", k))
			_ = graph.Render(chart.PNG, bufToWrite)

		}

		//llreference

		currentAverages = [3][]string{{"", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}, avg, media}
		matchDatas[0] = []string{"", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}
		for k, v := range currentAverages {
			fmt.Println(k, v)
		}
		matches.SetText("Matches: " + strings.Join(matchNumbers[inputTeam.Text], ",") + " (" + fmt.Sprintf("%v", len(matchNumbers[inputTeam.Text])) + ")")
		fmt.Println(matchNumbersR)
		importantGeneralData.Refresh()
	})
	matchButton := widget.NewButton("LOOKUP", func() {
		fmt.Println(matchNumbersR)
		if inputMatch.Text == "" {
			return
		}
		vals, ok := matchNumbersR[inputMatch.Text]
		if !ok {
			return
		}
		kTeamNumberVData := map[string][][]string{}
		for ind, teamName := range vals {
			avg := []string{"Averages: "}
			for _, v := range x {
				if v[1] == teamName[1] {
					avg = append(avg, v[1:]...)
				}
			}
			media := []string{"Medians: "}
			for _, v := range medians {
				if v[2] == teamName[1] {
					media = append(media, v[1:]...)
				}
				//comment node
			}
			avg[4] = inputMatch.Text
			media[4] = inputMatch.Text
			kTeamNumberVData[teamName[1]] = [][]string{avg, media}
			tcpa := []string{"Performance: "}
			tcpa = append(tcpa, teamName...)
			kTeamNumberVData[teamName[1]] = append(kTeamNumberVData[teamName[1]], tcpa)

			matchDatas[(ind*3)+1] = avg
			matchDatas[ind*3+2] = media
			matchDatas[0] = []string{"ID", "TeamName", "TeamNumber", "Match", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}
			for k, v := range matchDatas {
				fmt.Println(k, v)
			}

		}
		//.env js libray me omw to make serverside js unordered hashmap
		val := 1
		for _, v := range kTeamNumberVData {
			for _, vala := range v {
				matchDatas[val] = vala
				val++
			}
		}
		importantMatchData.Refresh()
		//
	})

	teamChart := canvas.NewImageFromFile("graphs/graph.png")
	renderChartButton := widget.NewButtonWithIcon("Render Charts", theme.DocumentIcon(), func() {
		teamCharts.Show()
	})
	teamChart.Resize(fyne.NewSize(2048, 612))
	teamChart.FillMode = canvas.ImageFillOriginal
	vsplit := container.NewVSplit(container.NewVBox(inputTeam, teamButton, renderChartButton, matches), importantGeneralData)
	secondvsplit := container.NewVSplit(container.NewVBox(inputMatch, matchButton, matches), importantMatchData)
	vsplit.SetOffset(0)
	secondvsplit.SetOffset(0)
	matchLookup.SetContent(secondvsplit)
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
			return len(tcp), len(tcp[0]) //llvm
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
	cont := container.NewVBox(label, divider, widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
		llvm.Refresh()
		x = generateAverages(tcp)
		medians = generateMedians(tcp)
		medians = generateMedians(tcp)
		medianTable.Refresh()
		averageTable.Refresh()
	}), matchScheduleShowButton, widget.NewButtonWithIcon("Display Raw", theme.GridIcon(), func() {
		settings.Show()
	}), widget.NewButtonWithIcon("Display Averages", theme.RadioButtonIcon(), func() {
		averageWindow.Show()
	}), widget.NewButtonWithIcon("Display Medians", theme.RadioButtonIcon(), func() {
		medianWindow.Show()
	}),
		widget.NewButtonWithIcon("Match Lookup", theme.SearchIcon(), func() {
			matchLookup.Show()
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
			//llvm
			llvm.Show()
		}), widget.NewButtonWithIcon("Import", theme.InfoIcon(), func() {
			//mut value allocated
			jwtauth := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
				if err == nil && reader != nil {

					fmt.Println(reader.Path())
					//aloloc to memory address
					alert := dialog.NewConfirm("Import?", fmt.Sprintf("Are you sure you want to import from directory '%v'?", reader.Name()), func(yes bool) {
						//lld
						if yes {
							dir, _ := os.Open(reader.Path())
							filevec, err := dir.Readdir(0)
							if err != nil {
								return
							}
							for _, val := range filevec {
								file, err := os.Open(reader.Path() + "/" + val.Name())
								if err != nil {
									return
								}
								buffer := make([]byte, val.Size())
								toAdd := map[string]string{}
								_, _ = file.Read(buffer)
								split := strings.Split(string(buffer), "\n")
								for _, v := range split {
									splittwo := strings.Split(v, ":")
									if len(splittwo) == 2 {
										toAdd[splittwo[0]] = strings.Trim(splittwo[1], " ")
									} else {
										toAdd[splittwo[0]] = ""
									}
								}
								vale := createNewElem(toAdd)
								db.Create(&vale)
								//lld linker scripts mov rsi, rdi
								tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
								llvm.Refresh()
								x = generateAverages(tcp)
								medians = generateMedians(tcp)
								medians = generateMedians(tcp)
								medianTable.Refresh()
								averageTable.Refresh()

							}

							//a random valid that actually exsists google oauth key scraped from the internet is:  "client_id": "32555940559.apps.googleusercontent.com",
						}
					}, current)
					alert.Show()

				}
				fmt.Println(err)
			}, current)
			jwtauth.Show()
		}))
	//clamps

	settings.SetContent(llvm)
	//comemt node
	//ast

	//jwt reference
	//tcp reference

	//open syscall
	//migrate data and t	erm[late
	//expressions ast
	//db.Raw("GET WHERE name = 1")
	//invoke com[
	//db.Create(&data.Schema{
	//	TeamName:      "llvm",
	//	TeamNumber:    0,
	//llvm comment
	//lvm
	//	MatchNumber:   0,
	//	AutoAmps:      2,
	//	AutoSpeaker:   2,
	//	AutoLeave:     false,
	//	AutoMiddle:    false,
	//	TeleopAmps:    2,
	//	TeleopSpeaker: 2,
	//	Chain:         false,
	//	Harmony:       false,
	//	Trap:          false,
	//	Park:          false,
	//	Ground:        false,
	//	Feeder:        false,
	//	Mobility:      false,
	//	Penalties:     2,
	//	TechPenalties: 2,
	//	GroundPickup:  false,
	//	StartingPos:   2,
	//	Defense:       false,
	//	CenterRing:    false,
	//	Notes:         "",
	//})

	current.SetContent(cont)
	current.ShowAndRun() //defer?
	//llvm go sadck jwt auth
	//lld linker go sdk
	//gorm orm sql wrapper

}

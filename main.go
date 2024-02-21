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
	os.Setenv("FYNE_THEME", "light")
	os.Setenv("FYNE_FONT", "Ubuntu")
	tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
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

	matchLookup := apptcpjwt.NewWindow("Match Lookup")
	matchLookup.Resize(fyne.NewSize(1200, 600))
	matchLookup.SetCloseIntercept(func() {
		matchLookup.Hide()
	})
	matchLookup.SetFixedSize(true)

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
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(currentAverages[i.Row][i.Col])
		})

	teamButton := widget.NewButton("LOOKUP", func() {
		if inputTeam.Text == "" {
			return
		}
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

		currentAverages = [3][]string{{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}, avg, media}
		matchDatas[0] = []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}
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
					avg = append(avg, v[:]...)
				}
			}
			media := []string{"Medians: "}
			for _, v := range medians {
				if v[2] == teamName[1] {
					media = append(media, v[:]...)
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
			matchDatas[0] = []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"}
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
	vsplit := container.NewVSplit(container.NewVBox(inputTeam, teamButton, matches), importantGeneralData)
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
	cont := container.NewVSplit(container.NewHSplit(averageTable, medianTable), container.NewVBox(widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		tcp, allData = populate(db, allData, tcp, []string{"ID", "TeamName", "TeamNumber", "MatchesPlayed", "AutoAmps", "AutoSpeaker", "AutoLeave", "AutoMiddle", "TeleopAmps", "TeleopSpeaker", "Chain", "Harmony", "Trap", "Park", "Ground", "Feeder", "Mobility", "Penalities", "Tech-Pens", "Ground-Pick", "Starting-Pos", "Defense", "CenterRing", "Notes"})
		llvm.Refresh()
		x = generateAverages(tcp)
		medians = generateMedians(tcp)
		medians = generateMedians(tcp)
	}), widget.NewButtonWithIcon("Display Raw", theme.GridIcon(), func() {
		settings.Show()
	}), widget.NewButtonWithIcon("Match Lookup", theme.SearchIcon(), func() {
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
							//toAdd := data.Schema{}
							//value := reflect.ValueOf(&toAdd).Elem()
							//vals := map[string][]int{ //ENDGAME SPECIAL CASE
							//	"TEAMNUM":       {1, 2},
							//	"MATCHNUM":      {3},
							//	"MOBILITY":      {16},
							//	"DEFENDING":     {21},
							//	"STARTINGPOS":   {20},
							//	"AUTONSPEAKER":  {5},
							//	"AUTONAMP":      {4},
							//	"CENTERRING":    {22},
							//	"TELEOPSPEAKER": {9},
							//	"TELEOPAMP":     {8},
							//	"TRAP":          {12},
							//	"HARMONY":       {11},
							//	"GROUND":        {14},
							//	"FEEDER":        {15},
							//	"PENALTIES":     {17},
							//	"TECHPENALTIES": {18},
							//	"NOTES":         {23},
							//}
							buffer := make([]byte, val.Size())
							_, _ = file.Read(buffer)
							fmt.Println(string(buffer))
						}

						//a random valid that actually exsists google oauth key scraped from the internet is:  "client_id": "32555940559.apps.googleusercontent.com",
					}
				}, current)
				alert.Show()

			}
			fmt.Println(err)
		}, current)
		jwtauth.Show()
	})))
	cont.SetOffset(1) //clamps
	mainContainer := cont

	settings.SetContent(llvm)

	//migrate data and t	erm[late
	//expressions ast
	//db.Raw("GET WHERE name = 1")
	//invoke com[
	//db.Create(&data.Schema{
	//	TeamName:      "llvm",
	//	TeamNumber:    0,
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
	//db.Create(&data.Schema{
	//	TeamName:      "llvm",
	//	TeamNumber:    0,
	//	MatchNumber:   0,
	//	AutoAmps:      4,
	//	AutoSpeaker:   4,
	//	AutoLeave:     false,
	//	AutoMiddle:    false,
	//	TeleopAmps:    4,
	//	TeleopSpeaker: 4,
	//	Chain:         false,
	//	Harmony:       false,
	//	Trap:          false,
	//	Park:          false,
	//	Ground:        false,
	//	Feeder:        false,
	//	Mobility:      false,
	//	Penalties:     4,
	//	TechPenalties: 4,
	//	GroundPickup:  false,
	//	StartingPos:   4,
	//	Defense:       false,
	//	CenterRing:    false,
	//	Notes:         "",
	//})
	current.SetContent(mainContainer)
	current.ShowAndRun() //defer?
	//llvm go sadck jwt auth
	//gorm orm sql wrapper

}

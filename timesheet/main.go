package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type user struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type project struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type timeEntry struct {
	SpentDate string   `json:"spent_date"`
	User      *user    `json:"user"`
	Project   *project `json:"project"`
	Hours     float64  `json:"hours"`
	Notes     string   `json:"notes"`
}

type timesheet struct {
	TimeEntries []timeEntry `json:"time_entries"`
}

type proj struct {
	ID   int32  `json:"id"`
	Code string `json:"code"`
}

type projects struct {
	Projects []proj `json:"projects"`
}

type row struct {
	SpentDate string
	User      string
	Project   string
	Hours     string
	Notes     string
}

type config struct {
	From          *string `json:"from"`
	To            *string `json:"to"`
	AccountID     string  `json:"accountID"`
	Authorization string  `json:"authorization"`
	UA            string  `json:"ua"`
	URLs          *urls   `json:"urls"`
	User          string  `json:"spUser"`
	DateInput     string  `json:"dateInput"`
	DateOutput    string  `json:"dateOutput"`
	Header        string  `json:"header"`
}

type urls struct {
	TimeEntries string `json:"timeEntries"`
	Projects    string `json:"projects"`
}

func main() {
	var excel = flag.String("excel", "", "excel filename to import")

	from := flag.String("from", "", "start date (YYYY-MM-DD) to import (default: last week's Monday)")
	to := flag.String("to", "", "end date (YYYY-MM-DD) to import (default: last week's Friday)")

	flag.Parse()

	cfg := readConfig(from, to)

	if *excel != "" {
		readFromExcel(*excel)
		return
	}

	//fmt.Println(cfg, *cfg.From, *cfg.To)
	readFromHarvest(cfg)
}

func readConfig(from, to *string) config {
	var cfg = config{}
	b, err := ioutil.ReadFile("./.config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t := time.Now()
	if *from == "monthly" {
		md := t.Day()
		//special case, build monthly timesheet for previous month
		*from = time.Now().AddDate(0, -1, -1*md+1).Format(cfg.DateInput)
		*to = time.Now().AddDate(0, 0, -1*md).Format(cfg.DateInput)
	} else if *from == "" {
		wd := int(t.Weekday())
		//last week's Monday
		*from = time.Now().AddDate(0, 0, -1*wd-6).Format(cfg.DateInput)
		if *to == "" {
			//last week's Friday
			*to = time.Now().AddDate(0, 0, -1*wd-2).Format(cfg.DateInput)
		}
	}

	//if only to was blank, then use previous day
	if *to == "" {
		*to = time.Now().AddDate(0, 0, -1).Format(cfg.DateInput)
	}

	*cfg.From = *from
	*cfg.To = *to

	return cfg
}

func readFromHarvest(cfg config) {
	ts, err := getTimesheet(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r, err := processTimesheet(ts, cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outputTimesheet(r, cfg)
}

func getTimesheet(cfg config) (*timesheet, error) {
	resp, err := requestHarvestTimeEntries(cfg)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ts, err := readTimesheetFromResponse(body)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func requestHarvestTimeEntries(cfg config) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v?from=%v&to=%v", cfg.URLs.TimeEntries, *cfg.From, *cfg.To), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Harvest-Account-ID", cfg.AccountID)
	req.Header.Add("Authorization", cfg.Authorization)
	req.Header.Add("User-Agent", cfg.UA)
	return client.Do(req)
}

func readTimesheetFromResponse(b []byte) (*timesheet, error) {
	var t *timesheet
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func outputTimesheet(r []row, cfg config) {
	fmt.Println(cfg.Header)
	for _, row := range r {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", row.SpentDate, row.Project, row.Notes, row.Hours, row.User)
	}
}

func getProjects(cfg config) (*projects, error) {
	resp, err := requestHarvestProjects(cfg)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	p, err := readProjectsFromResponse(body)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func requestHarvestProjects(cfg config) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", cfg.URLs.Projects, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Harvest-Account-ID", cfg.AccountID)
	req.Header.Add("Authorization", cfg.Authorization)
	req.Header.Add("User-Agent", cfg.UA)
	return client.Do(req)
}

func readProjectsFromResponse(b []byte) (*projects, error) {
	var p *projects
	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getProjectCode(p *projects, id int32) string {
	for _, v := range p.Projects {
		if v.ID == id {
			return v.Code
		}
	}
	return ""
}

func processTimesheet(t *timesheet, cfg config) ([]row, error) {
	var rr []row
	projects, err := getProjects(cfg)
	if err != nil {
		return nil, err
	}
	for i := len(t.TimeEntries) - 1; i >= 0; i-- {
		e := t.TimeEntries[i]
		sd, _ := time.Parse(cfg.DateInput, e.SpentDate)
		r := row{
			SpentDate: sd.Format(cfg.DateOutput),
			User:      cfg.User,
			Project:   getProjectCode(projects, e.Project.ID),
			Hours:     strconv.FormatFloat(e.Hours, 'f', 2, 64),
			Notes:     e.Notes,
		}
		rr = append(rr, r)
	}
	return rr, nil
}

func readFromExcel(filename string) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Date\tProject\tNotes\tHours\tPerson")
	for _, sheet := range xlFile.Sheets {
		for rx, row := range sheet.Rows {
			if rx > 0 {
				date := strings.Replace(row.Cells[0].String(), "-", "/", 3)
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n", date, row.Cells[3], row.Cells[5], row.Cells[6], row.Cells[11])
			}
		}
	}
}

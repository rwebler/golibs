//Package harvest gets data from
//the Harvest API
package harvest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"timesheet/config"
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

//Read converts time entries from
//Harvest to Tempo
func Read(cfg config.Config) {
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

func getTimesheet(cfg config.Config) (*timesheet, error) {
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

func requestHarvestTimeEntries(cfg config.Config) (*http.Response, error) {
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

func outputTimesheet(r []config.Row, cfg config.Config) {
	fmt.Println(cfg.Header)
	for _, row := range r {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", row.SpentDate, row.Project, row.Notes, row.Hours, row.User)
	}
}

func getProjects(cfg config.Config) (*projects, error) {
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

func requestHarvestProjects(cfg config.Config) (*http.Response, error) {
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

func processTimesheet(t *timesheet, cfg config.Config) ([]config.Row, error) {
	var rr []config.Row
	projects, err := getProjects(cfg)
	if err != nil {
		return nil, err
	}
	for i := len(t.TimeEntries) - 1; i >= 0; i-- {
		e := t.TimeEntries[i]
		sd, _ := time.Parse(cfg.DateInput, e.SpentDate)
		r := config.Row{
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

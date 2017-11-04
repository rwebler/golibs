// Package config contains basic structures
// needed for the timesheet import to work
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Row represents the final state
// of a row that will be imported into
// Tempo.
type Row struct {
	SpentDate string
	User      string
	Project   string
	Hours     string
	Notes     string
}

// Config contains the values
// used throughout timesheet
type Config struct {
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

//Read gets values from .config.json
//into the the Config struct
func Read(from, to *string) Config {
	var cfg = Config{}
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

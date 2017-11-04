// Package excel has a function
// to read Harvest Excel sheets
package excel

import (
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

//Read converts a Harvest Excel timesheet
// into a Tempo timesheet
func Read(filename string) {
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

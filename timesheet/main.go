package main

import (
	"flag"

	"timesheet/config"
	"timesheet/excel"
	"timesheet/harvest"
)

func main() {
	var ef = flag.String("excel", "", "excel filename to import")

	from := flag.String("from", "", "start date (YYYY-MM-DD) to import (default: last week's Monday)")
	to := flag.String("to", "", "end date (YYYY-MM-DD) to import (default: last week's Friday)")

	flag.Parse()

	if *ef != "" {
		excel.Read(*ef)
		return
	}

	cfg := config.Read(from, to)

	harvest.Read(cfg)
}

# Timesheet

## Converts from Harvest spreadsheets and time entries into Tempo timesheets

### Install

1. Get the code
2. Copy `.config.json.default` into `.config.json`. In this file:
  1. Fill in your Harvest account ID and authorization code on the eponymous fields
  2. Fill in your Tempo Person data on the `spUser` field
  3. Provide a better description inside the `ua` field
3. Build the code

### Usage

From `./timesheet -h`:

```
Usage of ./timesheet:
  -excel string
    	excel filename to import
  -from string
    	start date (YYYY-MM-DD) to import (default: last week's Monday)
  -to string
    	end date (YYYY-MM-DD) to import (default: last week's Friday)
```

`./timesheet -from=monthly` will export a timesheet for the previous month

### Notes

* I pipe the output of timesheet to `pbcopy` on macOS, so that I'll just need to paste it into Tempo. Since we're able to review the timesheet in Tempo before submitting, it's safe to do so.

* This config works for me because I use Harvest's `code` attribute in a Project for SP's project number (using the `[XXX] Name` format).
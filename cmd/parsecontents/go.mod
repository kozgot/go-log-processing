module github.com/kozgot/go-log-processing/cmd/parsecontents

go 1.15
require (
	github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0
	github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0
)

replace github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0 => ../filterlines
replace github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0 => ../parsedates
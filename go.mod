module github.com/kozgot/go-log-processing

go 1.15

require (
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/dustin/go-humanize v1.0.0
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/kozgot/go-log-processing/cmd/elasticsearch v0.0.0
	github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0
	github.com/kozgot/go-log-processing/cmd/parsecontents v0.0.0
	github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0
)

replace github.com/kozgot/go-log-processing/cmd/elasticsearch v0.0.0 => ./cmd/elasticsearch
replace github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0 => ./cmd/filterlines
replace github.com/kozgot/go-log-processing/cmd/parsecontents v0.0.0 => ./cmd/parsecontents
replace github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0 => ./cmd/parsedates

module github.com/kozgot/go-log-processing

go 1.15

require (
	github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0
	github.com/kozgot/go-log-processing/cmd/parsecontents v0.0.0
	github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0
	github.com/streadway/amqp v1.0.0
)

replace github.com/kozgot/go-log-processing/cmd/filterlines v0.0.0 => ./cmd/filterlines

replace github.com/kozgot/go-log-processing/cmd/parsecontents v0.0.0 => ./cmd/parsecontents

replace github.com/kozgot/go-log-processing/cmd/parsedates v0.0.0 => ./cmd/parsedates

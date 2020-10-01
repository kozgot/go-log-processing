# Diagnostic log processing and analysis

This repository contains the source code of my independent laboratory project, which is Diagnostic log data processing

## System architecture
![alt text](https://github.com/kozgot/go-log-processing/blob/master/images/abra.PNG)

The system is going to follow an event-driven microservice architecture.

### Elasticsearch & Kibana
Elasticsearch and Kibana are part of the “ELK stack.” It is a software stack frequently used for data storage and visualization in business intelligence systems. This software stack was chosen for the storage, analysis and visualization of the data.

### Results uploader
This component will be responsible for uploading the parsed input to Elasticsearch.

### Azure blob storage
This component will provide a communication platform for the other components. Each component is going to retrieve its input, and write its output here.

### DLMS file parser
This component will be responsible for the parsing of DLMS messages. It is going to written in Python to make parsing the binary format easier. All other components are going to be written in Go.

### Postprocessing
This component will be responsible for the additional processing of imported files and messages (eg. filtering duplications).

### Log import
This component will be responsible for importing data from the textual log files. It migth upload the imported data directly to Elasticsearch. For now, implementation will focus mainly on this component. 

## Golang
Most of the components are going to be implemented using the Go language. I chose this language for the following resons.

Go makes it easier (than for example Java or Python) to write correct, clear and efficient code. 
 - The core language consists of a few simple, orthogonal features that can be combined in a relatively small number of ways. This makes it easier to learn the language, and to read and write programs.
 - Go is strongly and statically typed with no implicit conversions, but the syntactic overhead is still surprisingly small. 
 - Programs are constructed from packages that offer clear code separation and allow efficient management of dependencies. 
 - Structurally typed interfaces provide runtime polymorphism through dynamic dispatch.
 - Concurrency is an integral part of Go, supported by goroutines, channels and the select statement.
 
 ## Project layout
 The structure of the project follows some guidelines from https://github.com/golang-standards/project-layout .

## Development
For development, I use the VSCode IDE with the Go extension.

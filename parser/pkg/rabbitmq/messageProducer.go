package rabbitmq

import "github.com/kozgot/go-log-processing/parser/pkg/models"

// MessageProducer encapsulates methods used to communicate with rabbitMQ server.
type MessageProducer interface {
	PublishStringMessage(indexName string)
	PublishEntry(line models.ParsedLogEntry)
	OpenChannelAndConnection()
	CloseChannelAndConnection()
}

package rabbitmq

import "github.com/kozgot/go-log-processing/postprocessor/pkg/models"

// MessageProducer encapsulates methods used to publish data for ES uploader service.
type MessageProducer interface {
	PublishEvent(event models.SmcEvent)
	PublishConsumption(cons models.ConsumtionValue)
	Connect()
	CloseChannelAndConnection()
}

package processing

import (
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// ConsumptionProcessor encapsulates consumption processsing data and logic.
type ConsumptionProcessor struct {
	consumptionValues    []models.ConsumtionValue
	indexValues          []models.IndexValue
	messageProducer      rabbitmq.MessageProducer
	consumptionIndexName string
}

// NewConsumptionProcessor creates a new consmptionprocessor instance.
func NewConsumptionProcessor(
	consumptionValues []models.ConsumtionValue,
	indexValues []models.IndexValue,
	messageProducer rabbitmq.MessageProducer,
	consumptionIndexName string) *ConsumptionProcessor {
	consumptionProcessor := ConsumptionProcessor{
		consumptionValues:    consumptionValues,
		indexValues:          indexValues,
		messageProducer:      messageProducer,
		consumptionIndexName: consumptionIndexName,
	}

	return &consumptionProcessor
}

// ProcessConsumptionAndIndexValues performs further processing on to retrieve consumption values for SMCs.
func (consumptionProcessor *ConsumptionProcessor) ProcessConsumptionAndIndexValues() {
	for _, cons := range consumptionProcessor.consumptionValues {
		smcUID := consumptionProcessor.findRelatedSmc(cons)
		if smcUID != "" {
			cons.SmcUID = smcUID
			consumptionProcessor.messageProducer.PublishConsumption(cons, consumptionProcessor.consumptionIndexName)
		}
	}
}

func (consumptionProcessor *ConsumptionProcessor) findRelatedSmc(cons models.ConsumtionValue) string {
	for _, index := range consumptionProcessor.indexValues {
		if index.ServiceLevel == cons.ServiceLevel &&
			index.ReceiveTime == cons.ReceiveTime &&
			index.PreviousTime == cons.StartTime {
			return index.SmcUID
		}
	}
	return ""
}

package processing

import (
	"github.com/kozgot/go-log-processing/postprocessor/internal/rabbitmq"
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
	"github.com/streadway/amqp"
)

func ProcessConsumptionAndIndexValues(
	consumptionValues []models.ConsumtionValue,
	indexValues []models.IndexValue,
	channel *amqp.Channel) {
	consumptionsBySmcUID := make(map[string][]models.ConsumtionValue)
	for _, cons := range consumptionValues {
		smcUID := findRelatedSmc(cons, indexValues)
		if smcUID != "" {
			initConsumptionArrayIfNeeded(consumptionsBySmcUID, smcUID)
			cons.SmcUID = smcUID
			consumptionsBySmcUID[smcUID] = append(consumptionsBySmcUID[smcUID], cons)
			saveConsumptionToDB(cons, channel)
		}
	}
}

func saveConsumptionToDB(cons models.ConsumtionValue, channel *amqp.Channel) {
	rabbitmq.SendConsumptionToElasticUploader(cons, channel, "consumption")
}

func initConsumptionArrayIfNeeded(consumptionsBySmcUID map[string][]models.ConsumtionValue, uid string) {
	_, ok := consumptionsBySmcUID[uid]
	if !ok {
		consumptionsBySmcUID[uid] = []models.ConsumtionValue{}
	}
}

func findRelatedSmc(cons models.ConsumtionValue, indexValues []models.IndexValue) string {
	for _, index := range indexValues {
		if index.ServiceLevel == cons.ServiceLevel &&
			index.ReceiveTime == cons.ReceiveTime &&
			index.PreviousTime == cons.StartTime {
			return index.SmcUID
		}
	}
	return ""
}

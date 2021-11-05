package processing

import (
	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

func (processor *EntryProcessor) ProcessConsumptionAndIndexValues() {
	consumptionsBySmcUID := make(map[string][]models.ConsumtionValue)
	for _, cons := range processor.consumptionValues {
		smcUID := findRelatedSmc(cons, processor.indexValues)
		if smcUID != "" {
			initConsumptionArrayIfNeeded(consumptionsBySmcUID, smcUID)
			cons.SmcUID = smcUID
			consumptionsBySmcUID[smcUID] = append(consumptionsBySmcUID[smcUID], cons)
			processor.esUploader.SendConsumptionToElasticUploader(cons)
		}
	}
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

package testutils

import (
	"reflect"
	"testing"

	"github.com/kozgot/go-log-processing/postprocessor/pkg/models"
)

// AssertEqualSmcData fails the test if the given smc datas are not equal.
func AssertEqualSmcData(actual *models.SmcData, expected *models.SmcData, t *testing.T, testIndex int) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected SMC Data does not match actual SMC Data in test no. %d.", testIndex+1)
	}
}

// AssertEqualSmcEvent fails the test if the given smc events are not equal.
func AssertEqualSmcEvent(actual *models.SmcEvent, expected *models.SmcEvent, t *testing.T, testIndex int) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected SMC Event does not match actual SMC Event in test no. %d.", testIndex+1)
	}
}

// AssertEqualConsumption fails the test if the given consumption values are not equal.
func AssertEqualConsumption(
	actual *models.ConsumtionValue,
	expected *models.ConsumtionValue,
	t *testing.T,
	testIndex int) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected Consumtion Value does not match actual Consumtion Value in test no. %d.", testIndex+1)
	}
}

// AssertEqualIndex fails the test if the given index values are not equal.
func AssertEqualIndex(actual *models.IndexValue, expected *models.IndexValue, t *testing.T, testIndex int) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected Index Value does not match actual Index Value in test no. %d.", testIndex+1)
	}
}

package serviceintegrationtests

import (
	"testing"

	"github.com/kozgot/go-log-processing/elasticuploader/tests/testutils"
)

func TestServiceIntegrationWithElasticsearch(t *testing.T) {
	testesclient := testutils.NewTestEsClientWrapper()
	testesclient.QueryIndex("smc")
}

func TestServiceIntegrationWithRabbitMQ(t *testing.T) {

}

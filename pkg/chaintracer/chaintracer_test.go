package chaintracer_test

import (
	"reflect"
	"testing"

	"github.com/primevprotocol/oracle/pkg/chaintracer"
)

func TestDataPull(t *testing.T) {
	tracer := chaintracer.NewIncrementingTracer(18293308)
	_, builder, _ := tracer.RetrieveDetails()

	if !reflect.DeepEqual("titanbuilder", builder) {
		t.Error("winning builder is not titanbuilder for block 18293308")
	}

}

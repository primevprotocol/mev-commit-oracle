package chaintracer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/primevprotocol/oracle/pkg/chaintracer"
)

func TestInfura(t *testing.T) {
	tracer := chaintracer.NewIncrementingTracer(18293308)
	for ; ; tracer.IncrementBlock() {
		blockData, builder, _ := tracer.RetrieveDetails()
		fmt.Println(blockData.Transactions[0])
		fmt.Println(builder)
		time.Sleep(1 * time.Second)
	}

}

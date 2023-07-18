package otlp

/*import (
	"context"
	"strings"
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"github.com/elastic/apm-data/model/modelpb"
	"go.opentelemetry.io/collector/pdata/plog"
	"golang.org/x/sync/semaphore"
)

func FuzzOtel(f *testing.F) {
	p := NewConsumer(ConsumerConfig{
		Semaphore: semaphore.NewWeighted(20),
		Processor: modelpb.ProcessBatchFunc(func(ctx context.Context, b *modelpb.Batch) error { return nil }),
	})

	f.Fuzz(func(t *testing.T, input []byte) {
		fuzzer := fuzz.NewConsumer(input)
		fuzzer.AllowUnexportedFields()
		fuzzer.DisallowUnknownTypes = true
		fuzzer.MaxDepth = 20
		fuzzer.NilChance = 0

		l := plog.Logs{}
		if err := fuzzer.GenerateStruct(&l); err != nil {
			if strings.Contains(err.Error(), "not enough bytes") {
				return
			}
			if strings.Contains(err.Error(), "json: unsupported value: ") {
				return
			}
			t.Fatalf("failed to generate data with input %v: %v", input, err)
			return
		}

		t.Logf("%+v", l.ResourceLogs())

		if err := p.ConsumeLogs(context.Background(), l); err != nil {
			//		if err != nil && !strings.Contains(err.Error(), "validation error") {
			t.Fatalf("failed to handle stream: %v", err)
		}
	})
}*/

package rumv3

import (
	"net/http"
	"testing"

	"github.com/elastic/apm-data/input/elasticapm/internal/modeldecoder/nullable"
	"github.com/elastic/apm-data/model/modelpb"
	fuzz "github.com/google/gofuzz"
)

func fuzzEvent[T any](mapFn func(*T, *modelpb.APMEvent)) func(*testing.T, []byte) {
	return func(t *testing.T, data []byte) {
		var e modelpb.APMEvent
		inputEvent := new(T)

		t.Cleanup(func() {
			if t.Failed() {
				t.Logf("failed to map input event: %+v", inputEvent)
			}
		})

		fuzzer := fuzz.NewFromGoFuzz(data).NilChance(.1).NumElements(1, 10).MaxDepth(5).AllowUnsupportedFields(true).Funcs(
			func(b *nullable.String, c fuzz.Continue) {
				if c.RandBool() {
					b.Set("str")
				}
			},
			func(b *nullable.Int, c fuzz.Continue) {
				if c.RandBool() {
					b.Set(1)
				}
			},
			func(b *nullable.Float64, c fuzz.Continue) {
				if c.RandBool() {
					b.Set(2.1)
				}
			},
			func(b *nullable.Bool, c fuzz.Continue) {
				if c.RandBool() {
					b.Set(true)
				}
			},
			func(b *nullable.HTTPHeader, c fuzz.Continue) {
				if c.RandBool() {
					b.Set(make(http.Header))
					c.Fuzz(&b.Val)
				}
			},
			func(b *nullable.TimeMicrosUnix, c fuzz.Continue) {
				if c.RandBool() {
					c.Fuzz(&b.Val)
					b.Set(b.Val)
				}
			},
			func(b *nullable.Interface, c fuzz.Continue) {
				if c.RandBool() {
					c.Fuzz(&b.Val)
					b.Set(b.Val)
				}
			},
		)
		fuzzer.Fuzz(inputEvent)
		mapFn(inputEvent, &e)
	}
}

func FuzzError(f *testing.F) {
	f.Fuzz(fuzzEvent(mapToErrorModel))
}

func FuzzSpan(f *testing.F) {
	f.Fuzz(fuzzEvent(mapToSpanModel))
}

func FuzzTransaction(f *testing.F) {
	f.Fuzz(fuzzEvent(mapToTransactionModel))
}

func FuzzMetricset(f *testing.F) {
	fn := func(m *transactionMetricset, e *modelpb.APMEvent) {
		mapToTransactionMetricsetModel(m, e)
	}
	f.Fuzz(fuzzEvent(fn))
}

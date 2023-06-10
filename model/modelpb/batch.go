package modelpb

import "context"

// BatchProcessor can be used to process a batch of events, giving the
// opportunity to update, add or remove events.
type BatchProcessor interface {
	// ProcessBatch is called with a batch of events for processing.
	//
	// Processing may involve anything, e.g. modifying, adding, removing,
	// aggregating, or publishing events.
	//
	// The caller should not assume the batch to be valid after the
	// method has returned.
	// If the batch needs to be processed asynchronously or kept around,
	// the processor must create a copy of the slice.
	ProcessBatch(context.Context, *Batch) error
}

// ProcessBatchFunc is a function type that implements BatchProcessor.
type ProcessBatchFunc func(context.Context, *Batch) error

// ProcessBatch calls f(ctx, b)
func (f ProcessBatchFunc) ProcessBatch(ctx context.Context, b *Batch) error {
	return f(ctx, b)
}

// Batch is a collection of APM events.
type Batch []*APMEvent

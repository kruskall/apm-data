package modelpb

func SpanProcessor() *Processor {
	return &Processor{
		Name:  "transaction",
		Event: "span",
	}
}

func (p *Processor) IsSpan() bool {
	return p.Name == "transaction" && p.Event == "span"
}

func TransactionProcessor() *Processor {
	return &Processor{
		Name:  "transaction",
		Event: "transaction",
	}
}

func (p *Processor) IsTransaction() bool {
	return p.Name == "transaction" && p.Event == "transaction"
}

func ErrorProcessor() *Processor {
	return &Processor{
		Name:  "error",
		Event: "error",
	}
}

func (p *Processor) IsError() bool {
	return p.Name == "error" && p.Event == "error"
}

func LogProcessor() *Processor {
	return &Processor{
		Name:  "log",
		Event: "log",
	}
}

func (p *Processor) IsLog() bool {
	return p.Name == "log" && p.Event == "log"
}

func MetricsetProcessor() *Processor {
	return &Processor{
		Name:  "metric",
		Event: "metric",
	}
}

func (p *Processor) IsMetricset() bool {
	return p.Name == "metric" && p.Event == "metric"
}

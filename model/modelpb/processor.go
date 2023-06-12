package modelpb

func SpanProcessor() *Processor {
	return &Processor{
		Name:  "transaction",
		Event: "span",
	}
}

func (p *Processor) IsSpan() bool {
	return p.GetName() == "transaction" && p.GetEvent() == "span"
}

func TransactionProcessor() *Processor {
	return &Processor{
		Name:  "transaction",
		Event: "transaction",
	}
}

func (p *Processor) IsTransaction() bool {
	return p.GetName() == "transaction" && p.GetEvent() == "transaction"
}

func ErrorProcessor() *Processor {
	return &Processor{
		Name:  "error",
		Event: "error",
	}
}

func (p *Processor) IsError() bool {
	return p.GetName() == "error" && p.GetEvent() == "error"
}

func LogProcessor() *Processor {
	return &Processor{
		Name:  "log",
		Event: "log",
	}
}

func (p *Processor) IsLog() bool {
	return p.GetName() == "log" && p.GetEvent() == "log"
}

func MetricsetProcessor() *Processor {
	return &Processor{
		Name:  "metric",
		Event: "metric",
	}
}

func (p *Processor) IsMetricset() bool {
	return p.GetName() == "metric" && p.GetEvent() == "metric"
}

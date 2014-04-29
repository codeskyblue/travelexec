package reporters

type multiReport struct {
	rs []Reporter
}

func NewMultiReporters(rs ...Reporter) Reporter {
	m := multiReport{rs}
	return &m
}

func (r *multiReport) Before(smy *SetupSummary) {
	for _, r := range r.rs {
		r.Before(smy)
	}
}

func (r *multiReport) After(smy *TeardownSummary) {
	for _, r := range r.rs {
		r.After(smy)
	}
}

func (r *multiReport) Close() {
	for _, r := range r.rs {
		r.Close()
	}
}

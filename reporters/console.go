package reportes

import (
	"fmt"
	"time"
)

type ConsoleReport struct {
	m      map[string]time.Time
	total  int
	failed int
}

func NewConsoleReport() *ConsoleReport {
	return &ConsoleReport{}
}

func (r *ConsoleReport) Before(smy *SetupSummary) {
	r.m[smy.Name] = time.Now()
	fmt.Printf("begin %s - %s\n", smy.Name, smy.Cmd)
}

func (r *ConsoleReport) After(smy *TeardownSummary) {
	timeused := time.Since(r.m[smy.Name])
	_ = timeused.Seconds()
	fmt.Println("time used: %s, %s", timeused, smy.Error)
	r.total += 1
	if smy.Error != nil {
		r.failed += 1
	}
}

func (r *ConsoleReport) Close() {
	fmt.Printf("close console report %d/%d\n", r.failed, r.total)
}

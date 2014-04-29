package reporters

type SetupSummary struct {
	Name string
	Cmd  string
}

type TeardownSummary struct {
	Name   string
	Error  error
	Output string
}

type Reporter interface {
	Before(*SetupSummary)
	After(*TeardownSummary)
	Close()
}

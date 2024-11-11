package analyze

type ResourceDumpAnalyzer struct {
	data DataSource
}

func NewResourceDumpAnalyzer(data DataSource) *ResourceDumpAnalyzer {
	return &ResourceDumpAnalyzer{data: data}
}

func (a *ResourceDumpAnalyzer) Run() error {
	return nil
}

func (a *ResourceDumpAnalyzer) Shutdown() {
}

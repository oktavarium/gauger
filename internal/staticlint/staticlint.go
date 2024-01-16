package staticlint

import (
	"strings"

	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"honnef.co/go/tools/staticcheck"
)

// Run - функция инициализации и запуска линтера
func Run() {
	var analyzers []*analysis.Analyzer
	analyzers = append(
		analyzers,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		sortslice.Analyzer,
		unmarshal.Analyzer,
		errcheck.Analyzer,
		ExitCheckAnalyzer,
	)

	//правила для statickcheck
	checks := map[string]bool{
		"SA":    true,
		"S1016": true,
		"S1028": true,
	}
	for _, v := range staticcheck.Analyzers {
		for k, ok := range checks {
			if ok {
				if strings.HasPrefix(v.Analyzer.Name, k) {
					analyzers = append(analyzers, v.Analyzer)
				}
				continue
			}
		}
	}

	multichecker.Main(
		analyzers...,
	)
}

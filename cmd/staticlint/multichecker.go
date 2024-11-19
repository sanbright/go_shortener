package main

import (
	"fmt"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"log"
	"sanbright/go_shortener/cmd/staticlint/noexitanalyzer"
	"sanbright/go_shortener/cmd/staticlint/revive"
	"strings"
)

// main - мультичекер
func main() {
	var mychecks []*analysis.Analyzer

	mychecks = append(mychecks,
		noexitanalyzer.NoExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" || a.Analyzer.Name == "ST1000" {
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	reviveOutput, err := revive.RunRevive()
	if err != nil {
		log.Fatalf("ошибка запуска revive: %v\n%s", err, reviveOutput)
	}

	fmt.Println("Revive Output:")
	fmt.Println(strings.TrimSpace(reviveOutput))

	multichecker.Main(
		mychecks...,
	)
}

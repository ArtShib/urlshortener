// Модуль staticlint собирается в multichecker, состоящий из:
//
// - статических анализаторов пакета golang.org/x/tools/go/analysis/passes,
//
// - анализаторов SA* пакета staticcheck.io
//
// - анализатора exitcheck, который проверяет, что в функции main() нет вызова os.Exit
//
// Для сборки multichecker вызвать команду go build.
// Запуск командой ./staticlint.exe ./...
//
// по умолчанию запускаются все анализаторы, настроить можно с помощью флагов:
//
// выбор отдельных анализаторов: -ИМЯ_АНАЛИЗАТОРА
//
// исключение анализатора из проверки: -ИМЯ_АНАЛИЗАТОРА=false
package main

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// в методе main добавляем все необходимые анализаторы в multichecker.main
func main() {
	checks := map[string]bool{
		"ST1017": true,
		"S1010":  true,
	}

	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") || checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	mychecks = append(mychecks, printf.Analyzer)
	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, structtag.Analyzer)

	mychecks = append(mychecks, ExitCheckAnalyzer)

	fmt.Println(mychecks)

	multichecker.Main(
		mychecks...,
	)
}

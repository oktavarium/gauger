package staticlint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ExitCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./...")
}

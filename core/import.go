package core

type Package struct {
	file     string
	isImport bool
}

var gImportList map[string]*Package

func doImport(s string) string {
	return ""
}

func doScan() {

}

func doPreprocessing(s string) {
}

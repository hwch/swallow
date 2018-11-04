package core

type Package struct {
	file     string
	isImport bool
}

var g_import_list map[string]*Package

func doImport(s string) string {
	return ""
}

func doScan() {

}

func doPreprocessing(s string) {
}

package testPackage

/*
	import (
		"not_exist_package1"
	)
*/
func ImportTester() string {
	return `import (
		"not_exist_package"
	)` + "import ()"
}

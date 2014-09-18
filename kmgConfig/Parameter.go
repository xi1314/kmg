package kmgConfig

type Parameter struct {
	DatabaseUsername   string
	DatabasePassword   string
	DatabaseHost       string
	DatabaseDbName     string
	DatabaseTestDbName string

	MemcacheHostList []string

	SessionPrefix     string
	SessionExpiration string
}

package common

type CheckType int

const (
	InvalidType  CheckType = 0
	SqlException           = 1 << 0
	Sql                    = 1 << 1
	AllType                = Sql | SqlException
)

var buildinCheckTypes = []CheckType{SqlException}

func CheckTypeToString(ct CheckType) string {
	switch ct {
	case SqlException:
		return "sql_exception"
	case Sql:
		return "sql"
	default:
		return "unknown"
	}
}

func CheckStringToType(key string) CheckType {
	switch key {
	case "sql_exception":
		return SqlException
	case "sql":
		return Sql
	case "all":
		return AllType
	default:
		return InvalidType
	}
}

func BuildinActionScript() string {
	if len(buildinCheckTypes) > 0 {
		var bcond string
		for _, ct := range buildinCheckTypes {
			bcond += (" && key === '" + CheckTypeToString(ct) + "'")
		}
		script := `JSON.stringify(Object.keys(RASP.algorithmConfig || {})
		.filter(key => typeof key === 'string' && typeof RASP.algorithmConfig[key] === 'object' && typeof RASP.algorithmConfig[key].action === 'string'`
		script += bcond
		script += `).map(key => [key, RASP.algorithmConfig[key].action]))`
		return script
	}
	return ""
}

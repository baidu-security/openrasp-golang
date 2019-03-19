package common

type CheckType int

const (
	InvalidType  CheckType = 0
	SqlException           = 1 << 0
	Sql                    = 1 << 1
	AllType                = Sql | SqlException
)

var checkTypeName = map[CheckType]string{
	Sql:          "sql",
	SqlException: "sql_exception",
	AllType:      "all",
}

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

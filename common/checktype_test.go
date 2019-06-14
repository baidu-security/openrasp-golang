package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTypeToString(t *testing.T) {
	assert.Equal(t, CheckTypeToString(Sql), "sql", "they should be equal")
	assert.Equal(t, CheckTypeToString(SqlException), "sql_exception", "they should be equal")
	assert.Equal(t, CheckTypeToString(InvalidType), "unknown", "they should be equal")
}

func TestCheckStringToType(t *testing.T) {
	assert.EqualValues(t, CheckStringToType("sql"), Sql, "they should be equal")
	assert.EqualValues(t, CheckStringToType("sql_exception"), SqlException, "they should be equal")
	assert.EqualValues(t, CheckStringToType("all"), AllType, "they should be equal")
	assert.EqualValues(t, CheckStringToType("doom"), InvalidType, "they should be equal")
}

func TestBuildinActionScript(t *testing.T) {
	script := BuildinActionScript()
	assert.Equal(t, script, "JSON.stringify(Object.keys(RASP.algorithmConfig || {})\n\t\t.filter(key => typeof key === 'string' && typeof RASP.algorithmConfig[key] === 'object' && typeof RASP.algorithmConfig[key].action === 'string' && key === 'sql_exception').map(key => [key, RASP.algorithmConfig[key].action]))", "they should be equal")
	buildinCheckTypes = []CheckType{}
	script = BuildinActionScript()
	assert.Equal(t, script, "", "they should be equal")
}

package testhelper

import (
	"fmt"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

// Dialect containers metadata that differ across SQL database
type Dialect int

// Enum values for Dialect
const (
	DialectUnknown Dialect = iota
	DialectMySQL
	DialectPostgres
)

func (d Dialect) String() string {
	switch d {
	case DialectMySQL:
		return "mysql"
	case DialectPostgres:
		return "postgres"
	}
	return ""
}

type dialectParams struct {
	repo   string
	tag    string
	envs   []string
	urlFmt string
	portID string
}

func (p *dialectParams) URL(res *dockertest.Resource) string {
	return fmt.Sprintf(p.urlFmt, res.GetPort(p.portID))
}

var databaseName = "gormzap_test"
var databasePass = "secret"

var dialectParamsByDialect = map[Dialect]dialectParams{
	DialectPostgres: {
		repo: "postgres",
		tag:  "9.6.5-alpine",
		envs: []string{
			"POSTGRES_PASSWORD=" + databasePass,
			"POSTGRES_DB=" + databaseName,
		},
		urlFmt: fmt.Sprintf("postgres://postgres:%s@localhost:%%s/%s?sslmode=disable", databasePass, databaseName),
		portID: "5432/tcp",
	},
	DialectMySQL: {
		repo: "mysql",
		tag:  "5.7.19-alpine",
		envs: []string{
			"MYSQL_ROOT_PASSWORD=" + databasePass,
		},
		urlFmt: fmt.Sprintf("root:%s@(localhost:%%s)/mysql", databasePass),
		portID: "3306/tcp",
	},
}

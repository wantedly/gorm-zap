package testhelper

import (
	"database/sql"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

// MustCreatePool creates new docker remove API client instance
func MustCreatePool() *DockerPool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	return &DockerPool{
		pool: pool,
	}
}

// DockerPool wraps dockertest.Pool
type DockerPool struct {
	pool *dockertest.Pool
}

// DockerConnection contains connections to a docker container
type DockerConnection struct {
	MustClose func()
	URL       string
	Dialect   string
}

// MustCreateDB creates new database container
func (dp *DockerPool) MustCreateDB(d Dialect) *DockerConnection {
	params, ok := dialectParamsByDialect[d]
	if !ok {
		panic("Unknown dialect")
	}

	res, err := dp.pool.Run(params.repo, params.tag, params.envs)
	if err != nil {
		panic(err)
	}

	url := params.URL(res)

	if err := dp.pool.Retry(func() error {
		var err error
		db, err := sql.Open(d.String(), url)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		panic(err)
	}

	return &DockerConnection{
		URL:     url,
		Dialect: d.String(),
		MustClose: func() {
			if err := dp.pool.Purge(res); err != nil {
				panic(err)
			}
		},
	}
}

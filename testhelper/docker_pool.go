package testhelper

import (
	"database/sql"

	dockertest "gopkg.in/ory-am/dockertest.v3"
)

func MustCreatePool() *DockerPool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	return &DockerPool{
		pool: pool,
	}
}

type DockerPool struct {
	pool *dockertest.Pool
}

type DockerConnection struct {
	MustClose func()
	URL       string
	Dialect   string
}

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

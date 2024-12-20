package storage

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type Storage struct {
	logger *zerolog.Logger
	engine *xorm.Engine
}

func NewStorage(engineType, engineUrl string, logger *zerolog.Logger) (*Storage, error) {
	if logger == nil {
		logger = &log.Logger
	}

	eng, err := xorm.NewEngine(parseDriverName(engineType), engineUrl)
	if err != nil {
		return nil, err
	}
	eng.SetLogger(&xormLogger{logger: logger})

	if err := eng.Ping(); err != nil {
		return nil, err
	}

	return &Storage{logger: logger, engine: eng}, nil
}

func parseDriverName(engineType string) string {
	if engineType == "postgresql" {
		return "pgx"
	}
	if engineType == "mariadb" {
		return "mysql"
	}
	return engineType
}

func (db *Storage) Init() error {
	return db.engine.CreateTables(&Account{})
}

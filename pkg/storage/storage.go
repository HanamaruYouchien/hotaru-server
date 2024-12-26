package storage

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"

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

	cacher := caches.NewLRUCacher(caches.NewMemoryStore(), 1000)
	eng.SetDefaultCacher(cacher)

	if err := eng.Ping(); err != nil {
		return nil, err
	}

	// sync db
	if err := eng.Sync(&Account{}, &Device{}, &Profile{}, &Room{}, &RoomAlias{}, &AccountRoom{}); err != nil {
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
	if err := db.engine.Sync(&Account{}, &Device{}, &Profile{}, &Room{}, &RoomAlias{}, &AccountRoom{}); err != nil {
		return err
	}
	db.engine.Exec("ALTER TABLE device ADD CONSTRAINT device_localpart_fkey FOREIGN KEY (localpart) REFERENCES account(localpart) ON DELETE CASCADE ON UPDATE CASCADE;")
	db.engine.Exec("ALTER TABLE profile ADD CONSTRAINT profile_localpart_fkey FOREIGN KEY (localpart) REFERENCES account(localpart) ON DELETE CASCADE ON UPDATE CASCADE;")
	db.engine.Exec("ALTER TABLE room_alias ADD CONSTRAINT room_alias_room_id_fkey FOREIGN KEY (room_id) REFERENCES room(room_id) ON DELETE CASCADE ON UPDATE CASCADE;")
	db.engine.Exec("ALTER TABLE account_room ADD CONSTRAINT account_room_localpart_fkey FOREIGN KEY (localpart) REFERENCES account(localpart) ON DELETE CASCADE ON UPDATE CASCADE;")
	db.engine.Exec("ALTER TABLE account_room ADD CONSTRAINT account_room_room_id_fkey FOREIGN KEY (room_id) REFERENCES room(room_id) ON DELETE CASCADE ON UPDATE CASCADE;")
	return nil
}

package db

import (
	"github.com/gocraft/dbr/v2"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"path"
)

// NewSqlite 创建一个 sqlite db [filepath]db存储路径 [sqlDir]sql脚本路径
func NewSqlite(filepath, sqlDir string) *dbr.Session {
	err := os.Mkdir(path.Dir(filepath), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	conn, err := dbr.Open("sqlite3", filepath, nil)
	if err != nil {
		panic(err)
	}
	session := conn.NewSession(nil)
	migrations := &migrate.FileMigrationSource{
		Dir: sqlDir,
	}
	_, err = migrate.Exec(session.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	return session
}

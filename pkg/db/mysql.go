package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/gocraft/dbr/v2"
	migrate "github.com/rubenv/sql-migrate"
	"net/http"
	"sort"
	"strings"
	"time"
)

// NewMySQL 创建一个MySQL db，[addr]db 存储路径 [sqlDir]sql脚本目录 [migration]标志是否需要执行数据库迁移
func NewMySQL(addr string, sqlDir string, migration bool) *dbr.Session {
	// 打开一个 MySQL 数据库连接
	conn, err := dbr.Open("mysql", addr, nil)
	if err != nil {
		panic(err)
	}
	// 最大的连接数
	conn.SetMaxOpenConns(2000)
	// 连接池中最大空闲连接数
	conn.SetMaxIdleConns(2000)
	// mysql 默认超时时间为 60*60*8=28800 SetConnMaxLifetime 设置为小于数据库超时时间
	conn.SetConnMaxLifetime(time.Second * 60 * 60 * 4)
	// 验证数据库的连接是否正常
	conn.Ping()

	// 创建一个新的会话对象
	session := conn.NewSession(nil)

	// 如果需要进行数据迁移
	if migration {
		// 数据迁移
		Migration(sqlDir, session)
	}

	return session
}

// Migration 数据迁移
func Migration(sqlDir string, session *dbr.Session) {
	migrations := &FileDirMigrationSource{
		Dir: sqlDir,
	}

	_, err := migrate.Exec(session.DB, "mysql", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
}

type byID []*migrate.Migration

func (b byID) Len() int           { return len(b) }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byID) Less(i, j int) bool { return b[i].Less(b[j]) }

// FileDirMigrationSource 文件目录源 遇到目录进行递归获取
type FileDirMigrationSource struct {
	Dir string
}

// FindMigrations 查找数据迁移脚本
func (f FileDirMigrationSource) FindMigrations() ([]*migrate.Migration, error) {
	// 创建一个文件系统抽象，表示从f.Dir(sql脚本存放的目录)
	filesystem := http.Dir(f.Dir)
	// 初始化切片
	migrations := make([]*migrate.Migration, 0, 100)
	err := f.findMigrations(filesystem, &migrations)
	if err != nil {
		return nil, err
	}

	// 确保 sql脚本是有序的
	sort.Sort(byID(migrations))

	return migrations, nil
}

// 查找所有sql脚本
func (f FileDirMigrationSource) findMigrations(dir http.FileSystem, migrations *[]*migrate.Migration) error {
	// 打开当前目录的根路径
	file, err := dir.Open("/")
	if err != nil {
		return err
	}

	// 读取所有文件信息
	files, err := file.Readdir(0)
	if err != nil {
		return err
	}

	for _, info := range files {
		// 查找所有以.sql结尾的文件
		if strings.HasSuffix(info.Name(), ".sql") {
			file, err := dir.Open(info.Name())
			if err != nil {
				return fmt.Errorf("Error while opening %s: %s", info.Name(), err)
			}

			migration, err := migrate.ParseMigration(info.Name(), file)
			if err != nil {
				return fmt.Errorf("Error while parsing %s: %s", info.Name(), err)
			}

			// 添加结果到切片
			*migrations = append(*migrations, migration)
		} else if info.IsDir() {
			// 如果是目录，递归调用findMigrations，继续查找子目录中的迁移脚本
			err = f.findMigrations(http.Dir(fmt.Sprintf("%s/%s", f.Dir, info.Name())), migrations)
		}
	}

	return nil
}

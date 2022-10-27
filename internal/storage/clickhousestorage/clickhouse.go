package clickhousestorage

import (
	"database/sql"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	_ "github.com/mailru/go-clickhouse/v2"
	"log"
)

type StorageClickHouse struct {
	conn *sql.DB
}

func initConnection(conn *sql.DB) *StorageClickHouse {
	return &StorageClickHouse{conn: conn}
}

func Init() (*StorageClickHouse, error) {
	connect, err := sql.Open("chhttp", "http://default:asdqwe123@127.0.0.1:8123")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		log.Fatal(err)
	}

	return initConnection(connect), nil
}

func (s *StorageClickHouse) Add(item storage.Item) (string, error) {
	tx, err := s.conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO items (*) VALUES (?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(item.Id, item.CampaignId, item.Name,
		item.Description, item.Priority, item.Removed, item.CreatedAt); err != nil {
		log.Fatal(err)
		return "YOU LOX", err
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return "YOU LOX2", err
	}
	return "Ok with add to click", nil
}

package clickhousestorage

import (
	"github.com/GeorgeShibanin/hezzl_test5/internal/config"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/roistat/go-clickhouse"
	"log"
)

const (
	PatchItemsLog  = `INSERT INTO items VALUES ()`
	AddItemsLog    = ``
	DeleteItemsLog = ``
)

type StorageClickHouse struct {
	conn *clickhouse.Conn
}

func initConnection(conn *clickhouse.Conn) *StorageClickHouse {
	return &StorageClickHouse{conn: conn}
}

func Init() (*StorageClickHouse, error) {
	clickHouseClient := clickhouse.NewHttpTransport()
	clickHouseConn := clickhouse.NewConn(config.ClickHouse, clickHouseClient)
	err := clickHouseConn.Ping()
	log.Println("connected to clickhouse")
	if err != nil {
		panic(err)
	}
	return initConnection(clickHouseConn), nil
}

func (s *StorageClickHouse) Add(item storage.Item) (string, error) {
	query, err := clickhouse.BuildInsert("items",
		clickhouse.Columns{"id", "campaign_id", "name", "description", "priority", "removed", "created_ad"},
		clickhouse.Row{item.Id, item.CampaignId, item.Name, item.Description, item.Priority, item.Removed, item.CreatedAt},
	)
	if err != nil {
		log.Println("LOL HERE", err)
	}
	err = query.Exec(s.conn)
	if err != nil {
		log.Println(err)
	}

	if err == nil {
		err = query.Exec(s.conn)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			return "OK", nil
		}
		log.Println("trouble while add to clickhouse")
		return "NO YOU FAILED", err
	}
	log.Println("trouble while add to clickhouse")
	return "NO YOU FAILED", err
}

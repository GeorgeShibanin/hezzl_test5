package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	GetAllItems = `SELECT * FROM items`

	InsertItem = `INSERT INTO items (campaign_id, name, description, removed, created_at) 
						VALUES ($1, $2, $3, $4, $5) 
						RETURNING id, campaign_id, name, description, priority, removed, created_at`

	PatchItems = `UPDATE items SET name = $1, description = $2 WHERE id = $3 AND campaign_id = $4
 						RETURNING id, campaign_id, name, description, priority, removed, created_at`

	GetById = `SELECT id, description FROM items WHERE id = $1 AND campaign_id = $2`

	DeleteItem = `DELETE FROM items WHERE id = $1 AND campaign_id = $2
						RETURNING id, campaign_id, name, description, priority, removed, created_at`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

type StoragePostgres struct {
	conn *pgx.Conn
}

func initConnection(conn *pgx.Conn) *StoragePostgres {
	return &StoragePostgres{conn: conn}
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*StoragePostgres, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}

	return initConnection(conn), nil
}

func (s *StoragePostgres) PostItem(ctx context.Context, campaignId storage.CampaignId, name storage.Name) (storage.Item, error) {
	item := storage.Item{}

	err := s.conn.QueryRow(ctx, InsertItem, campaignId, name, "interesting", false, time.Now().UTC().Format(time.RFC3339)).
		Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Item{}, fmt.Errorf("item not exist - %w", storage.StorageError)
		}
		return storage.Item{}, fmt.Errorf("cant add new item - %w", err)
	}

	return item, nil
}

func (s *StoragePostgres) PatchItem(ctx context.Context, id storage.Id, campaignId storage.CampaignId, name storage.Name, description storage.Description, flag int) (storage.Item, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.Item{}, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	item := storage.Item{}
	err = tx.QueryRow(ctx, GetById, id, campaignId).Scan(&item.Id, &item.Description)
	if err != nil {
		return storage.Item{}, fmt.Errorf("errors.item.notFound - %w", err)
	}
	if flag == 0 {
		item.Description = string(description)
	}

	err = tx.QueryRow(ctx, PatchItems, name, &item.Description, id, campaignId).
		Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Item{}, fmt.Errorf("item not exist - %w", storage.StorageError)
		}
		return storage.Item{}, fmt.Errorf("cant add new item - %w", err)
	}
	return item, nil
}

func (s *StoragePostgres) GetItems(ctx context.Context) ([]storage.Item, error) {
	items := make([]storage.Item, 0)
	rows, err := s.conn.Query(ctx, GetAllItems)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := storage.Item{}
		if err := rows.Scan(&item.Id, &item.CampaignId, &item.Name,
			&item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {
			log.Fatal(err)
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *StoragePostgres) DeleteItem(ctx context.Context, id storage.Id, campaignId storage.CampaignId) (storage.Item, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.Item{}, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	item := &storage.Item{}
	err = tx.QueryRow(ctx, GetById, id, campaignId).Scan(&item.Id, &item.Description)
	if err != nil {
		//Оформить ошибку по пункту 7 в задании
		return storage.Item{}, fmt.Errorf("errors.item.notFound - %w", err)
	}
	err = tx.QueryRow(ctx, DeleteItem, id, campaignId).
		Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt)
	if err != nil {
		//Обработай ошибку
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Item{}, fmt.Errorf("item not exist - %w", storage.StorageError)
		}
		return storage.Item{}, fmt.Errorf("cant delete item - %w", err)
	}
	item.Removed = true
	return *item, nil
}

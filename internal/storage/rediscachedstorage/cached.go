package rediscachedstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GeorgeShibanin/hezzl_test5/internal/storage"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Storage struct {
	conn   storage.Storage
	client *redis.Client
}

func Init(redisClient *redis.Client, persistentStorage storage.Storage) (*Storage, error) {
	return &Storage{
		conn:   persistentStorage,
		client: redisClient,
	}, nil
}

func (s *Storage) GetItems(ctx context.Context) ([]storage.Item, error) {
	//получаем значение из cache
	get := s.client.Get(ctx, "items")
	switch items, err := get.Result(); {
	case err == redis.Nil:
		// continue execution
	case err != nil:
		return nil, fmt.Errorf("%w: failed to get value from redis due to error %s", storage.StorageError, err)
	default:
		log.Printf("Successfully obtained items from cache")
		data := []storage.Item{}
		json.Unmarshal([]byte(items), &data)
		return data, nil
	}
	log.Printf("Loading post from persistent storage")
	//получаем значение из базы если значение нет в cache
	allItems, err := s.conn.GetItems(ctx)
	if err != nil {
		return allItems, err
	}
	//устанавливаем значение в cache
	allItemsMarshall, err := json.Marshal(allItems)
	if err != nil {
		//Обработать ошибку
		log.Fatal(err)
		return nil, err
	}
	err = s.client.Set(ctx, "items", allItemsMarshall, time.Minute).Err()
	if err != nil {
		log.Printf("Failed to insert key items into cache due to an error: %s\n", err)
	}
	return allItems, nil
}

func (s *Storage) DeleteItem(ctx context.Context, id storage.Id, campaignId storage.CampaignId) (storage.Item, error) {
	return s.conn.DeleteItem(ctx, id, campaignId)
}

func (s *Storage) PatchItem(ctx context.Context, id storage.Id, campaignId storage.CampaignId, name storage.Name, description storage.Description) (storage.Item, error) {
	return s.conn.PatchItem(ctx, id, campaignId, name, description)
}

func (s *Storage) PostItem(ctx context.Context, campaignId storage.CampaignId, name storage.Name) (storage.Item, error) {
	return s.conn.PostItem(ctx, campaignId, name)
}

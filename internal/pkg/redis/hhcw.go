package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"kawalrealcount/internal/data/model"
	"time"
)

func (repo repo) GetHHCW(key string) (model.HHCWEntity, error) {
	ctx := context.Background()

	res, err := repo.client.Get(ctx, key).Bytes()
	if err != nil {
		return model.HHCWEntity{}, err
	}

	var data model.HHCWEntity

	if err = json.Unmarshal(res, &data); err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return model.HHCWEntity{}, err
	}

	return data, nil
}

func (repo repo) PutHHCW(key string, data model.HHCWEntity, expiry time.Duration) error {
	ctx := context.Background()

	// Marshal the Person struct into binary data
	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
		return nil
	}

	return repo.client.Set(ctx, key, buf, expiry).Err()
}

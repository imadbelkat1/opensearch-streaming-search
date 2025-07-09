package services

import (
	"context"
	"log"

	"internship-project/internal/models"
)

type UpdateApiService struct {
	Client *HackerNewsApiClient
}

func NewUpdateApiService(client *HackerNewsApiClient) *UpdateApiService {
	return &UpdateApiService{Client: client}
}

func (s *UpdateApiService) FetchUpdates(ctx context.Context) (models.Update, error) {
	var update models.Update
	err := s.Client.Get(ctx, "/updates.json", update)
	if err != nil {
		log.Printf("Error fetching updates: %v", err)
		return models.Update{}, err
	}

	return update, nil
}

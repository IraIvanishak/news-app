package services

import (
	"context"
	"fmt"

	"github.com/IraIvanishak/news-app/internal/adapter/external/unsplash"
	"github.com/IraIvanishak/news-app/internal/core/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type UnsplashService struct {
	client *unsplash.Client
	logger *zap.Logger
}

var _ ports.UnsplashService = (*UnsplashService)(nil)

// NewUnsplashService creates a new Unsplash service
func NewUnsplashService(client *unsplash.Client, logger *zap.Logger) *UnsplashService {
	return &UnsplashService{
		client: client,
		logger: logger,
	}
}

// SearchPhotos searches for photos on Unsplash and transforms them to domain objects
func (s *UnsplashService) SearchPhotos(ctx context.Context, query string) ([]ports.UnsplashPhoto, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Call Unsplash API with landscape orientation and 9 photos per page
	response, err := s.client.SearchPhotos(ctx, query, "landscape", 9)
	if err != nil {
		s.logger.Error("Failed to search photos on Unsplash",
			zap.String("query", query),
			zap.Error(err),
		)
		return nil, fmt.Errorf("photo search temporarily unavailable")
	}

	// Transform API response to domain objects
	photos := make([]ports.UnsplashPhoto, 0, len(response.Results))
	appName := viper.GetString("UNSPLASH_APP_NAME")

	for _, result := range response.Results {
		// Append UTM parameters to photographer URL
		photographerURL := fmt.Sprintf("%s?utm_source=%s&utm_medium=referral",
			result.User.Links.HTML,
			appName,
		)

		photo := ports.UnsplashPhoto{
			ID:               result.ID,
			RegularURL:       result.URLs.Regular,
			SmallURL:         result.URLs.Small,
			ThumbURL:         result.URLs.Thumb,
			Description:      getDescription(result),
			PhotographerName: result.User.Name,
			PhotographerURL:  photographerURL,
			DownloadLocation: result.Links.DownloadLocation,
		}
		photos = append(photos, photo)
	}

	s.logger.Info("Successfully fetched photos from Unsplash",
		zap.String("query", query),
		zap.Int("count", len(photos)),
	)

	return photos, nil
}

// TriggerDownload triggers the download endpoint to track photo usage
func (s *UnsplashService) TriggerDownload(ctx context.Context, downloadURL string) error {
	if downloadURL == "" {
		s.logger.Warn("Download URL is empty, skipping download tracking")
		return nil
	}

	err := s.client.TriggerDownload(ctx, downloadURL)
	if err != nil {
		// Log but don't fail the operation
		s.logger.Warn("Failed to trigger download tracking",
			zap.String("downloadURL", downloadURL),
			zap.Error(err),
		)
		return nil // Non-blocking error
	}

	s.logger.Info("Successfully triggered download tracking",
		zap.String("downloadURL", downloadURL),
	)

	return nil
}

// getDescription returns the best available description for a photo
func getDescription(result unsplash.PhotoResult) string {
	if result.Description != "" {
		return result.Description
	}
	if result.AltDesc != "" {
		return result.AltDesc
	}
	return "Photo from Unsplash"
}

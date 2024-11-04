package link_generator

import (
	"context"
	"errors"
	"main/internal/storage"
	"main/internal/utils"
)

type LinkGenerator struct {
	storage storage.Storage
}

func NewLinkGenerator(storage storage.Storage) *LinkGenerator {
	return &LinkGenerator{storage: storage}
}

// TODO: add tests
func (lg *LinkGenerator) GenerateLinkAlias(ctx context.Context, link string, lifetimeSeconds int, alias string, linkLength int) (string, error) {
	if alias == "" {
		for i := 0; i < 5; i++ {
			alias = utils.RandomString(linkLength)
			exists, err := lg.storage.Exists(ctx, alias)
			if err != nil {
				return "", err
			}
			if !exists {
				break
			}
			if i == 4 {
				return "", errors.New("cannot create unique alias for encoding")
			}
		}
	} else {
		exists, err := lg.storage.Exists(ctx, alias)
		if err != nil {
			return "", err
		}
		if exists {
			return "", errors.New("alias already exists")
		}
	}

	err := lg.storage.Set(ctx, alias, link, lifetimeSeconds)
	if err != nil {
		return "", err
	}

	return alias, nil
}

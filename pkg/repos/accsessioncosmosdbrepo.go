package repos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
)

type CosmosSessionRepo struct {
	container *azcosmos.ContainerClient
}

func NewCosmosSessionRepo(connStr, dbName, containerName string) (*CosmosSessionRepo, error) {
	client, err := azcosmos.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating cosmos client: %w", err)
	}

	dbClient, err := client.NewDatabase(dbName)
	if err != nil {
		return nil, fmt.Errorf("getting database client: %w", err)
	}

	containerClient, err := dbClient.NewContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("getting container client: %w", err)
	}

	return &CosmosSessionRepo{container: containerClient}, nil
}

func (r *CosmosSessionRepo) UpsertSessions(sessions []*models.Session) error {
	for _, s := range sessions {
		_, err := r.UpsertSession(s)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (r *CosmosSessionRepo) UpsertSession(session *models.Session) (string, error) {
	ctx := context.Background()

	data, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshalling session: %w", err)
	}

	pk := azcosmos.NewPartitionKeyString(session.Player)

	_, err = r.container.UpsertItem(ctx, pk, data, nil)
	if err != nil {
		return "", fmt.Errorf("upserting item: %w", err)
	}

	return session.Id, nil
}

func (r *CosmosSessionRepo) GetSessions(limit int, offset int) ([]models.Session, error) {
	ctx := context.Background()
	query := "SELECT * FROM c ORDER BY c.startTime DESC"

	queryPager := r.container.NewQueryItemsPager(
		query,
		azcosmos.NewPartitionKeyString("anonymous"),
		&azcosmos.QueryOptions{
			PageSizeHint: int32(limit),
		},
	)

	var sessions []models.Session
	skipped := 0

	for queryPager.More() {
		page, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("querying sessions: %w", err)
		}

		for _, item := range page.Items {
			if offset > 0 && skipped < offset {
				skipped++
				continue
			}

			var s models.Session
			if err := json.Unmarshal(item, &s); err != nil {
				return nil, fmt.Errorf("unmarshalling session: %w", err)
			}

			sessions = append(sessions, s)

			if limit != -1 && len(sessions) >= limit {
				return sessions, nil
			}
		}
	}

	return sessions, nil
}

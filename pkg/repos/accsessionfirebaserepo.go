package repos

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
	"google.golang.org/api/iterator"
)

type FirebaseSessionRepo struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewFirebaseSessionRepo(
	projectID string,
	databaseName string,
	collectionName string,
) (*FirebaseSessionRepo, error) {
	client, err := firestore.NewClientWithDatabase(context.Background(), projectID, databaseName)
	if err != nil {
		return nil, fmt.Errorf(
			"initializing Firestore client with database %q: %w",
			databaseName,
			err,
		)
	}

	collection := client.Collection(collectionName)

	return &FirebaseSessionRepo{
		client:     client,
		collection: collection,
	}, nil
}

func (r *FirebaseSessionRepo) UpsertSessions(sessions []*models.Session) error {
	for _, s := range sessions {
		if _, err := r.upsertSession(s); err != nil {
			// TODO: plaster
			if s.LapsCompleted == 0 && len(s.Laps)-1 > 0 {
				fmt.Printf("HAD TO PUT A PLASTER ON THE NUMBER OF LAPS ISSUE")
				s.LapsCompleted = int32(len(s.Laps)) - 1
			}
			fmt.Printf("%v", err)
			return err
		}
	}
	return nil
}

func (r *FirebaseSessionRepo) upsertSession(session *models.Session) (string, error) {
	ctx := context.Background()

	data, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshalling session: %w", err)
	}

	var docData map[string]interface{}
	if err := json.Unmarshal(data, &docData); err != nil {
		return "", fmt.Errorf("unmarshalling to map: %w", err)
	}

	_, err = r.collection.Doc(session.Id).Set(ctx, docData)
	if err != nil {
		return "", fmt.Errorf("upserting session to Firestore: %w", err)
	}

	return session.Id, nil
}

func (r *FirebaseSessionRepo) GetRecentSessions() ([]*models.Session, error) {
	ctx := context.Background()
	query := r.collection.OrderBy("startTime", firestore.Desc).Limit(20)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting recent sessions: %w", err)
	}

	sessions := make([]*models.Session, 0, len(docs))
	for _, doc := range docs {
		var session models.Session
		data, err := json.Marshal(doc.Data())
		if err != nil {
			return nil, fmt.Errorf("marshalling Firestore data: %w", err)
		}
		if err := json.Unmarshal(data, &session); err != nil {
			return nil, fmt.Errorf("unmarshalling to session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *FirebaseSessionRepo) CleanUpSessions() error {
	ctx := context.Background()
	query := r.collection.
		Where("isActive", "==", false).
		Where("lapsCompleted", "<=", 5)

	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

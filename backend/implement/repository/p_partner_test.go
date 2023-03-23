package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sinulingga23/dynamic-form/backend/db"
	"github.com/sinulingga23/dynamic-form/backend/model"
	"github.com/stretchr/testify/assert"
)

func TestPPartnerRepositoryImpl_Create(t *testing.T) {
	db, errConnectDB := db.ConnectDB()
	if errConnectDB != nil {
		t.Fatalf("errConnectDB: %v", errConnectDB)
	}

	pPartnerRepositoryImpl := NewPPartnerRepositoryImpl(db)

	t.Run("should success insert an object into p_partner table.", func(t *testing.T) {
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.PPartner{
			Id:          uuid.NewString(),
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
			UpdatedAt:   sql.NullTime{Valid: false},
		})

		assert.Nil(t, errCreate)
	})

	t.Run("should failed insert an object into p_partner table when id is empty.", func(t *testing.T) {
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.PPartner{
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
			UpdatedAt:   sql.NullTime{Valid: false},
		})

		assert.Error(t, errCreate)
	})

	t.Run("should failed insert an object into p_partner table when id is not uuid.", func(t *testing.T) {
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.PPartner{
			Id:          "34739473949",
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
			UpdatedAt:   sql.NullTime{Valid: false},
		})

		assert.Error(t, errCreate)
	})
}

func TestPPartnerRepositoryImpl_FindOne(t *testing.T) {
	db, errConnectDB := db.ConnectDB()
	if errConnectDB != nil {
		t.Fatalf("errConnectDB: %v", errConnectDB)
	}

	pPartnerRepositoryImpl := NewPPartnerRepositoryImpl(db)

	t.Run("should success find an object from p_partner table based on id.", func(t *testing.T) {
		wantId := uuid.NewString()
		wantName := "Sumber Jaya Maju"
		wantDesription := "Ini adalah deskripsi sumber jaya"
		wantCreatedAt := time.Now()
		wantUpdatedAt := sql.NullTime{}

		// init data
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.PPartner{
			Id:          wantId,
			Name:        wantName,
			Description: wantDesription,
			CreatedAt:   wantCreatedAt,
		})
		if errCreate != nil {
			t.Fatalf("errCreate: %v", errCreate)
		}

		pPartner, errFindOne := pPartnerRepositoryImpl.FindOne(context.TODO(), wantId)
		if errFindOne != nil {
			t.Fatalf("errFindOne: %v", errFindOne)
		}

		assert.Equal(t, wantId, pPartner.Id)
		assert.Equal(t, wantName, pPartner.Name)
		assert.Equal(t, wantDesription, pPartner.Description)
		assert.Equal(t, wantCreatedAt, wantCreatedAt)
		assert.Equal(t, wantUpdatedAt, pPartner.UpdatedAt)
	})

	t.Run("should failed find an object from p_partner table when id is empty.", func(t *testing.T) {

		gotPPartner, gotErrFindOne := pPartnerRepositoryImpl.FindOne(context.TODO(), "")

		assert.Error(t, gotErrFindOne)
		assert.Empty(t, gotPPartner.Id)
		assert.Empty(t, gotPPartner.Name)
		assert.Empty(t, gotPPartner.Description)
		assert.Equal(t, time.Time{}, gotPPartner.CreatedAt)
		assert.Equal(t, sql.NullTime{}, gotPPartner.UpdatedAt)
	})

	t.Run("should failed find an object from p_partner table when id is not exist on the records.", func(t *testing.T) {
		wantError := sql.ErrNoRows

		gotPPartner, gotErrFindOne := pPartnerRepositoryImpl.FindOne(context.TODO(), "00606566-8a44-4c58-8cac-ec09fd6595ac")

		assert.Equal(t, wantError, gotErrFindOne)
		assert.Equal(t, model.PPartner{}, gotPPartner)
	})

	t.Run("should failed find an object from p_partner table when id is not uuid.", func(t *testing.T) {
		gotPPartner, gotErrFindOne := pPartnerRepositoryImpl.FindOne(context.TODO(), "sdssdsdsd")

		assert.Error(t, gotErrFindOne)
		assert.Empty(t, gotPPartner.Id)
		assert.Empty(t, gotPPartner.Name)
		assert.Empty(t, gotPPartner.Description)
		assert.Equal(t, time.Time{}, gotPPartner.CreatedAt)
		assert.Equal(t, sql.NullTime{}, gotPPartner.UpdatedAt)
	})
}

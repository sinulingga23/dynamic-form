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
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Id:          uuid.NewString(),
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
		})

		assert.Nil(t, errCreate)
	})

	t.Run("should failed insert an object into p_partner table when id is empty.", func(t *testing.T) {
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
		})

		assert.Error(t, errCreate)
	})

	t.Run("should failed insert an object into p_partner table when id is not uuid.", func(t *testing.T) {
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Id:          "34739473949",
			Name:        fmt.Sprintf("New Partner: %v", time.Now().Unix()),
			Description: fmt.Sprintf("New Partner Description: %v", time.Now().Unix()),
			CreatedAt:   time.Now(),
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
		errCreate := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
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

func TestPPartnerRepositoryImpl_FindPPartnersByIds(t *testing.T) {
	db, errConnectDB := db.ConnectDB()
	if errConnectDB != nil {
		t.Fatalf("errConnectDB: %v", errConnectDB)
	}

	pPartnerRepositoryImpl := NewPPartnerRepositoryImpl(db)

	t.Run("should success find the ppartners based on ids", func(t *testing.T) {
		wantId1 := uuid.NewString()
		wantId2 := uuid.NewString()
		wantId3 := uuid.NewString()

		wantName1 := "New Partner 1"
		wantName2 := "New Partner 2"
		wantName3 := "New Partner 3"

		wantDescription1 := "New Description Partner 1"
		wantDescription2 := "New Description Partner 2"
		wantDescription3 := "New Description Partner 3"

		wantPPartners := []model.PPartner{
			model.PPartner{
				Id:          wantId1,
				Name:        wantName1,
				Description: wantDescription1,
			},
			model.PPartner{
				Id:          wantId2,
				Name:        wantName2,
				Description: wantDescription2,
			},
			model.PPartner{
				Id:          wantId3,
				Name:        wantName3,
				Description: wantDescription3,
			},
		}

		// init data
		errCreate1 := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Id:          wantId1,
			Name:        "New Partner 1",
			Description: "New Description Partner 1",
			CreatedAt:   time.Now(),
		})
		if errCreate1 != nil {
			t.Fatalf("errCreate1: %v", errCreate1)
		}

		errCreate2 := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Id:          wantId2,
			Name:        "New Partner 2",
			Description: "New Description Partner 2",
			CreatedAt:   time.Now(),
		})
		if errCreate2 != nil {
			t.Fatalf("errCreate2: %v", errCreate2)
		}

		errCreate3 := pPartnerRepositoryImpl.Create(context.TODO(), model.CreatePPartner{
			Id:          wantId3,
			Name:        "New Partner 3",
			Description: "New Description Partner 3",
			CreatedAt:   time.Now(),
		})
		if errCreate2 != nil {
			t.Fatalf("errCreate3: %v", errCreate3)
		}

		gotPPartners, errFindPPartnersByIds := pPartnerRepositoryImpl.FIndPPartnersByIds(context.TODO(), []string{wantId1, wantId2, wantId3})
		if errFindPPartnersByIds != nil {
			t.Fatalf("errFindPPartnersByIds: %v", errFindPPartnersByIds)
		}

		assert.Equal(t, len(gotPPartners), len([]string{wantId1, wantId2, wantId3}))
		for i := 0; i < len(gotPPartners); i++ {
			assert.Equal(t, wantPPartners[i].Id, gotPPartners[i].Id)
			assert.Equal(t, wantPPartners[i].Name, gotPPartners[i].Name)
			assert.Equal(t, wantPPartners[i].Description, gotPPartners[i].Description)
		}
	})

	t.Run("should failed find the ppartners when the ids is empty", func(t *testing.T) {

	})

	t.Run("should failed find the ppartners when the item for each ids not exist on the records", func(t *testing.T) {

	})

	t.Run("should failed find the ppartners when the item for each ids is not valid uuid", func(t *testing.T) {

	})

	t.Run("should failed find the ppartners when the item for each ids either not exists, not valid uuid, and empty", func(t *testing.T) {

	})
}

package tests

import (
"cloud.google.com/go/spanner"
"context"
"github.com/10antz-inc/cp-service-go/ssorm"
"testing"
)

func TestUpdateModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 23
	insert.FirstName = "first21"
	insert.LastName = "last21"
	
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Create(&insert, txn)
		_, err = db.Build().DeleteModel(&insert, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

func TestUpdateMap(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 23
	insert.FirstName = "first21"
	insert.LastName = "last21"

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Create(&insert, txn)
		params := map[string]interface{}{"SingerId": 23, "LastName": "test111"}
		err = db.Build().UpdateMap(&Singers{}, params, txn)
		_, err = db.Build().DeleteModel(&insert, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

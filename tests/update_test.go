package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/10antz-inc/ssorm"
	"testing"
)

func TestUpdateModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 12
	insert.FirstName = "updateModel"
	insert.LastName = "updateFlastNameModel"

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := db.Model(&insert).Update(ctx, txn)
		fmt.Println(count)
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
	insert.SingerId = 12
	insert.FirstName = "updateName"
	insert.LastName = "updateName"
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		params := map[string]interface{}{"LastName": "testMap"}
		count, err := db.Model(&insert).UpdateMap(ctx, txn, params)
		fmt.Println(count)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

func TestUpdateWhere(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()
	insert := Singers{}
	insert.SingerId = 12
	insert.FirstName = "updateName"
	insert.LastName = "updateName"
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		params := map[string]interface{}{"LastName": "testWhreMap"}
		count, err := db.Model(&insert).Where("SingerId > ?", 13).UpdateWhere(ctx, txn, params)
		fmt.Println(count)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

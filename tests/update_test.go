package tests

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/10antz-inc/ssorm"
)

func TestUpdateModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 12
	insert.FirstName = "updateModel"
	insert.LastName = spanner.NullString{StringVal: "last21",Valid: true}
	insert.TagIDs = []spanner.NullString{{StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c050", Valid: true}, {StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c051", Valid: true}}
	insert.Numbers = []int64{1, 2, 3}

	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.Model(&insert).Update(ctx, txn)
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
	update := Singers{}
	update.SingerId = 12
	update.FirstName = "updateName"
	update.LastName = spanner.NullString{StringVal: "last21",Valid: true}
	update.Numbers = []int64{10, 11, 12}
	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		update.LastName = spanner.NullString{StringVal: "last21",Valid: true}
		columns := []string{"LastName", "FirstName", "Numbers"}
		count, err := ssorm.Model(&update).UpdateColumns(ctx, txn, columns)
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
	insert.LastName = spanner.NullString{StringVal: "last21",Valid: true}
	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		params := map[string]interface{}{"TagIds": []spanner.NullString{{StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c052", Valid: true}, {StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c053", Valid: true}}}
		count, err := ssorm.Model(&insert).Where("SingerId > ?", 13).UpdateParams(ctx, txn, params)
		fmt.Println(count)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

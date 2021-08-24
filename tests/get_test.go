package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/10antz-inc/cp-service-go/ssorm"
	"testing"
)

func TestGetAllColumn(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	singer := Singers{}
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Where("SingerId in (?)", []int{12, 13, 14}).First(&singer, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when get singer, got %v", err)
	}
}

func TestGetColumn(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	singer := Singers{}
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Select([]string{"SingerId,FirstName"}).Where("SingerId in (?)", []int{12, 13, 14}).First(&singer, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when get singer, got %v", err)
	}
}

func TestGetAllColumnReadOnly(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	rtx := client.ReadOnlyTransaction()
	defer rtx.Close()

	singer := Singers{}
	db := ssorm.CreateDB()
	err := db.Build().Where("SingerId in (?)", []int{12, 13, 14}).First(&singer, ctx, rtx)
	
	if err != nil {
		t.Fatalf("Error happened when get singer, got %v", err)
	}
}

func TestGetColumnReadOnly(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	rtx := client.ReadOnlyTransaction()
	defer rtx.Close()

	singer := Singers{}
	db := ssorm.CreateDB()
	err := db.Build().Select([]string{"SingerId,FirstName"}).Where("SingerId in (?)", []int{12, 13, 14}).First(&singer, ctx, rtx)

	if err != nil {
		t.Fatalf("Error happened when get singer, got %v", err)
	}
}

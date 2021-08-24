package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/10antz-inc/cp-service-go/ssorm"
	"testing"
)

func TestSelectColumn(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var singers []*Singers
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Select([]string{"SingerId,FirstName"}).Where("SingerId in (?) and FirstName = ?", []int{12, 13, 14, 15}, "Dylan").Limit(1).Order("FirstName, LastName desc").Find(&singers, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

func TestSelectAllColumn(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var singers []*Singers
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Build().Where("SingerId in (?)", []int{12, 13, 14}).First(singers, ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

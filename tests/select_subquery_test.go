package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/10antz-inc/ssorm"
	"testing"
)

func TestSubQueryFirst(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	subSingers := Singer{}

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Model(&subSingers).TableName("Singers").AddSub(Albums{}, "SingerId = ?", 12).AddSub(Concerts{}, "SingerId = ?", 12).First(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}


func TestSubQueryFind(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var subSingers []*Singer

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Model(&subSingers).Where("SingerId > ?", 12).TableName("Singers").AddSub(Albums{}, "").AddSub(Concerts{}, "SingerId > ?", 12).Find(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

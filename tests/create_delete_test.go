package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/10antz-inc/ssorm"
	"testing"
)

func TestCreateDeleteModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 23
	insert.FirstName = "first21"
	insert.LastName = "last21"
	var singers []*Singers

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Model(&singers).Find(ctx, txn)
		_, err = db.Model(&insert).Insert(ctx, txn)

		_, err = db.Model(&insert).Update(ctx, txn)
		err = db.Model(&singers).Find(ctx, txn)

		_, err = db.Model(&insert).Where("SingerId = ?", 23).DeleteWhere(ctx, txn)

		err = db.Model(&singers).Find(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

func TestCreateDeleteWhere(t *testing.T) {
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
		count, err := db.Model(&insert).Insert(ctx, txn)
		fmt.Println(count)
		if err != nil {
			t.Fatalf("Error happened when create singer, got %v", err)
		}
		_, err = db.Model(&insert).Where("SingerId = ?", 23).DeleteWhere(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

func TestDelete(t *testing.T) {
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
		count, err := db.Model(&insert).DeleteModel(ctx, txn)
		fmt.Println(count)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

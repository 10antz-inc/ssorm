package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/10antz-inc/ssorm"
	"testing"
	"time"
)

func TestInsertDeleteModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 25
	insert.FirstName = "first21"
	insert.LastName = "last21"
	insert.TestTime = spanner.NullTime{time.Now(), true}
	var singers []*Singers

	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Model(&singers).Find(ctx, txn)
		err = db.Model(&singers).Find(ctx, txn)
		_, err = db.Model(&insert).Insert(ctx, txn)
		err = db.Model(&singers).Find(ctx, txn)
		_, err = db.Model(&insert).Update(ctx, txn)
		err = db.Model(&singers).Find(ctx, txn)

		_, err = db.Model(&insert).DeleteModel(ctx, txn)

		err = db.Model(&insert).First(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.Model(&singers).Find(ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

func TestInsertDeleteWhere(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 23
	insert.FirstName = "first21"
	insert.LastName = "last21"
	//insert.TestTime = time.Now()
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

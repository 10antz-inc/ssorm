package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/10antz-inc/ssorm"
)

func TestInsertDeleteModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 25
	//insert.FirstName = "first21"
	insert.LastName = spanner.NullString{StringVal: "last21", Valid: true}
	insert.TagIDs = []spanner.NullString{{StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c050", Valid: true}, {StringVal: "a3eb54bd-0138-4c22-b858-41bbefc5c051", Valid: true}}
	insert.Numbers = []int64{1, 2, 3}

	var singers []*Singers

	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.Model(&singers).Find(ctx, txn)
		err = ssorm.Model(&singers).Find(ctx, txn)
		_, err = ssorm.Model(&insert).Insert(ctx, txn)
		err = ssorm.Model(&singers).Find(ctx, txn)
		insert.TestTime = spanner.NullTime{Time: time.Now(), Valid: true}
		_, err = ssorm.Model(&insert).Update(ctx, txn)
		err = ssorm.Model(&singers).Find(ctx, txn)

		_, err = ssorm.Model(&insert).DeleteModel(ctx, txn)

		err = ssorm.Model(&insert).First(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		return err
	})

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.Model(&singers).Find(ctx, txn)
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
	//insert.FirstName = "first21"
	insert.LastName = spanner.NullString{StringVal: "last21", Valid: true}
	//insert.TestTime = time.Now()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.SimpleQueryWrite(ctx, txn, "Delete From Singers where SingerId = @singerId", map[string]interface{}{
			"singerId": 1,
		})
		count, err = ssorm.Model(&insert).Insert(ctx, txn)
		fmt.Println(count)
		if err != nil {
			t.Fatalf("Error happened when create singer, got %v", err)
		}
		_, err = ssorm.Model(&insert).Where("SingerId = ?", 23).DeleteWhere(ctx, txn)
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
	//insert.FirstName = "first21"
	insert.LastName = spanner.NullString{StringVal: "last21", Valid: true}

	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.Model(&insert).DeleteModel(ctx, txn)
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

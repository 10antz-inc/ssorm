package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/10antz-inc/ssorm"
	"testing"
)

func TestSimpleQueryRead(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var singers []*Singers
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.SimpleQueryRead(ctx, txn, "select * from singers limit 10", nil, &singers)
		return err
	})

	rtx := client.ReadOnlyTransaction()
	defer rtx.Close()
	singerRead := Singers{}
	err = ssorm.SimpleQueryRead(ctx, rtx, "select * from singers limit 1", nil, &singerRead)

	if err != nil {
		t.Fatalf("Error happened when count singers, got %v", err)
	}
}

func TestSimpleQueryWrite(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.SimpleQueryWrite(ctx, txn, "update singers set FirstName = \"TestSimpleQueryWrite\" where singerId =12", nil)
		fmt.Println(count)
		return err
	})

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.SimpleQueryWrite(ctx, txn, "INSERT Singers (SingerId, FirstName, LastName, TagIds, Numbers, UpdateTime, CreateTime) VALUES (1200, 'Melissa', 'Garcia', [\"a3eb54bd-0138-4c22-b858-41bbefc5c050\", \"a3eb54bd-0138-4c22-b858-41bbefc5c051\"], [1, 2], CURRENT_TIMESTAMP() ,CURRENT_TIMESTAMP());", nil)
		fmt.Println(count)
		return err
	})

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		count, err := ssorm.SimpleQueryWrite(ctx, txn, "DELETE FROM Singers where SingerId=1200", nil)
		fmt.Println(count)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when count singers, got %v", err)
	}
}

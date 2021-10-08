package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/10antz-inc/ssorm"
	"github.com/10antz-inc/ssorm/utils"
	"testing"
)

func TestSelectColumnReadWrite(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"

	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var singers []*Singers
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.Model(&singers).Select([]string{"SingerId,FirstName"}).Where("SingerId in ? and FirstName = ?", []int{12, 13, 14, 15}, "Dylan").Limit(1).Order("FirstName, LastName desc").Find(ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

func TestSelectAllColumnReadWrite(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"

	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	var singers []*Singers
	test := Singers{}
	utils.GetDeleteColumnName(&singers)
	utils.GetDeleteColumnName(&test)

	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.Model(&singers).Where("SingerId in ?", []int{12, 13, 14}).Find(ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

func TestSelectColumnReadOnly(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"

	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()
	rtx := client.ReadOnlyTransaction()
	defer rtx.Close()

	var singers []*Singers

	err := ssorm.Model(&singers).Select([]string{"SingerId,FirstName"}).Where("SingerId in ? and FirstName = ?", []int{12, 13, 14, 15}, "Dylan").Limit(1).Order("FirstName, LastName desc").Find(ctx, rtx)

	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

func TestSelectAllColumnReadOnly(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"

	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	rtx := client.ReadOnlyTransaction()
	defer rtx.Close()

	var singers []*Singers

	err := ssorm.Model(&singers).Where("SingerId in ?", []int{12, 13, 14}).Find(ctx, rtx)
	err = ssorm.SimpleQueryRead(ctx, rtx, "select * from Singers where singerId in UNNEST(@singerids)", map[string]interface{}{
		"singerids": []int{12, 13, 14},
	}, &singers)
	if err != nil {
		t.Fatalf("Error happened when search singers, got %v", err)
	}
}

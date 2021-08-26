package tests

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/10antz-inc/ssorm"
	"testing"
)

func TestSoftDeleteModel(t *testing.T) {
	url := "projects/spanner-emulator/instances/test/databases/test"
	ctx := context.Background()

	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	insert := Singers{}
	insert.SingerId = 12
	insert.FirstName = "updateModel"
	insert.LastName = "updateFlastNameModel"

	var singers []*Singers
	singer := Singers{}
	var count int64
	var subSingers []*Singer
	db := ssorm.CreateDB()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := db.SoftDeleteModel(&singers).Find(ctx, txn)
		err = db.SoftDeleteModel(&subSingers).Where("SingerId > ?", 12).TableName("Singers").AddSub(Albums{}, "").AddSub(Concerts{}, "SingerId > ?", 12).Find(ctx, txn)
		err = db.SoftDeleteModel(&singers).Where("SingerId = 13").Find(ctx, txn)
		err = db.SoftDeleteModel(insert).Where("SingerId in (?)", []int{12, 13, 14, 15}).Count(ctx, txn, &count)

		insert.SingerId = 100
		_, err = db.SoftDeleteModel(&insert).Insert(ctx, txn)
		_, err = db.SoftDeleteModel(&insert).Update(ctx, txn)
		_, err = db.SoftDeleteModel(&insert).DeleteModel(ctx, txn)
		err = db.SoftDeleteModel(&singer).Where("SingerId = ?", 25).First(ctx, txn)
		_, err = db.SoftDeleteModel(&insert).Where("SingerId = ?", 25).DeleteWhere(ctx, txn)
		_, err = db.Model(&insert).DeleteModel(ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}
		fmt.Println(count)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

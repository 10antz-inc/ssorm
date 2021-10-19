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
	//insert.FirstName = "updateModel"
	insert.LastName =spanner.NullString{StringVal: "updateFlastNameModel",Valid: true}

	var singers []*Singers
	singer := Singers{}
	var count int64
	var subSingers []*Singer
	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := ssorm.SoftDeleteModel(&singers).Find(ctx, txn)
		err = ssorm.SoftDeleteModel(&subSingers).Where("SingerId > ?", 12).TableName("Singers").AddSub(Albums{}, "").AddSub(Concerts{}, "SingerId > ?", 12).Find(ctx, txn)
		err = ssorm.SoftDeleteModel(&singers).Where("SingerId = 13").Find(ctx, txn)
		err = ssorm.SoftDeleteModel(insert).Where("SingerId in ?", []int{12, 13, 14, 15}).Count(ctx, txn, &count)

		insert.SingerId = 100
		_, err = ssorm.SoftDeleteModel(&insert).Insert(ctx, txn)
		_, err = ssorm.SoftDeleteModel(&insert).Update(ctx, txn)
		_, err = ssorm.SoftDeleteModel(&insert).DeleteModel(ctx, txn)
		err = ssorm.SoftDeleteModel(&singer).Where("SingerId = ?", 25).First(ctx, txn)
		_, err = ssorm.SoftDeleteModel(&insert).Where("SingerId = ?", 25).DeleteWhere(ctx, txn)
		_, err = ssorm.Model(&insert).DeleteModel(ctx, txn)
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

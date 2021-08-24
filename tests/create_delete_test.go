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
		err := db.Build().Find(&singers, ctx, txn)
		_, err = db.Build().Insert(&insert, ctx, txn)
		insert.SingerId = 12
		insert.LastName = "1322"
		count, err := db.Build().Update(&insert, ctx, txn)
		insert.SingerId = 23
		count, err = db.Build().Update(&insert, ctx, txn)
		params := map[string]interface{}{"LastName": "testMap"}
		count, err = db.Build().UpdateMap(&insert, ctx, params, txn)
		fmt.Println(count)
		insert.LastName = "1322111111111"

		err = db.Build().Find(&singers, ctx, txn)

		var rowCount int64
		fmt.Println(rowCount)
		//err = db.Build().InsertOrUpdate(&insert, txn)

		err = db.Build().Find(&singers, ctx, txn)
		//stmt = spanner.Statement{
		//	SQL: `DELETE FROM Singers where SingerId = 24`,
		//}
		//rowCount, err = txn.Update(ctx, stmt)
		//fmt.Println(rowCount)

		_, err = db.Build().Model(&insert).Where("SingerId = ?", 23).DeleteWhere(ctx, txn)

		err = db.Build().Find(&singers, ctx, txn)
		if err != nil {
			t.Fatalf("Error happened when delete singer, got %v", err)
		}

		//stmt = spanner.Statement{
		//	SQL: `DELETE FROM Singers where SingerId = 23`,
		//}
		//rowCount, err = txn.Update(ctx, stmt)
		//fmt.Println(rowCount)

		return err
	})

	if err != nil {
		t.Fatalf("Error happened when create singer, got %v", err)
	}
}

//
//func TestCreateDeleteWhere(t *testing.T) {
//	url := "projects/spanner-emulator/instances/test/databases/test"
//	ctx := context.Background()
//
//	client, _ := spanner.NewClient(ctx, url)
//	defer client.Close()
//
//
//
//	insert := Singers{}
//	insert.SingerId = 23
//	insert.FirstName = "first21"
//	insert.LastName = "last21"
//
//	db := ssorm.CreateDB()
//	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
//		err := db.Build().Insert(&insert, txn)
//		if err != nil {
//			t.Fatalf("Error happened when create singer, got %v", err)
//		}
//		_, err = db.Build().Model(&insert).Where("SingerId = ?",23).DeleteWhere( ctx, txn)
//		if err != nil {
//			t.Fatalf("Error happened when delete singer, got %v", err)
//		}
//
//		return err
//	})
//
//	if err != nil {
//		t.Fatalf("Error happened when create singer, got %v", err)
//	}
//}
//
//func TestDelete(t *testing.T) {
//	url := "projects/spanner-emulator/instances/test/databases/test"
//	ctx := context.Background()
//
//	client, _ := spanner.NewClient(ctx, url)
//	defer client.Close()
//
//
//
//	insert := Singers{}
//	insert.SingerId = 23
//	insert.FirstName = "first21"
//	insert.LastName = "last21"
//
//	db := ssorm.CreateDB()
//	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
//		count, err := db.Build().DeleteModel(&insert, ctx, txn)
//		fmt.Println(count)
//		if err != nil {
//			t.Fatalf("Error happened when delete singer, got %v", err)
//		}
//
//		return err
//	})
//
//	if err != nil {
//		t.Fatalf("Error happened when create singer, got %v", err)
//	}
//}

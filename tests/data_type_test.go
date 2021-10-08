package tests

import (
	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"context"
	"github.com/10antz-inc/ssorm"
	"testing"
	"time"
)

func TestDataTypeBool(t *testing.T) {
	url := "projects/spanner-emulator/instances/dev/databases/kagura"
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, url)
	defer client.Close()

	dataType := DataTypes{}
	dataType.DataTypesId = 3
	dataType.FirstName = "first21"
	dataType.TestTime = spanner.NullTime{Time: time.Now(), Valid: true}
	dataType.ArrayString = []spanner.NullString{{StringVal: "arr_str_1", Valid: true}, {StringVal: "arr_str_2", Valid: true}}
	dataType.ArrayInt64 = []int64{1, 2, 3}
	dataType.BoolValue = true
	dataType.FloatValue = 3.003
	dataType.ArrayFloat64 = []float64{1.01, 2.02, 3.03}
	dataType.DateValue, _ = civil.ParseDate("2021-09-30")
	//
	getDataType := DataTypes{}
	getDataTypeWhere := DataTypes{}
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		_, err := ssorm.Model(&dataType).Insert(ctx, txn)
		err = ssorm.Model(&getDataType).Where("DataTypesId = ?", 3).First(ctx, txn)
		getDataType.DateValue, _ = civil.ParseDate("2021-09-28")
		_, err = ssorm.Model(&getDataType).Where("DataTypesId = ?", 3).Update(ctx, txn)
		err = ssorm.Model(&getDataTypeWhere).Where("DateValue = ?", "2021-09-28").First(ctx, txn)
		ssorm.Model(&dataType).DeleteModel(ctx, txn)
		return err
	})

	if err != nil {
		t.Fatalf("Error happened when count singers, got %v", err)
	}
}

package task

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"learn_gy/util"
	"log"
	"time"
)

type Student struct {
	gorm.Model
	Name  string `gorm:"column:name"`
	Age   int    `gorm:"column:age"`
	Grade string `gorm:"column:grade"`
}

type Account struct {
	gorm.Model
	Balance decimal.Decimal `gorm:"type:decimal(20,2);not null;comment:交易金额（元）"`
}

type Transaction struct {
	gorm.Model
	FromAccountID string          `gorm:"column:from_account_id "` //源账户ID
	ToAccountID   string          `gorm:"column:to_account_id  "`  //目标账户ID
	Amount        decimal.Decimal `gorm:"type:decimal(20,2);not null;comment:交易金额（元）"`
}

// Crud 基本CRUD操作
func CrudTest(db *gorm.DB) {
	//db.AutoMigrate(&Student{})
	//db.Create(&Student{Name: "王五", Grade: "三年级", Age: 18})

	// query 操作
	var students []Student
	db.Where("age > ?", 18).Find(&students)
	fmt.Println(students) //
	json, err := json.Marshal(students)
	if err != nil {
		log.Fatal("Failed to marshal student:", err)
	}
	fmt.Println(string(json))
	// update 操作
	errorUpdate := util.UpdateByModel(db, &Student{}, map[string]interface{}{"name": "张三"}, map[string]interface{}{"grade": "六年级"})
	if errorUpdate != nil {
		log.Fatal("Failed to marshal student:", errorUpdate)
	}
	// delete 操作
	errDel := util.DeleteByModel(db, &Student{}, "age < ?", 15)
	if errDel != nil {
		log.Fatal("Failed to delete student:", errDel)
	}
}

func TransactionTest(db *gorm.DB) {
	//db.AutoMigrate(&Account{})
	//accountA := &Account{Balance: decimal.NewFromInt(100)}
	//accountB := &Account{Balance: decimal.NewFromInt(100)}
	//util.CreateByModel(db, &accountA)
	//util.CreateByModel(db, &accountB)
	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("SET innodb_lock_wait_timeout = 1")
		var accountA, accountB Account
		if err := tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&accountA, 1).Error; err != nil {
			return fmt.Errorf("获取账户A失败: %v", err)
		}
		if err := tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&accountB, 2).Error; err != nil {
			return fmt.Errorf("获取账户B失败: %v", err)
		}
		transferAmount := decimal.NewFromInt(100)
		if accountA.Balance.Cmp(transferAmount) < 0 {
			return fmt.Errorf("账户A余额不足，请充值!")
		}
		accountA.Balance = accountA.Balance.Sub(transferAmount)
		accountB.Balance = accountB.Balance.Add(transferAmount)
		txA := tx.Debug().Model(&accountA).Updates(Account{Balance: accountA.Balance})
		if txA.Error != nil {
			return fmt.Errorf("账户A更新失败！: %v", txA.Error)
		}
		txB := tx.Debug().Model(&accountB).Updates(Account{Balance: accountB.Balance})
		if txB.Error != nil {
			return fmt.Errorf("账户B更新失败！: %v", txA.Error)
		}
		time.Sleep(2000 * time.Millisecond)
		return nil
	})

	if err != nil {
		fmt.Println("Transaction failed:", err)
	} else {
		fmt.Println("Transaction completed successfully.")
	}
}

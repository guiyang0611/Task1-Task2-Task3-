package task

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// SqlXTest1 题目1：模型定义
func SqlXTest1(db *sqlx.DB) {
	var employees []Employee
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	for _, emp := range employees {
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪资: %.2f\n", emp.ID, emp.Name, emp.Department, emp.Salary)
	}

	//薪资最高员工
	var employee Employee
	err = db.Get(&employee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	fmt.Println("薪资最高员工: ")
	fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪资: %.2f\n", employee.ID, employee.Name, employee.Department, employee.Salary)

}

// SqlXTest2 题目2：实现类型安全映射
func SqlXTest2(db *sqlx.DB) {
	var books []Book
	err := db.Select(&books, "SELECT  id, title, author, price FROM books WHERE price > ?", 50)
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	fmt.Println("价格大于 50 的书籍: ")
	for _, book := range books {
		fmt.Printf("ID: %d, 标题: %s, 作者: %s, 价格: %.2f\n", book.ID, book.Title, book.Author, book.Price)
	}

}

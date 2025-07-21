package task

import (
	"learn_gy/config"
)

func TestTask3() {
	//初始化数据库
	err := config.InitDb()
	if err != nil {
		panic(err)
	}
	//测试CRUD
	//CrudTest(config.GormDB)
	//// 测试事务
	//for i := 0; i < 100; i++ {
	//	go TransactionTest(config.GormDB)
	//}
	//time.Sleep(3 * time.Second)
	//// 题目1：模型定义
	//SqlXTest1(config.SqlxDB)
	////题目2：实现类型安全映射
	//SqlXTest2(config.SqlxDB)
	//博客系统
	TestBlog(config.GormDB)
}

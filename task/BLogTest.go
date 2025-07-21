package task

import (
	"fmt"
	"gorm.io/gorm"
	"learn_gy/util"
	"log"
	"strconv"
)

func TestBlog(db *gorm.DB) {
	//saveBLog(db)
	//queryBLog(db)
	delBLog(db, 25)
}

type CommentStatus int

const (
	CommentDisabled CommentStatus = iota
	CommentEnabled
)

func (s CommentStatus) Desc() string {
	return [...]string{"无评论", "有评论"}[s]
}

type User struct {
	gorm.Model
	Name      string
	Age       int
	PostCount int
	Posts     []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	Status   CommentStatus
	Comments []Comment `gorm:"foreignKey:PostID"`
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	UserID  uint
}

func (*User) TableName() string {
	return "users"
}

func (*Post) TableName() string {
	return "posts"
}

func (*Comment) TableName() string {
	return "comments"
}

type Result struct {
	Post
	CommentCount uint
}

func (post *Post) AfterCreate(tx *gorm.DB) (err error) {
	var user User
	tx.First(&user, post.UserID)
	user.PostCount++
	return tx.Save(&user).Error
}

func (comment *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var post Post
	tx.First(&post, comment.PostID)
	fmt.Println("开始更新文章评论状态----", post)
	var commentCount int64
	tx.Debug().Model(&Comment{}).Where(&Comment{PostID: comment.PostID}).Select("count(*) As comment_count").Count(&commentCount)
	fmt.Println("文章评论数量为：", commentCount)
	if commentCount > 0 {
		return nil
	}
	post.Status = CommentDisabled
	return tx.Save(&post).Error
}

func saveBLog(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})
	for i := 1; i <= 5; i++ {
		age := i * 10
		name := "张三" + strconv.Itoa(i)
		user := User{Name: name, Age: age}
		util.CreateByModel(db, &user)
		for j := 1; j <= 3; j++ {
			post := Post{Title: user.Name + "的文章标题" + strconv.Itoa(i*j),
				Content: user.Name + "的文章内容" + strconv.Itoa(i*j), UserID: user.ID, Status: CommentEnabled}
			util.CreateByModel(db, &post)
			for k := 1; k <= 5; k++ {
				comment := Comment{Content: post.Title + "的评论内容" + strconv.Itoa(i*j*k), PostID: post.ID, UserID: user.ID}
				util.CreateByModel(db, &comment)
			}
		}
	}
}

func queryBLog(db *gorm.DB) {
	var user User
	tx := db.Preload("Posts").Preload("Posts.Comments").First(&user, 1)
	if err := tx.Error; err != nil {
		log.Fatal("Failed to marshal user:", err)
	}
	fmt.Printf("用户: %s, 年龄: %d\n", user.Name, user.Age)
	for _, post := range user.Posts {
		fmt.Printf("  文章: %s\n", post.Title)
		for _, comment := range post.Comments {
			fmt.Printf("    评论: %s\n", comment.Content)
		}
	}
	var result Result
	err := db.Model(&Post{}).
		Select("posts.*, count(comments.id) AS comment_count").
		Joins("LEFT JOIN comments on posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		log.Fatalln("查询失败:", err)
	}
	fmt.Printf("最多评论的文章: %s, 评论数: %d\n", result.Title, result.CommentCount)
}

func delBLog(db *gorm.DB, commentId uint) {
	var comment Comment
	db.First(&comment, commentId)
	err := util.DeleteByModel(db, &comment, "id = ?", commentId)
	if err != nil {
		log.Fatalln("删除失败:", err)
	}
	fmt.Println("删除成功", comment)

	//查询文章的评论状态和数量
	var post Post
	db.Preload("Comments").First(&post, comment.PostID)
	fmt.Println("文章的评论状态:", post.Status.Desc(), ",文章的评论数量:", len(post.Comments))

}

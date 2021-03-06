package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

type Line struct {
	Id       int       `gorm:"size:11;primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name     string    `gorm:"size:255;DEFAULT NULL" json:"name"`
	Content  string    `gorm:"size:255;DEFAULT NULL" json:"content"`
	Ip       string    `gorm:"size:255;DEFAULT NULL" json:"ip"`
	tag      string    `gorm:"size:255;DEFAULT NULL" json:"tag"`
	CreateAt time.Time `json:"createat"`
}

func init() {
	var err error
	DB, err = gorm.Open("sqlite3", "./mtimeline.db")
	if err != nil {
		fmt.Printf("------- sqlite3 connect error %v", err)
	}
	if DB.Error != nil {
		fmt.Printf("------- sqlite3 database error %v", DB.Error)
	}
	DB.SingularTable(true)
	DB.AutoMigrate(&Line{})
}

// DB GET
func (This *Line) GetAllHosts() [][]string {
	var line []Line
	DB.Order("create_at desc").Find(&line)
	var lines [][]string
	for _, l := range line {
		var oneLine []string
		oneLine = append(oneLine, l.CreateAt.Format("2006-01-02 15:04:05"), l.Name, l.Content, l.Ip)
		lines = append(lines, oneLine)
	}
	return lines
}

// 新增页
func LinePost(c *gin.Context) {
	fmt.Println(c.ClientIP())
	var line Line

	line.Name = c.PostFormArray("name")[0]
	line.Content = c.PostFormArray("content")[0]
	line.Ip = c.ClientIP()
	line.CreateAt = time.Now()
	DB.Create(&line)
	c.Redirect(http.StatusMovedPermanently, "/")
}

// 展示页
func FrontPage(c *gin.Context) {
	var line Line
	lines := line.GetAllHosts()
	c.HTML(
		http.StatusOK,
		"frontpage.html",
		gin.H{
			"title":      "mtimeline",
			"linesValue": lines,
		},
	)

}

func main() {
	defer DB.Close()

	fmt.Println("----------->mtimeline!")

	gin.ForceConsoleColor()
	router := gin.Default()                                                        // Logger & Recovery
	router.LoadHTMLGlob("templates/*")                                             // html模板
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string { // 日志格式化
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.GET("/", FrontPage)
	router.POST("/line", LinePost)

	router.Run(":6969")
}

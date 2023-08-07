package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

const SUCCESS = 0
const FAILED = -1

// GetIndex
// @Tags 首页
// @Success 200 {string} HelloWorld!
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello golang web!",
		"code":    SUCCESS,
		"data":    nil,
	})
}

// GetUserList
// @Summary 用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code":    SUCCESS,
		"data":    data,
	})
}

// CreateUser
// @Summary 创建用户
// @Tags 用户模块
// @Param name formData string false "姓名"
// @Param password formData string true "密码"
// @Param repassword formData string true "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if (models.FindUserByName(user.Name).ID) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户名已存在！",
			"code":    FAILED,
			"data":    nil,
		})
		return
	}
	if password != repassword {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    FAILED,
			"message": "两次密码不一致!",
			"data":    nil,
		})
		return
	}
	user.PassWord = password
	models.CreateUser(user)
	fmt.Println("===========================", user.Name)
	c.JSON(http.StatusCreated, gin.H{
		"code":    SUCCESS,
		"message": "新建用户成功",
		"data":    nil,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @Param id formData string false "姓名"
// @Param name formData string false "姓名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param phone formData string true "手机号码"
// @Success 201 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Email = c.PostForm("email")
	user.Phone = c.PostForm("phone")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    FAILED,
			"message": err,
			"data":    nil,
		})
	} else {
		data := models.UpdateUser(user)
		fmt.Println("===============", user)
		fmt.Println("===============", data.Error)
		c.JSON(http.StatusCreated, gin.H{
			"code":    FAILED,
			"message": "修改成功!",
			"data":    nil,
		})
	}
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id query string true "用户ID"
// @Success 201 {string} json{"code","message"}
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	fmt.Println("==================", id)
	if models.FindUserById(id).ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    FAILED,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	models.DeleteUser(id, user)
	c.JSON(http.StatusCreated, gin.H{
		"code":    SUCCESS,
		"message": "删除成功",
		"data":    nil,
	})
}

// LoginUser
// @Summary 用户登录
// @Tags 用户模块
// @Param name formData string false "姓名"
// @Param password formData string true "密码"
// @Success 201 {string} json{"code","message"}
// @Router /user/loginUser [post]
func LoginUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user = models.FindUserByNameAndPwd(user.Name, user.PassWord)
	if user.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    FAILED,
			"message": "用户不存在",
			"data":    nil,
		})
	} else {
		user.Identity = utils.RandomString(10)
		data := models.UpdateUser(user)
		fmt.Println("===============", user)
		fmt.Println("===============", data.Error)
		c.JSON(http.StatusCreated, gin.H{
			"code":    SUCCESS,
			"message": "登录成功！",
			"data":    user,
		})
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	//ReadBufferSize:  1024,
	//WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("err===============", err)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println("Close===============", err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("MsgHandler=============", err)
		}
		tm := time.Now()
		fmt.Println("msg==============", msg)
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println("message-err======", err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

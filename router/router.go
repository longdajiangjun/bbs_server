package router

import (
	"bbs_server/controller/admin"
	"bbs_server/controller/feedback"
	"bbs_server/controller/index"
	"bbs_server/controller/post"
	"bbs_server/controller/user"
	"bbs_server/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init() {
	router := gin.Default()
	// CrossDomain跨域处理，options请求处理
	router.Use(middleware.CrossDomain())
	//router.Use(middleware.TlsHandler())
	// v1群组对任何人开放
	v1 := router.Group("/v1")
	{
		v1.POST("/login", user.Login)
		v1.POST("/register", user.Register)
		v1.GET("/index", index.GetInfo)
		v1.GET("/posts", post.GetPosts)
		v1.GET("/post", post.GetPost)
		v1.POST("/admin", admin.Login)
		v1.GET("/notices/get", user.GetNotice)
		v1.GET("/zhiding/get", user.GetZhiDing)
	}

	v2 := router.Group("/v2")
	// v2群组使用中间件AuthMiddleWare，需要token权限才能请求到
	v2.Use(middleware.AuthMiddleWare())
	{
		v2.POST("/publish", post.Publish)
		v2.POST("/isload", user.IsLoad)
		v2.POST("/reply1", post.Reply1)
		v2.POST("/reply2", post.Reply2)
		v2.POST("/signin", user.Signin)
		v2.DELETE("/delpost", post.DelPost)
		v2.POST("/support", post.Support)
		v2.POST("/cancel", post.Cancel)
		v2.POST("/complaint", feedback.Complaint)
		v2.POST("/user/myposts", user.GetMyPosts)
		v2.POST("gift/showGiftPack", user.ShowGiftPack)
		v2.POST("gift/getGiftPack", user.GetGiftPack)
	}
	// 管理员路由组（管理员请求）
	adminAPI := router.Group("/admin")
	adminAPI.Use(middleware.AuthAdmin())
	{
		adminAPI.POST("/count", admin.Count)
		adminAPI.POST("/search", admin.UserSearch)
		adminAPI.POST("/addBlackList", admin.AddBlackList)
		adminAPI.POST("/removeBlackList", admin.RemoveBlackList)
		adminAPI.POST("/getBlackList", admin.GetBlackList)
		adminAPI.POST("/getFeedList0", admin.GetFeedList0)
		adminAPI.POST("/getFeedList1", admin.GetFeedList1)
		adminAPI.POST("/delFeedBack", admin.DelFeedBack)
		adminAPI.POST("/agreeFeedBack", admin.AgreeFeedBack)
		adminAPI.POST("/notices/save", admin.AddNotice)
		adminAPI.POST("/zhiding", admin.ZhiDing)
		adminAPI.POST("/sendGiftPack", admin.SendGiftPack)
		adminAPI.POST("/getGiftPacks", admin.GetGiftPacks)
		adminAPI.POST("/delGiftPack", admin.DelGiftPack)
		adminAPI.POST("/countRedeemCodes", admin.CountRedeemCodes)
		adminAPI.POST("/addTopic", admin.AddTopic)
	}
	//router.RunTLS(":8000", "./ssl/1_www.yinghuo2018.com_bundle.crt", "./ssl/2_www.yinghuo2018.com.key")
	router.Run(":8000")
}
                                                         
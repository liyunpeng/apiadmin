package routes

import (
	"apiadmin/controllers"
	"apiadmin/middleware"
	"apiadmin/models"
	"github.com/kataras/iris/v12"
)

func Register(api *iris.Application) {
	main := api.Party("/", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	{
		api.Get("/", func(ctx iris.Context) { // 首页模块
			_ = ctx.View("index.html")
		})

		api.Get("/apiDoc", func(ctx iris.Context) {
			_ = ctx.View("apiDoc/index.html")
		})

		v1 := main.Party("/v1")
		{
			v1.Post("/admin/login", controllers.UserLogin)
			v1.PartyFunc("/admin", func(admin iris.Party) {

				casbinMiddleware := middleware.New(models.Enforcer)                  //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
				admin.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				admin.Get("/logout", controllers.UserLogout).Name = "退出"
				//v1.Use(irisyaag.New())
				admin.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers).Name = "用户列表123"
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"
					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"
					users.Get("/profile", controllers.GetProfile).Name = "个人信息"
				})
				admin.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
					roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
					roles.Post("/", controllers.CreateRole).Name = "创建角色"
					roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
					roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
				})
				admin.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions).Name = "权限列表"
					permissions.Get("/{id:uint}", controllers.GetPermission).Name = "权限详情"
					permissions.Post("/import", controllers.ImportPermission).Name = "导入权限"
					//permissions.Post("/", controllers.CreatePermission).Name = "创建权限"
					//permissions.Put("/{id:uint}", controllers.UpdatePermission).Name = "编辑权限"
					//permissions.Delete("/{id:uint}", controllers.DeletePermission).Name = "删除权限"
				})
			})
		}

		api.Any("/payload", controllers.Payload)
	}

}

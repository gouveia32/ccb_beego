package routers

import (
	"ccb_beego/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//websocket
	beego.Router("/websocketwidget/index", &controllers.WebsocketWidgetController{}, "*:Index")
	beego.Router("/websocketwidget/ws", &controllers.WebsocketWidgetController{}, "Get:Get")

	//Roteamento da função do usuário
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//roteamento de recursos
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")

	//Ordem de alteração rápida
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//Painel de seleção geral
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	beego.Router("/resource/chooseIcon", &controllers.ResourceController{}, "Get:ChooseIcon")

	//Lista de menus (incluindo áreas) que o usuário tem permissão para gerenciar
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//Roteamento de usuário em segundo plano
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "POST:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get,Post:Edit")
	beego.Router("/backenduser/delete", &controllers.BackendUserController{}, "Post:Delete")

	//Centro do usuário em segundo plano
	beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/usercenter/passwordsave", &controllers.UserCenterController{}, "Post:PasswordSave")

	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/index2", &controllers.HomeController{}, "*:Index2")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/home/weather", &controllers.HomeController{}, "*:GetWeather")

	//Cliente
	beego.Router("/cliente/index", &controllers.ClienteController{}, "*:Index")
	beego.Router("/cliente/datagrid", &controllers.ClienteController{}, "Get,Post:DataGrid")
	beego.Router("/cliente/datalist", &controllers.ClienteController{}, "Post:DataList")
	beego.Router("/cliente/edit/?:id", &controllers.ClienteController{}, "Get,Post:Edit")
	beego.Router("/cliente/delete", &controllers.ClienteController{}, "Post:Delete")

	//Linha
	beego.Router("/linha/index", &controllers.LinhaController{}, "*:Index")
	beego.Router("/linha/datagrid", &controllers.LinhaController{}, "Get,Post:DataGrid")
	beego.Router("/linha/datalist", &controllers.LinhaController{}, "Post:DataList")
	beego.Router("/linha/edit/?:id", &controllers.LinhaController{}, "Get,Post:Edit")
	beego.Router("/linha/delete", &controllers.LinhaController{}, "Post:Delete")

	beego.Router("/", &controllers.HomeController{}, "*:Index")
}

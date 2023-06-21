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

	//Fornecedor
	beego.Router("/fornecedor/index", &controllers.FornecedorController{}, "*:Index")
	beego.Router("/fornecedor/datagrid", &controllers.FornecedorController{}, "Get,Post:DataGrid")
	beego.Router("/fornecedor/datalist", &controllers.FornecedorController{}, "Post:DataList")
	beego.Router("/fornecedor/edit/?:id", &controllers.FornecedorController{}, "Get,Post:Edit")
	beego.Router("/fornecedor/delete", &controllers.FornecedorController{}, "Post:Delete")

	//Linha
	beego.Router("/linha/index", &controllers.LinhaController{}, "*:Index")
	beego.Router("/linha/datagrid", &controllers.LinhaController{}, "Get,Post:DataGrid")
	beego.Router("/linha/datalist", &controllers.LinhaController{}, "Post:DataList")
	beego.Router("/linha/edit/?:codigo", &controllers.LinhaController{}, "Get,Post:Edit")
	beego.Router("/linha/delete", &controllers.LinhaController{}, "Post:Delete")

	//Bordado
	beego.Router("/bordado/index", &controllers.BordadoController{}, "*:Index")
	beego.Router("/bordado/datagrid", &controllers.BordadoController{}, "Get,Post:DataGrid")
	beego.Router("/bordado/datalist", &controllers.BordadoController{}, "Post:DataList")
	beego.Router("/bordado/edit/?:id", &controllers.BordadoController{}, "Get,Post:Edit")
	beego.Router("/bordado/delete", &controllers.BordadoController{}, "Post:Delete")
	beego.Router("/bordado/lerdst", &controllers.BordadoController{}, "Post:LerDst")

	//Roteamento do grupo
	beego.Router("/grupo/index", &controllers.GrupoController{}, "*:Index")
	beego.Router("/grupo/datagrid", &controllers.GrupoController{}, "Get,Post:DataGrid")
	beego.Router("/grupo/edit/?:id", &controllers.GrupoController{}, "Get,Post:Edit")
	beego.Router("/grupo/delete", &controllers.GrupoController{}, "Post:Delete")
	beego.Router("/grupo/datalist", &controllers.GrupoController{}, "Post:DataList")

	//Roteamento do catalogo
	beego.Router("/catalogo/index", &controllers.CatalogoController{}, "*:Index")
	beego.Router("/catalogo/datagrid", &controllers.CatalogoController{}, "Get,Post:DataGrid")
	beego.Router("/catalogo/edit/?:id", &controllers.CatalogoController{}, "Get,Post:Edit")
	beego.Router("/catalogo/delete", &controllers.CatalogoController{}, "Post:Delete")
	beego.Router("/catalogo/datalist", &controllers.CatalogoController{}, "Post:DataList")
	beego.Router("/catalogo/allocate", &controllers.CatalogoController{}, "Post:Allocate")
	beego.Router("/catalogo/updateseq", &controllers.CatalogoController{}, "Post:UpdateSeq")

	beego.Router("/", &controllers.HomeController{}, "*:Index")
}

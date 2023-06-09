
package controllers

import (
	"encoding/json"

	"ccb_beego/enums"
	"ccb_beego/models"

	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

// GrupoController 角色管理
type GrupoController struct {
	BaseController
}

// Prepare 参考beego官方文档说明
func (c *GrupoController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

// Index 角色管理首页
func (c *GrupoController) Index() {
	c.Data["pageTitle"] = "Ger. de Funções"

	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false

	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "grupo/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "grupo/index_footerjs.html"

	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("GrupoController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("GrupoController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("GrupoController", "Allocate")
}

// DataGrid 角色管理首页 表格获取数据
func (c *GrupoController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.GrupoQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.GrupoPageList(&params)

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	c.Data["json"] = result
	c.ServeJSON()
}

// DataList 角色列表
func (c *GrupoController) DataList() {
	var params = models.GrupoQueryParam{}

	//获取数据列表和总数
	data := models.GrupoDataList(&params)

	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *GrupoController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}

	Id, _ := c.GetInt(":id", 0)
	m := models.Grupo{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("Os dados são inválidos, atualize e tente novamente")
		}
	}
	c.Data["m"] = m
	c.setTpl("grupo/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "grupo/edit_footerjs.html"
}

// Save 添加、编辑页面 保存
func (c *GrupoController) Save() {
	var err error
	m := models.Grupo{}

	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter dados", m.Id)
	}

	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "Gravação com sucesso", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "Atualizado", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "Falha ao modificar", m.Id)
		}
	}
}

// Delete 批量删除
func (c *GrupoController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))

	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.GrupoBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("o ítem %d foi excluído", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", 0)
	}
}

// Allocate
func (c *GrupoController) Allocate() {
	grupoId, _ := c.GetInt("id", 0)
	strs := c.GetString("ids")

	o := orm.NewOrm()
	m := models.Grupo{Id: grupoId}
	if err := o.Read(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Os dados são inválidos, atualize e tente novamente", "")
	}

	var relations []models.GrupoBordadoRel
	for _, str := range strings.Split(strs, ",") {
		if _, err := strconv.Atoi(str); err == nil {
			relation := models.GrupoBordadoRel{Grupo: &m}
			relations = append(relations, relation)
		}
	}

	if len(relations) > 0 {
		//批量添加
		if _, err := o.InsertMulti(len(relations), relations); err == nil {
			c.jsonResult(enums.JRCodeSucc, "Salvo com sucesso", "")
		}
	}

	c.jsonResult(0, "Falha ao Salvar", "")
}

func (c *GrupoController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.GrupoOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "Dados inválidos selecionados", 0)
	}

	value, _ := c.GetInt("value", 0)
	oM.Seq = value

	o := orm.NewOrm()
	if _, err := o.Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "Modificado com sucesso", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha ao Modificar", oM.Id)
	}
}

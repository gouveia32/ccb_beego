// Copyright 2018 ccb_beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"ccb_beego/enums"
	"ccb_beego/models"

	"io/ioutil"

	"github.com/beego/beego/v2/client/orm"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//O controle de permissão é comentado aqui, então a verificação de login é necessária aqui
	c.checkLogin()
}

func (c *ResourceController) Index() {
	c.Data["pageTitle"] = "Gestão de Recursos"

	//Precisa de controle de permissão
	c.checkAuthor()

	//Ative um item no menu à esquerda da página
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/index_footerjs.html"

	//Controle de permissão do botão na página
	c.Data["canEdit"] = c.checkActionAuthor("ResourceController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("ResourceController", "Delete")
}

// TreeGrid 获取所有资源的列表
func (c *ResourceController) TreeGrid() {
	tree := models.ResourceTreeGrid()

	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// UserMenuTree 获取用户有权管理的菜单、区域列表
func (c *ResourceController) UserMenuTree() {
	userid := c.curUser.Id

	//获取用户有权管理的菜单列表（包括区域）
	tree := models.ResourceTreeGridByUserId(userid, 1)

	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *ResourceController) ParentTreeGrid() {
	Id, _ := c.GetInt("id", 0)
	tree := models.ResourceTreeGrid4Parent(Id)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// UrlFor2LinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2LinkOne(urlfor string) string {
	if len(urlfor) == 0 {
		return ""
	}

	// ResourceController.Edit,:id,1
	strs := strings.Split(urlfor, ",")
	if len(strs) == 1 {
		return c.URLFor(strs[0])
	} else if len(strs) > 1 {
		var values []interface{}
		for _, val := range strs[1:] {
			values = append(values, val)
		}

		return c.URLFor(strs[0], values...)
	}
	return ""
}

// UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *ResourceController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

// Edit 资源编辑页面
func (c *ResourceController) Edit() {
	//需要权限控制
	c.checkAuthor()

	//如果是POST请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}

	Id, _ := c.GetInt(":id", 0)
	m := &models.Resource{}
	var err error
	if Id == 0 {
		m.Seq = 100
	} else {
		m, err = models.ResourceOne(Id)
		if err != nil {
			c.pageError("Os dados são inválidos, atualize e tente novamente")
		}
	}

	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}

	//获取可以成为当前节点的父节点的列表
	c.Data["parents"] = models.ResourceTreeGrid4Parent(Id)

	//转换地址
	m.LinkUrl = c.UrlFor2LinkOne(m.UrlFor)

	c.Data["m"] = m
	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}

	c.setTpl("resource/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "resource/edit_footerjs.html"
}

// Save 资源添加编辑 保存
func (c *ResourceController) Save() {
	var err error
	o := orm.NewOrm()
	parent := &models.Resource{}
	m := models.Resource{}
	parentId, _ := c.GetInt("Parent", 0)

	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter dados", m.Id)
	}

	//获取父节点
	if parentId > 0 {
		parent, err = models.ResourceOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.jsonResult(enums.JRCodeFailed, "Nó pai inválido", "")
		}
	}

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

// Delete 删除
func (c *ResourceController) Delete() {
	//需要权限控制
	c.checkAuthor()
	Id, _ := c.GetInt("Id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "Dados inválidos selecionados", 0)
	}

	query := orm.NewOrm().QueryTable(models.ResourceTBName())
	if _, err := query.Filter("id", Id).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("excluído com sucesso"), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", 0)
	}
}

// Select 通用选择面板
func (c *ResourceController) Select() {
	//获取调用者的类别 1表示 角色
	desttype, _ := c.GetInt("desttype", 0)

	//获取调用者的值
	destval, _ := c.GetInt("destval", 0)

	//返回的资源列表
	var selectedIds []string
	o := orm.NewOrm()
	if desttype > 0 && destval > 0 {
		//如果都大于0,则获取已选择的值，例如：角色，就是获取某个角色已关联的资源列表
		switch desttype {
		case 1:
			{
				role := models.Role{Id: destval}
				o.LoadRelated(&role, "RoleResourceRel")
				for _, item := range role.RoleResourceRel {
					selectedIds = append(selectedIds, strconv.Itoa(item.Resource.Id))
				}
			}
		}
	}

	c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	c.setTpl("resource/select.html", "shared/layout_pullbox.html")

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/select_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/select_footerjs.html"
}

func (c *ResourceController) ChooseIcon() {
	filename := "static/plugins/font-awesome/less/variables.less"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}

	var iconlist []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "@fa-var-") {
			tempStr := line[8:]
			idx := strings.Index(tempStr, ":")
			icon := tempStr[:idx]
			iconlist = append(iconlist, icon)
		}
	}
	c.Data["Iconlist"] = iconlist
	c.setTpl("resource/chooseicon.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "resource/chooseicon_headcssjs.html"
	c.LayoutSections["footerjs"] = "resource/chooseicon_footerjs.html"
}

// CheckUrlFor 填写UrlFor时进行验证
func (c *ResourceController) CheckUrlFor() {
	urlfor := c.GetString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.jsonResult(enums.JRCodeSucc, "analisado com sucesso", link)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Analisando Falha", link)
	}
}

func (c *ResourceController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)

	oM, err := models.ResourceOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "Dados inválidos selecionados", 0)
	}

	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	if _, err := orm.NewOrm().Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "Modificado com sucesso", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Modificar Falha", oM.Id)
	}
}

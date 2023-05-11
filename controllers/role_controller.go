// Copyright 2018 gardens Author. All Rights Reserved.
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
	"encoding/json"

	"gardens/enums"
	"gardens/models"

	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

// RoleController 角色管理
type RoleController struct {
	BaseController
}

// Prepare 参考beego官方文档说明
func (c *RoleController) Prepare() {
	//先执行
	c.BaseController.Prepare()

	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")

	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

// Index 角色管理首页
func (c *RoleController) Index() {
	c.Data["pageTitle"] = "Ger. de Funções"

	//是否显示更多查询条件的按钮
	c.Data["showMoreQuery"] = false

	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)

	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "role/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "role/index_footerjs.html"

	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("RoleController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("RoleController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("RoleController", "Allocate")
}

// DataGrid 角色管理首页 表格获取数据
func (c *RoleController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.RoleQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//获取数据列表和总数
	data, total := models.RolePageList(&params)

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	c.Data["json"] = result
	c.ServeJSON()
}

// DataList 角色列表
func (c *RoleController) DataList() {
	var params = models.RoleQueryParam{}

	//获取数据列表和总数
	data := models.RoleDataList(&params)

	//定义返回的数据结构
	c.jsonResult(enums.JRCodeSucc, "", data)
}

// Edit 添加、编辑角色界面
func (c *RoleController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}

	Id, _ := c.GetInt(":id", 0)
	m := models.Role{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("Os dados são inválidos, atualize e tente novamente")
		}
	}
	c.Data["m"] = m
	c.setTpl("role/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "role/edit_footerjs.html"
}

// Save 添加、编辑页面 保存
func (c *RoleController) Save() {
	var err error
	m := models.Role{}

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
func (c *RoleController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))

	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.RoleBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("o ítem %d foi excluído", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", 0)
	}
}

// Allocate 给角色分配资源界面
func (c *RoleController) Allocate() {
	roleId, _ := c.GetInt("id", 0)
	strs := c.GetString("ids")

	o := orm.NewOrm()
	m := models.Role{Id: roleId}
	if err := o.Read(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Os dados são inválidos, atualize e tente novamente", "")
	}

	//删除已关联的历史数据
	if _, err := o.QueryTable(models.RoleResourceRelTBName()).Filter("role__id", m.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", "")
	}

	var relations []models.RoleResourceRel
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			r := models.Resource{Id: id}
			relation := models.RoleResourceRel{Role: &m, Resource: &r}
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

func (c *RoleController) UpdateSeq() {
	Id, _ := c.GetInt("pk", 0)
	oM, err := models.RoleOne(Id)
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

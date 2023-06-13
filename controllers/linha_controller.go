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
	"encoding/json"
	"fmt"
	"ccb_beego/enums"
	"ccb_beego/models"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type LinhaController struct {
	BaseController
}

func (c *LinhaController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid", "DataList", "SelectPicker")
}

func (c *LinhaController) Index() {
	c.Data["pageTitle"] = "Linha"
	c.Data["showMoreQuery"] = true

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "linha/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "linha/index_footerjs.html"

	//Controle de permissão de botão na página
	c.Data["canEdit"] = c.checkActionAuthor("LinhaController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("LinhaController", "Delete")
}

func (c *LinhaController) DataGrid() {
	var params models.LinhaQueryParam

	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//fmt.Println("Params:", params)
	//fmt.Println("Params:", params.ContatoNome)
	data, total := models.LinhaPageList(&params)

	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	c.Data["json"] = result
	c.ServeJSON()
}

// Lista de seleção suspensa
func (c *LinhaController) SelectPicker() {
	var params = models.LinhaQueryParam{}
	params.Estado = c.GetString("Estado")
	data := models.LinhaDataList(&params)
	c.jsonResult(enums.JRCodeSucc, "", data)
}

func (c *LinhaController) DataList() {
	var params = models.LinhaQueryParam{}
	fmt.Println("Params:", params)
	data := models.LinhaDataList(&params)
	c.jsonResult(enums.JRCodeSucc, "", data)
}

func (c *LinhaController) Edit() {
	//fmt.Println("Method:", c.Ctx.Request.Method)
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}

	Id, _ := c.GetInt(":id", 0)

	m := models.Linha{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		//fmt.Println("Sexo: ", m.Sexo)
		if err != nil {
			c.pageError("Os dados são inválidos, atualize e tente novamente")
		}
	} else {
		m.Estado = enums.Enabled
	}

	c.Data["m"] = m
	c.setTpl("linha/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "linha/edit_footerjs.html"
}

// add | update
func (c *LinhaController) Save() {
	var err error

	m := models.Linha{}

	//fmt.Println("Estado:", m.Estado)

	//Há um controle bootstapswitch, que precisa ser pré-processado
	//c.preform()

	//Obter o valor no formulário
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter dados", m.Id)
	}

	m.CorHex = m.CorHex[1:] //remove &

	o := orm.NewOrm()

	//fmt.Println("SAVE:", m.Estado)

	if m.Id == 0 {
		to, err := o.Begin()
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", m.Id)
			return
		}

		//m.CriadoEm = time.Now()

		if _, err = o.Insert(&m); err == nil {
			if err = to.Commit(); err != nil {
				c.jsonResult(enums.JRCodeFailed, "Falha na Alteração", m.Id)
				to.Rollback()
			} else {
				c.jsonResult(enums.JRCodeSucc, "Gravação com sucesso", m.Id)
			}
		} else {
			if err = to.Rollback(); err != nil {
				c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", m.Id)
			} else {
				c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", m.Id)
			}
		}
	} else {

		if _, err = o.Update(&m,
			"Nome",
			"Codigo",
			"MaterialNome",
			"MaterialFabricante",
			"MaterialTipo",
			"CorHex",
			"Estoque1",
			"Estoque2",
			"Minimo",
			"Pedido",
			"Estado"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "Atualizado", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "Falha ao modificar!", m.Id)
		}
	}
}

func (c *LinhaController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.LinhaBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("Excluído com êxito %d item", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha da rxclusão", 0)
	}
}

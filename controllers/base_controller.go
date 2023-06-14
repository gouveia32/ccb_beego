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
	"strings"

	"github.com/beego/beego/v2/core/logs"

	"ccb_beego/enums"
	"ccb_beego/models"

	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
	controllerName		string             //nome do controle
	actionName     		string             //nome da ação
	curUser        		models.BackendUser //Usuário
}

func (c *BaseController) Prepare() {
	//valor adicionado
	c.controllerName, c.actionName = c.GetControllerAndAction()

	c.Data["siteApp"], _ = beego.AppConfig.String("site.app")
	c.Data["siteName"], _ = beego.AppConfig.String("site.name")
	c.Data["siteVersion"], _ = beego.AppConfig.String("site.version")

	//Obtenha dados da sessão e defina as informações do usuário
	c.adapterUserInfo()
}

// Obter informações do usuário da sessão
func (c *BaseController) adapterUserInfo() {
	a := c.GetSession("backenduser")
	if a != nil {
		c.curUser = a.(models.BackendUser)
		c.Data["backenduser"] = a
	}
}

// checkLogin判断用户是否登录，未登录则跳转至登录页面
// 一定要在BaseController.Prepare()后执行
func (c *BaseController) checkLogin() {
	if c.curUser.Id == 0 {
		//登录页面地址
		urlstr := c.URLFor("HomeController.Login") + "?url="

		//登录成功后返回的址为当前
		returnURL := c.Ctx.Request.URL.Path

		//如果ajax请求则返回相应的错码和跳转的地址
		if c.Ctx.Input.IsAjax() {
			//由于是ajax请求，因此地址是header里的Referer
			returnURL := c.Ctx.Input.Refer()
			c.jsonResult(enums.JRCode302, "Por favor, faça login", urlstr+returnURL)
		}
		c.Redirect(urlstr+returnURL, 302)
		c.StopRun()
	}
}

// 判断某 Controller.Action 当前用户是否有权访问
func (c *BaseController) checkActionAuthor(ctrlName, ActName string) bool {
	if c.curUser.Id == 0 {
		return false
	}

	//从session获取用户信息
	user := c.GetSession("backenduser")

	//类型断言
	v, ok := user.(models.BackendUser)
	if ok {
		//如果是超级管理员，则直接通过
		if v.IsSuper == true {
			return true
		}

		//遍历用户所负责的资源列表
		for i, _ := range v.ResourceUrlForList {
			urlfor := strings.TrimSpace(v.ResourceUrlForList[i])
			if len(urlfor) == 0 {
				continue
			}
			// TestController.Get,:last,xie,:first,asta
			strs := strings.Split(urlfor, ",")
			if len(strs) > 0 && strs[0] == (ctrlName+"."+ActName) {
				return true
			}
		}
	}
	return false
}

// checkLogin julga se o usuário tem direito de acessar determinado endereço, caso contrário, ele irá pular para a página de erro
// 一定要在BaseController.Prepare()后执行
// 会调用checkLogin
// 传入的参数为忽略权限控制的Action
func (c *BaseController) checkAuthor(ignores ...string) {
	//Determinar se o login deve ser feito primeiro
	c.checkLogin()

	//如果Action在忽略列表里，则直接通用
	for _, ignore := range ignores {
		if ignore == c.actionName {
			return
		}
	}

	hasAuthor := c.checkActionAuthor(c.controllerName, c.actionName)
	if !hasAuthor {
		logs.Debug(fmt.Sprintf("author control: path=%s.%s userid=%v  Sem acesso", c.controllerName, c.actionName, c.curUser.Id))

		//Se não houver permissões
		if !hasAuthor {
			if c.Ctx.Input.IsAjax() {
				c.jsonResult(enums.JRCode401, "Sem acesso", "")
			} else {
				c.pageError("Sem acesso")
			}
		}
	}
}

// SetBackendUser2Session 获取用户信息（包括资源UrlFor）保存至Session
// 被 HomeController.DoLogin 调用
func (c *BaseController) setBackendUser2Session(userId int) error {
	m, err := models.BackendUserOne(userId)
	if err != nil {
		return err
	}

	//Obtenha uma lista de todos os recursos disponíveis para este usuário
	resourceList := models.ResourceTreeGridByUserId(userId, 1000)
	for _, item := range resourceList {
		m.ResourceUrlForList = append(m.ResourceUrlForList, strings.TrimSpace(item.UrlFor))
	}
	c.SetSession("backenduser", *m)
	return nil
}

// definir modelo
// O primeiro parâmetro é o modelo, o segundo parâmetro é o layout
func (c *BaseController) setTpl(template ...string) {
	layout := "shared/layout_page.html"

	var tplName string
	switch {
	case len(template) == 1:
		tplName = template[0]
	case len(template) == 2:
		tplName = template[0]
		layout = template[1]
	default:
		//Não use as 10 letras de Controlador
		ctrlName := strings.ToLower(c.controllerName[0 : len(c.controllerName)-10])
		actionName := strings.ToLower(c.actionName)
		tplName = ctrlName + "/" + actionName + ".html"
	}

	c.Layout = layout
	c.TplName = tplName
}

func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	res := &models.JsonResult{Code: code, Msg: msg, Obj: obj}
	c.Data["json"] = res
	c.ServeJSON()
	c.StopRun()
}

// redirecionar
func (c *BaseController) redirect(url string) {
	c.Redirect(url, 302)
	c.StopRun()
}

// Redirecionar para a página de erro
func (c *BaseController) pageError(msg string) {
	errorurl := c.URLFor("HomeController.Error") + "/" + msg
	c.Redirect(errorurl, 302)
	c.StopRun()
}

// Redirecionar para a página de login
func (c *BaseController) pageLogin() {
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
	c.StopRun()
}

// func (this *BaseController) GetControllerAndAction() (string, string) {

// }

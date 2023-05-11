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
	"strings"

	"github.com/beego/beego/v2/core/logs"

	"fmt"
	"ccb_beego/enums"
	"ccb_beego/models"
	"ccb_beego/utils"
	"io/ioutil"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Data["pageTitle"] = "Home"

	//判断是否登录
	c.checkLogin()

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("home/index.html")
	/* 	c.LayoutSections = make(map[string]string)
	   	c.LayoutSections["headcssjs"] = "home/index_headcssjs.html"
	   	c.LayoutSections["footerjs"] = "home/index_footerjs.html" */
}

func (c *HomeController) Index2() {
	c.Data["pageTitle"] = "servidor de nuvem"

	c.checkLogin()

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/index_headcssjs2.html"
	c.LayoutSections["footerjs"] = "home/index_footerjs2.html"
}

func (c *HomeController) Page404() {
	c.setTpl("home/page404.html", "shared/layout_base.html")
}

func (c *HomeController) Error() {
	c.Data["error"] = c.GetString(":error")
	c.setTpl("home/error.html", "shared/layout_pullbox.html")
}

func (c *HomeController) Login() {
	siteApp, _ := beego.AppConfig.String("site.app")
	sitename, _ := beego.AppConfig.String("site.name")
	c.Data["pageTitle"] = siteApp + sitename + " - Entrar"
	c.Data["siteVersion"], _ = beego.AppConfig.String("site.version")

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
	c.setTpl("home/login.html", "shared/layout_base.html")
}

// 退出
func (c *HomeController) Logout() {
	user := models.BackendUser{}
	c.SetSession("backenduser", user)
	c.pageLogin()
}

// 登陆
func (c *HomeController) DoLogin() {
	remoteAddr := c.Ctx.Request.RemoteAddr
	addrs := strings.Split(remoteAddr, "::1")
	if len(addrs) > 1 {
		remoteAddr = "localhost"
	}

	username := strings.TrimSpace(c.GetString("UserName"))
	userpwd := strings.TrimSpace(c.GetString("UserPwd"))

	if err := models.LoginTraceAdd(username, remoteAddr, time.Now()); err != nil {
		logs.Error("LoginTraceAdd error.")
	}
	logs.Info(fmt.Sprintf("login: %s IP: %s", username, remoteAddr))

	if len(username) == 0 || len(userpwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "Nome de usuário e senha incorretos", "")
	}

	userpwd = utils.String2md5(userpwd)
	user, err := models.BackendUserOneByUserName(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFailed, "O usuário está desativado, entre em contato com o administrador", "")
		}
		//保存用户信息到session
		c.setBackendUser2Session(user.Id)

		//获取用户信息
		c.jsonResult(enums.JRCodeSucc, "login bem-sucedido", "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "Nome de usuário ou senha está incorreta", "")
	}
}

//************************************* A P I ******************************************************

// 获取配置文件信息
func (c *HomeController) GetConfigValue() {
	key := c.GetString("key")
	result := ""
	err := true
	if key == "siteApp" {
		result, _ = beego.AppConfig.String("site.app")
		err = false
	} else if key == "siteName" {
		result, _ = beego.AppConfig.String("site.name")
		err = false
	} else if key == "siteVersion" {
		result, _ = beego.AppConfig.String("site.version")
		err = false
	}

	if err {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter parâmetros", key)
	} else {
		c.jsonResult(enums.JRCodeSucc, "", result)
	}
}

func (c *HomeController) GetWeather() {
	url := "http://api.openweathermap.org/data/2.5/weather?q=Guangzhou&appid=dafef1a9be486b5015eb92330a0add7d"
	ch := make(chan string)

	go func(url string, ch chan string) {
		resp, err := http.Get(url)
		if err != nil {
			logs.Error("erro de sincronização do tempo." + err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			logs.Error("Erro ao sincronizar o clima: " + err.Error())
			return
		}
		ch <- string(body)
	}(url, ch)

	reuslt := <-ch
	gjsonData := gjson.Parse(reuslt)
	temp := gjsonData.Get("main.temp").Float() - 274.15
	c.jsonResult(enums.JRCodeSucc, "", temp)
}

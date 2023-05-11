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

	"ccb_beego/enums"
	"ccb_beego/models"
	"ccb_beego/utils"

	"github.com/beego/beego/v2/client/orm"
)

type UserCenterController struct {
	BaseController
}

func (c *UserCenterController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	c.checkLogin()
}

func (c *UserCenterController) Profile() {
	Id := c.curUser.Id
	m, err := models.BackendUserOne(Id)
	if m == nil || err != nil {
		c.pageError("Os dados são inválidos, atualize e tente novamente")
	}
	c.Data["hasAvatar"] = len(m.Avatar) > 0
	logs.Debug(m.Avatar)

	c.Data["m"] = m
	c.setTpl()

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "usercenter/profile_headcssjs.html"
	c.LayoutSections["footerjs"] = "usercenter/profile_footerjs.html"
}

func (c *UserCenterController) BasicInfoSave() {
	Id := c.curUser.Id
	oM, err := models.BackendUserOne(Id)
	if oM == nil || err != nil {
		c.jsonResult(enums.JRCodeFailed, "Os dados são inválidos, atualize e tente novamente", "")
	}
	m := models.BackendUser{}

	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter dados", m.Id)
	}

	oM.RealName = m.RealName
	oM.Mobile = m.Mobile
	oM.Email = m.Email
	oM.Avatar = c.GetString("ImageUrl")
	if len(oM.Avatar) == 0 {
		oM.Avatar = "/static/upload/tigger.png"
	}

	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao modificar", m.Id)
	} else {
		c.setBackendUser2Session(Id)
		c.jsonResult(enums.JRCodeSucc, "保存成功", m.Id)
	}
}

func (c *UserCenterController) PasswordSave() {
	Id := c.curUser.Id
	oM, err := models.BackendUserOne(Id)
	if oM == nil || err != nil {
		c.pageError("Os dados são inválidos, atualize e tente novamente")
	}
	oldPwd := strings.TrimSpace(c.GetString("UserPwd", ""))
	newPwd := strings.TrimSpace(c.GetString("NewUserPwd", ""))
	confirmPwd := strings.TrimSpace(c.GetString("ConfirmPwd", ""))
	md5str := utils.String2md5(oldPwd)

	if oM.UserPwd != md5str {
		c.jsonResult(enums.JRCodeFailed, "A senha original está incorreta", "")
	}

	if len(newPwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "请输入新密码", "")
	}

	if newPwd != confirmPwd {
		c.jsonResult(enums.JRCodeFailed, "两次输入的新密码不一致", "")
	}

	oM.UserPwd = utils.String2md5(newPwd)
	o := orm.NewOrm()
	if _, err := o.Update(oM); err != nil {
		c.jsonResult(enums.JRCodeFailed, "保存Falha", oM.Id)
	} else {
		c.setBackendUser2Session(Id)
		c.jsonResult(enums.JRCodeSucc, "保存成功", oM.Id)
	}
}

func (c *UserCenterController) UploadImage() {
	//这里type没有用，只是为了演示传值
	stype, _ := c.GetInt32("type", 0)
	if stype > 0 {
		f, h, err := c.GetFile("fileImageUrl")
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "上传Falha", "")
		}
		defer f.Close()

		filePath := "static/upload/" + h.Filename
		// 保存位置在 static/upload, 没有文件夹要先创建
		c.SaveToFile("fileImageUrl", filePath)
		c.jsonResult(enums.JRCodeSucc, "上传成功", "/"+filePath)
	} else {
		c.jsonResult(enums.JRCodeFailed, "上传Falha", "")
	}
}

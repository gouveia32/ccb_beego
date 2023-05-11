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

package sysinit

import (
	_ "gardens/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// Inicializar conexão de dados
func InitDatabase() {
	//Leia o arquivo de configuração e defina os parâmetros do banco de dados
	//classe de banco de dados
	dbType, err := beego.AppConfig.String("db_type")
	if err != nil {
		logs.Error("db_type is null.", err)
		return
	}

	//nome da conexão
	dbAlias, err := beego.AppConfig.String(dbType + "::db_alias")
	if err != nil {
		logs.Error("db_alias is null.", err)
		return
	}

	//Banco de dados de nomes
	dbName, err := beego.AppConfig.String(dbType + "::db_name")
	if err != nil {
		logs.Error("db_name is null.", err)
		return
	}

	//nome de usuário de conexão de banco de dados
	dbUser, err := beego.AppConfig.String(dbType + "::db_user")
	if err != nil {
		logs.Error("db_user is null.", err)
		return
	}

	//nome de usuário de conexão de banco de dados
	dbPwd, err := beego.AppConfig.String(dbType + "::db_pwd")
	if err != nil {
		logs.Error("db_pwd is null.", err)
		return
	}

	//IP do banco de dados (nome de domínio)
	dbHost, err := beego.AppConfig.String(dbType + "::db_host")
	if err != nil {
		logs.Error("db_host is null.", err)
		return
	}

	//porta do banco de dados
	dbPort, err := beego.AppConfig.String(dbType + "::db_port")
	if err != nil {
		logs.Error("db_port is null.", err)
		return
	}

	//dbName2 := "kxtimingdata"

	switch dbType {
	/* 	case "sqlite3":
	orm.RegisterDataBase(dbAlias, dbType, dbName) */
	case "mysql":
		dbCharset, err := beego.AppConfig.String(dbType + "::db_charset")
		if err != nil {
			logs.Error("db_charset is null.", err)
			return
		}

		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset+"&loc=Asia%2FShanghai")

		//orm.RegisterDataBase(dbName2, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName2+"?charset="+dbCharset+"&loc=Asia%2FShanghai")
	}

	//Se for o modo de desenvolvimento, exibir informações de comando
	runMode, _ := beego.AppConfig.String("RunMode")
	isDev := runMode == "dev"

	// true: drop table 后再建表
	force := false

	//Criação automática de tabelas
	orm.RunSyncdb("default", force, isDev)
	if isDev {
		orm.Debug = isDev
	}
}

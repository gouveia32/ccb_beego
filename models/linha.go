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

package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Linha struct {
	Id                 int    `orm:"column(id)" form:"Id"`
	Codigo             string `orm:"column(codigo)" form:"Codigo"`
	Nome               string `orm:"column(nome)" form:"Nome"`
	MaterialNome       string `orm:"column(material_nome)" form:"MaterialNome"`
	MaterialFabricante string `orm:"column(material_fabricante)" form:"MaterialFabricante"`
	MaterialTipo       string `orm:"column(material_ipo)" form:"MaterialTipo"`
	CorHex             string `orm:"column(cor_hex)" form:"CorHex"`
	Estoque1           int    `orm:"column(estoque_1)" form:"Estoque1"`
	Estoque2           int    `orm:"column(estoque_2)" form:"Estoque2"`
	Minimo             int    `orm:"column(minimo)" form:"Minimo"`
	Pedido             int    `orm:"column(pedido)" form:"Pedido"`
	Estado             int8   `orm:"column(estado)" form:"Estado"`
}

type LinhaQueryParam struct {
	BaseQueryParam
	Nome   string
	Codigo string
	Estado string
	Alerta string
}

func init() {
	orm.RegisterModel(new(Linha))
}

func LinhaTBName() string {
	return "linha"
}

func LinhaPageList(params *LinhaQueryParam) ([]*Linha, int64) {
	query := orm.NewOrm().QueryTable(LinhaTBName())
	data := make([]*Linha, 0)

	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Nome":
		sortorder = "Nome"
	case "Codigo":
		sortorder = "Codigo"
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	//fmt.Println("nome:", params.Nome)

	//fmt.Println("contato:", params.Contato)

	query = query.Filter("Nome__icontains", params.Nome)
	query = query.Filter("Codigo__icontains", params.Codigo)
	query = query.Filter("estado__istartswith", params.Estado)
	if (params.Alerta > "0") {
		query = query.Filter("estoque_1__lt","Minimo")
	}
	

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)

	return data, total
}

func LinhaDataList(params *LinhaQueryParam) []*Linha {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := LinhaPageList(params)
	return data
}

func LinhaBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(LinhaTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func LinhaOne(id int) (*Linha, error) {
	o := orm.NewOrm()
	m := Linha{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (this *Linha) TableName() string {
	return LinhaTBName()
}

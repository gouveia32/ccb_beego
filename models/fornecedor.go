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
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Fornecedor struct {
	Id            int       `orm:"column(id)" form:"Id"`
	Nome          string    `orm:"column(nome)" form:"Nome"`
	ContatoFuncao string    `orm:"column(contato_funcao)" form:"ContatoFuncao"`
	ContatoNome   string    `orm:"column(contato_nome)" form:"ContatoNome"`
	CgcCpf        string    `orm:"column(cgc_cpf)" form:"CgcCpf"`
	RazaoSocial   string    `orm:"column(razao_social)" form:"RazaoSocial"`
	InscrEstadual string    `orm:"columniInscr_estadual)" form:"InscrEstadual"`
	Endereco      string    `orm:"column(endereco)" form:"Endereco"`
	Cidade        string    `orm:"column(cidade)" form:"Cidade"`
	Uf            string    `orm:"column(uf)" form:"Uf"`
	Cep           string    `orm:"column(cep)" form:"Cep"`
	Telefone1     string    `orm:"column(telefone_1)" form:"Telefone1"`
	Telefone2     string    `orm:"column(telefone_2)" form:"Telefone2"`
	Telefone3     string    `orm:"column(telefone_3)" form:"Telefone3"`
	Email         string    `orm:"column(email)" form:"Email"`
	Obs           string    `orm:"column(obs)" form:"Obs"`
	CriadoEm      time.Time `orm:"column(criado_em)" form:"CriadoEm"`
	AlteradoEm    time.Time `orm:"auto_now;type(datetime);column(alterado_em)" form:"AlteradoEm"`
	Estado        int       `orm:"column(estado)" form:"Estado"`
}

type FornecedorQueryParam struct {
	BaseQueryParam
	Nome    string
	Contato string
	Estado  string
}

func init() {
	orm.RegisterModel(new(Fornecedor))
}

func FornecedorTBName() string {
	return "fornecedor"
}

func FornecedorPageList(params *FornecedorQueryParam) ([]*Fornecedor, int64) {
	query := orm.NewOrm().QueryTable(FornecedorTBName())
	data := make([]*Fornecedor, 0)

	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Nome":
		sortorder = "Nome"
	case "ContatoNome":
		sortorder = "ContatoNome"
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	//fmt.Println("nome:", params.Nome)

	//fmt.Println("contato:", params.Contato)

	query = query.Filter("Nome__icontains", params.Nome)
	query = query.Filter("ContatoNome__icontains", params.Contato)
	query = query.Filter("estado__istartswith", params.Estado)

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)

	return data, total
}

func FornecedorDataList(params *FornecedorQueryParam) []*Fornecedor {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := FornecedorPageList(params)
	return data
}

func FornecedorBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(FornecedorTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func FornecedorOne(id int) (*Fornecedor, error) {
	o := orm.NewOrm()
	m := Fornecedor{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}


func (this *Fornecedor) TableName() string {
	return FornecedorTBName()
}

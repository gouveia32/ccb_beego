
package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type GrupoQueryParam struct {
	BaseQueryParam
	NomeLike string
}


//Grupo
type Grupo struct {
	Id				int    `form:"Id"`
	Nome			string `form:"Nome"`
}

func init() {
	orm.RegisterModel(new(Grupo))
}

//Obtenha o nome da tabela correspondente ao Grupo
func GrupoTBName() string {
	return "grupo"
}

func (c *Grupo) TableName() string {
	return GrupoTBName()
}

//Obter dados paginados
func GrupoPageList(params *GrupoQueryParam) ([]*Grupo, int64) {
	query := orm.NewOrm().QueryTable(GrupoTBName())
	data := make([]*Grupo, 0)

	//classificação padrão
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("nome__istartswith", params.NomeLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

//Obter lista de grupos
func GrupoDataList(params *GrupoQueryParam) []*Grupo {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := GrupoPageList(params)
	return data
}

//exclusão em lote
func GrupoBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(GrupoTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func GrupoOne(id int) (*Grupo, error) {
	o := orm.NewOrm()
	m := Grupo{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

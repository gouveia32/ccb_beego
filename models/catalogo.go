
package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type CatalogoQueryParam struct {
	BaseQueryParam
	NomeLike string
}


//Catalogo
type Catalogo struct {
	Id				int    `form:"Id"`
	Nome			string `form:"Nome"`
	Seq             int
	CatalogoBordadoRel []*CatalogoBordadoRel `orm:"reverse(many)" json:"-"` // Configurar uma relação inversa de um para muitos
}

func init() {
	orm.RegisterModel(new(Catalogo))
}

//Obtenha o nome da tabela correspondente ao Catalogo
func CatalogoTBName() string {
	return "catalogo"
}

func (c *Catalogo) TableName() string {
	return CatalogoTBName()
}

//Obter dados paginados
func CatalogoPageList(params *CatalogoQueryParam) ([]*Catalogo, int64) {
	query := orm.NewOrm().QueryTable(CatalogoTBName())
	data := make([]*Catalogo, 0)

	//classificação padrão
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Seq":
		sortorder = "Seq"
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("nome__istartswith", params.NomeLike)
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

//Obter lista de catalogos
func CatalogoDataList(params *CatalogoQueryParam) []*Catalogo {
	params.Limit = -1
	params.Sort = "Seq"
	params.Order = "asc"
	data, _ := CatalogoPageList(params)
	return data
}

//exclusão em lote
func CatalogoBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(CatalogoTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func CatalogoOne(id int) (*Catalogo, error) {
	o := orm.NewOrm()
	m := Catalogo{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

package models

import (
	"fmt"
	"image/color"

	"github.com/beego/beego/v2/client/orm"
)

type Linha struct {
	//Id                 int    `orm:"column(id)" form:"Id"`

	Codigo             	string `orm:"pk; column(codigo)" form:"Codigo"`
	Nome               	string `orm:"column(nome)" form:"Nome"`
	MaterialNome       	string `orm:"column(material_nome)" form:"MaterialNome"`
	MaterialFabricante 	string `orm:"column(material_fabricante)" form:"MaterialFabricante"`
	MaterialTipo       	string `orm:"column(material_tipo)" form:"MaterialTipo"`
	CorHex             	string `orm:"column(cor_hex)" form:"CorHex"`
	CorRGB				color.Color	`orm:"-"`
	Estoque1           	int    `orm:"column(estoque_1)" form:"Estoque1"`
	Estoque2           	int    `orm:"column(estoque_2)" form:"Estoque2"`
	Minimo             	int    `orm:"column(minimo)" form:"Minimo"`
	Pedido             	int    `orm:"column(pedido)" form:"Pedido"`
	Estado             	int8   `orm:"column(estado)" form:"Estado"`
	New                	bool
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

	sortorder := "Codigo"
	switch params.Sort {
	/* 	case "Id":
	sortorder = "Id"
	*/case "Nome":
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
	fmt.Println("Alerta", params.Alerta)
	if params.Alerta == "1" {
		fmt.Println("Alerta: Aqui", params.Alerta)
		query = query.FilterRaw("estoque_1", "<= minimo")
	}

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)

	return data, total
}

func LinhaDataList(params *LinhaQueryParam) []*Linha {
	params.Limit = -1
	params.Sort = "Codigo"
	params.Order = "asc"
	data, _ := LinhaPageList(params)
	return data
}

func LinhaBatchDelete(cods []string) (int64, error) {
	query := orm.NewOrm().QueryTable(LinhaTBName())
	num, err := query.Filter("codigo__in", cods).Delete()
	return num, err
}

func LinhaOne(codigo string) (*Linha, error) {
	o := orm.NewOrm()
	m := Linha{Codigo: codigo}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (this *Linha) TableName() string {
	return LinhaTBName()
}


func ParseHexColor(s string) (c color.RGBA, err error) {
    c.A = 0xff
    switch len(s) {
    case 7:
        _, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
    case 4:
        _, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
        // Double the hex digits:
        c.R *= 17
        c.G *= 17
        c.B *= 17
    default:
        err = fmt.Errorf("invalid length, must be 7 or 4")

    }
    return
}
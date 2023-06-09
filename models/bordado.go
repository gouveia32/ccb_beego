package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Bordado struct {
	Id        int       `orm:"column(id)" form:"Id"`
	Arquivo   string  	 `orm:"size(40)"`
	Descricao string  	 `orm:"size(50)"`
	Caminho   string  	 `orm:"size(225)"`
	Disquete  string  	 `orm:"size(12)"`
	Bastidor  string   `orm:"size(12)"`
	GrupoId   int   	`orm:"column(grupo_id);default(1);" form:"GrupoId"`
	Preco     float64 
	Pontos    int64   
	Cores     int16   
	Largura   int64   
	Altura    int64   
	Metragem  int64   
	Aprovado  bool
	Alerta    bool
	Imagem    string  `form:imagem orm:"column(imagem);type(blob)`
	//Imagem      string  `orm:"column(imagem);type(blob)"                      description:"imagem derada a partir do dst" form:"imagem"`
	CorFundo    string `orm:"column(cor_fundo);nil" 					     description:"cor de fundo" form:"CorFundo"`
	ObsPublica 	string `orm:"size(1024)"`
	ObsRestrita string	`orm:"size(1024)"`
	CriadoEm      time.Time `orm:"column(criadoEm)" form:"CriadoEm"`
	AlteradoEm    time.Time `orm:"column(alteradoEm)" form:"AlteradoEm"`
	Estado      int
	CatalogoIds		[]int	`orm:"-" form:"CatalogoIds"`
	CatalogoBordadoRel	[]*CatalogoBordadoRel `orm:"reverse(many)"` // Configurar uma relação inversa de um para muitos
}

type BordadoQueryParam struct {
	BaseQueryParam
	ArquivoLike		string
	DescricaoLike	string
	Estado  		string
}

func init() {
	orm.RegisterModel(new(Bordado))
}

func BordadoTBName() string {
	return "bordado"
}

func BordadoPageList(params *BordadoQueryParam) ([]*Bordado, int64) {
	query := orm.NewOrm().QueryTable(BordadoTBName())
	data := make([]*Bordado, 0)

	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	case "Arquivo":
		sortorder = "Arquivo"
	case "Descricao":
		sortorder = "Descricao"
	}

	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}

	//fmt.Println("nome:", params.Nome)

	//fmt.Println("contato:", params.Contato)

	query = query.Filter("Arquivo__icontains", params.ArquivoLike)
	query = query.Filter("Descricao__icontains", params.DescricaoLike)
	query = query.Filter("estado__istartswith", params.Estado)

	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)

	return data, total
}

// Obter um único item por id
func BordadoOne(id int) (*Bordado, error) {
	o := orm.NewOrm()
	m := Bordado{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func BordadoDataList(params *BordadoQueryParam) []*Bordado {
	params.Limit = -1
	params.Sort = "Id"
	params.Order = "asc"
	data, _ := BordadoPageList(params)
	return data
}

func BordadoBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(BordadoTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

func (this *Bordado) TableName() string {
	return BordadoTBName()
}

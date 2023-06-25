package models

import (
	"log"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Linha e Relacionamento com o Bordado
type LinhaBordadoRel struct {
	Id       int
	Bordado  *Bordado `orm:"rel(fk);" ` //chave estrangeira
	Seq      int
	Linha    *Linha    `orm:"rel(fk);"` //chave estrangeira
	CriadoEm time.Time `orm:"auto_now_add;type(datetime);null"`
}

type LinhaBordadoQueryParam struct {
	BaseQueryParam
	BordadoId int
	LinhaId   string
}

// multiple fields index
func (l *LinhaBordadoRel) TableIndex() [][]string {
	return [][]string{
		[]string{"Bordado", "Seq"},
	}
}

func init() {
	orm.RegisterModel(new(LinhaBordadoRel))
}

// Tabela de relacionamento muitos-para-muitos de função e usuário
func LinhaBordadoRelTBName() string {
	return "linha_bordado_rel"
}

func LinhaBordadoPageList(bordado_id int) []*LinhaBordadoRel {
	o := orm.NewOrm()

	sql := "SELECT * FROM linha_bordado_rel WHERE  bordado_id = ?;"

	var members []*LinhaBordadoRel
	_, err := o.Raw(sql, bordado_id).QueryRows(&members)
	//fmt.Println("members:", members)
	if err != nil {
		log.Fatal(err)
	}
	return members
}

func (l *LinhaBordadoRel) TableName() string {
	return LinhaBordadoRelTBName()
}

func (l *LinhaBordadoRel) getLinhasByBordado(bordado_id int) ([]*Linha, int64) {
	query := orm.NewOrm().QueryTable(LinhaBordadoRelTBName())
	data := make([]*Linha, 0)

	query = query.Filter("Bordado", bordado_id)
	total, _ := query.Count()
	query.All(&data)

	return data, total
}

package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

//Catalogo e Relacionamento com o Bordado
type CatalogoBordadoRel struct {
	Id          int
	Catalogo		*Catalogo		`orm:"rel(fk)"`  //chave estrangeira
	Bordado		*Bordado	`orm:"rel(fk)" ` //chave estrangeira
	CriadoEm    time.Time	`orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(CatalogoBordadoRel))
}

//Tabela de relacionamento muitos-para-muitos de função e usuário
func CatalogoBordadoRelTBName() string {
	return "catalogo_bordado_rel"
}

func (c *CatalogoBordadoRel) TableName() string {
	return CatalogoBordadoRelTBName()
}

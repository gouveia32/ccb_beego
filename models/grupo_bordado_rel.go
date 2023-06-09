package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

//Grupo e Relacionamento com o Bordado
type GrupoBordadoRel struct {
	Id          int
	Grupo		*Grupo		`orm:"rel(fk)"`  //chave estrangeira
	Bordado		*Bordado	`orm:"rel(fk)" ` //chave estrangeira
	CriadoEm    time.Time	`orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(GrupoBordadoRel))
}

//Tabela de relacionamento muitos-para-muitos de função e usuário
func GrupoBordadoRelTBName() string {
	return "grupo_bordado_rel"
}

func (c *GrupoBordadoRel) TableName() string {
	return GrupoBordadoRelTBName()
}

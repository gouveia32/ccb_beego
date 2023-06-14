package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

//Linha e Relacionamento com o Bordado
type LinhaBordadoRel struct {
	Id          int
	Bordado		*Bordado	`orm:"rel(fk);" ` //chave estrangeira
	Seq			int	
	Linha		*Linha		`orm:"rel(fk);"`  //chave estrangeira
	CriadoEm    time.Time	`orm:"auto_now_add;type(datetime);null"`
}

// multiple fields index
func (l *LinhaBordadoRel) TableIndex() [][]string {
	return [][]string{
		[]string{"Bordado","Seq"},
	}
}


func init() {
	orm.RegisterModel(new(LinhaBordadoRel))
}

//Tabela de relacionamento muitos-para-muitos de função e usuário
func LinhaBordadoRelTBName() string {
	return "linha_bordado_rel"
}

func (c *LinhaBordadoRel) TableName() string {
	return LinhaBordadoRelTBName()
}

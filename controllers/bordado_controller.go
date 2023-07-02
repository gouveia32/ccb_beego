package controllers

import (
	"bytes"
	"ccb_beego/enums"
	"ccb_beego/models"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"

	//"github.com/go-playground/colors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/StephaneBunel/bresenham"

	"github.com/beego/beego/v2/client/orm"
)

type BordadoController struct {
	BaseController
}

func (c *BordadoController) Prepare() {
	c.BaseController.Prepare()
	c.checkAuthor("DataGrid", "DataList", "SelectPicker")
}

func (c *BordadoController) Index() {
	c.Data["pageTitle"] = "Bordado"
	c.Data["showMoreQuery"] = true

	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "bordado/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "bordado/index_footerjs.html"

	//Controle de permissão de botão na página
	c.Data["canEdit"] = c.checkActionAuthor("BordadoController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("BordadoController", "Delete")
}

func (c *BordadoController) DataGrid() {
	var params models.BordadoQueryParam

	json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	//fmt.Println("Params:", params)
	//fmt.Println("Params:", params.ContatoNome)
	data, total := models.BordadoPageList(&params)

	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data

	c.Data["json"] = result
	c.ServeJSON()
}

// Lista de seleção suspensa
func (c *BordadoController) SelectPicker() {
	var params = models.BordadoQueryParam{}
	params.Estado = c.GetString("Estado")
	data := models.BordadoDataList(&params)
	c.jsonResult(enums.JRCodeSucc, "", data)
}

func (c *BordadoController) DataList() {
	var params = models.BordadoQueryParam{}
	fmt.Println("Params:", params)
	data := models.BordadoDataList(&params)
	c.jsonResult(enums.JRCodeSucc, "", data)
}

// *
// *
// ***************** Edit **************************
func (c *BordadoController) Edit() {
	//fmt.Println("Method:", c.Ctx.Request.Method)
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}

	Id, _ := c.GetInt(":id", 0)
	m := &models.Bordado{}
	var err error
	if Id > 0 {
		m, err = models.BordadoOne(Id)
		if err != nil {
			c.pageError("Os dados são inválidos, atualize e tente novamente")
		}

		ct := orm.NewOrm()
		ct.LoadRelated(m, "CatalogoBordadoRel")
		ct.LoadRelated(m, "LinhaBordadoRel")
		//Ordena por bordado_it + seq
		sort.SliceStable(m.LinhaBordadoRel, func(i, j int) bool {
			return m.LinhaBordadoRel[i].Seq < m.LinhaBordadoRel[j].Seq
		})

		//fmt.Println("m.CatalogoBordadoRel:",m.CatalogoBordadoRel)
	} else {
		//Ativado por padrão ao adicionar bordados
		m.Estado = enums.Enabled
	}

	c.Data["imagem"] = m.Imagem
	//c.Data["img"] = CarregaDst(color.White)

	//c.Data["img"] = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg=="

	ufs := models.GetUFs()
	c.Data["ufs"] = ufs

	var params = models.GrupoQueryParam{}
	grupos := models.GrupoDataList(&params)

	//fmt.Println("grupos:",grupos)
	c.Data["grupos"] = grupos

	c.Data["m"] = m

	//Obtenha a lista de catalogoId associada
	var catalogoIds []string
	for _, item := range m.CatalogoBordadoRel {
		catalogoIds = append(catalogoIds, strconv.Itoa(item.Catalogo.Id))
	}

	//Obtenha a lista de linhaId associada
	linhas := make([]*models.LinhaBordadoRel, 0)
	for _, item := range m.LinhaBordadoRel {
		linha, _ := models.LinhaOne(item.Linha.Codigo)

		item.Linha.Codigo = linha.Codigo
		item.Linha.Nome = linha.Nome
		item.Linha.CorHex = linha.CorHex

		linhas = append(linhas, item)
	}

	lps := models.GetLp()

	//fmt.Println("lps:", lps)

	if len(linhas) < int(m.Cores) {
		for i := len(m.LinhaBordadoRel); i < int(m.Cores); i++ {
			linha, err := models.LinhaOne(lps[i])
			if err != nil {
				linha, _ = models.LinhaOne("5311")
			}
			//fmt.Println("linha:", linha.Codigo, " ", linha.Nome)

			item := models.LinhaBordadoRel{Bordado: m, Linha: linha, Seq: i + 1}
			linhas = append(linhas, &item)
		}
	}

	//fmt.Println("linhaIds final:", linhas)

	c.Data["linhas"] = linhas
	c.Data["catalogos"] = strings.Join(catalogoIds, ",")
	c.setTpl("bordado/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "bordado/edit_footerjs.html"
}

// add | update
func (c *BordadoController) Save() {
	var err error

	b := models.Bordado{}
	o := orm.NewOrm()

	//fmt.Println("Estado:", m.Estado)

	//Há um controle bootstapswitch, que precisa ser pré-processado
	//c.preform()

	//Obter o valor no formulário
	if err = c.ParseForm(&b); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao obter dados", b.Id)
	}

	//fmt.Printf("\n\nc: %+v\n\n", b) // annot call non-function r.Form (type url.Values)

	//fmt.Println("Imagem:", b.Imagem)

	//Excluir catalogos históricos associados
	if _, err := o.QueryTable(models.CatalogoBordadoRelTBName()).Filter("bordado__id", b.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", "")
	}

	//fmt.Println("AQUI")
	//adicionar relacionamento catalogo
	var relscat []models.CatalogoBordadoRel
	for _, catalogoId := range b.CatalogoIds {
		ct := models.Catalogo{Id: catalogoId}
		rel := models.CatalogoBordadoRel{Bordado: &b, Catalogo: &ct}
		relscat = append(relscat, rel)
	}

	if len(relscat) > 0 {
		//adicionar lote
		if _, err := o.InsertMulti(len(relscat), relscat); err != nil {
			c.jsonResult(enums.JRCodeFailed, "Falha ao Salvar", b.Id)
		}
	}

	//Excluir linhas históricos associados
	if _, err := o.QueryTable(models.LinhaBordadoRelTBName()).Filter("bordado__id", b.Id).Delete(); err != nil {
		c.jsonResult(enums.JRCodeFailed, "Falha ao excluir", "")
	}

	//fmt.Println("AQUI :  linhas", b.linhaCod)
	//adicionar relacionamento linha
	var seq = 1
	var relslin []models.LinhaBordadoRel
	for _, linhaCod := range b.LinhaCods {
		ln := models.Linha{Codigo: linhaCod}

		rel := models.LinhaBordadoRel{Bordado: &b, Linha: &ln, Seq: seq}
		relslin = append(relslin, rel)
		seq = seq + 1
	}

	if len(relslin) > 0 {
		//adicionar lote
		if _, err := o.InsertMulti(len(relslin), relslin); err != nil {
			c.jsonResult(enums.JRCodeFailed, "Falha ao Salvar", b.Id)
		}
	}

	if b.Id == 0 {
		to, err := o.Begin()
		if err != nil {
			c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", b.Id)
			return
		}

		b.CriadoEm = time.Now()
		//b.Estado = enums.Enabled

		if _, err = o.Insert(&b); err == nil {
			if err = to.Commit(); err != nil {
				c.jsonResult(enums.JRCodeFailed, "Falha na Alteração", b.Id)
				to.Rollback()
			} else {
				c.jsonResult(enums.JRCodeSucc, "Gravação com sucesso", b.Id)
			}
		} else {
			if err = to.Rollback(); err != nil {
				c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", b.Id)
			} else {
				c.jsonResult(enums.JRCodeFailed, "Falha ao adicionar", b.Id)
			}
		}
	} else {
		b.AlteradoEm = time.Now()
		fmt.Println("Grupo_id:", b.GrupoId)
		if _, err = o.Update(&b,
			"Arquivo",
			"Descricao",
			"Caminho",
			"Disquete",
			"Bastidor",
			"GrupoId",
			"Preco",
			"Pontos",
			"Cores",
			"Largura",
			"Altura",
			"Metragem",
			"Aprovado",
			"Alerta",
			"Imagem",
			//Imagem,
			"CorFundo",
			"ObsPublica",
			"ObsRestrita",
			"CriadoEm",
			"AlteradoEm",
			"ObsPublica",
			"ObsRestrita",
			"Estado"); err == nil {
			c.jsonResult(enums.JRCodeSucc, "Atualizado com sucesso", b.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "Falha ao modificar", b.Id)
		}
	}
}

func (c *BordadoController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}

	if num, err := models.BordadoBatchDelete(ids); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("Excluído com êxito %d item", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "Falha da rxclusão", 0)
	}
}

// *
// *
// ***************** CarregaDst **************************
func CarregaDst(cor color.Color) (resp string) {
	// Cria uma imagem com fundo branco
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	//draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Desenha um retângulo preto na imagem
	rect := image.Rect(50, 50, 100, 100)
	draw.Draw(img, rect, &image.Uniform{cor}, image.ZP, draw.Src)

	// Codifica a imagem em base64
	var buf bytes.Buffer
	png.Encode(&buf, img)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return b64
}

// *
// *
// ***************** DrawLine **************************
func DrawLine(x1, y1, x2, y2 int, cor color.Color) (resp string) {
	var imgRect = image.Rect(0, 0, 300, 300)
	var img = image.NewRGBA(imgRect)

	// draw line
	bresenham.DrawLine(img, x1, y1, x2, y2, cor)

	// Codifica a imagem em base64
	var buf bytes.Buffer
	png.Encode(&buf, img)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return b64

}

// *
// *
// ***************** LerDst **************************
func (c *BordadoController) LerDst() {
	cod_linha := c.GetString("cor")
	id, _ := c.GetInt("id", 0) //id do bordado

	seq := c.GetString("seq")

	//fileName := c.GetString("arq")

	fmt.Println("arq: ", seq)

	var arq string
	if seq != "" {
		arq = "c:/bordados/Bota.dst"
	} else {
		arq = cod_linha
	}

	data, err := ioutil.ReadFile(arq)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	//fmt.Println("Tamanho: ", binary.Size(data))

	var reg [3]int
	//var p int = -1
	var corHex = ""

	X := 0
	Y := 0
	salto := false

	var imgRect = image.Rect(0, 0, 500, 500)
	var img = image.NewRGBA(imgRect)

	if corHex == "" {
		l, _ := models.LinhaOne(cod_linha)
		corHex = l.CorHex
		//fmt.Println("NOVA:", corHex)
	}

	if id > 0 {

		linhas := models.LinhaBordadoPageList(id)

		fmt.Println("Bord id:", id)
		fmt.Println("linhas: ", linhas)

		cor1 := color.Black
		fmt.Printf("data=%d ", binary.Size(data))

		//sData := string(data)
		//fmt.Println("sData", sData)

		Xmais, _ := strconv.Atoi(string(data[41:46]))
		Xmenos, _ := strconv.Atoi(string(data[50:55]))
		Ymais, _ := strconv.Atoi(string(data[59:64]))
		Ymenos, _ := strconv.Atoi(string(data[68:73]))

		Largura := Xmais + Xmenos
		Altura := Ymais + Ymenos
		NrPontos, _ := strconv.Atoi(string(data[23:30]))
		Cores, _ := strconv.Atoi(string(data[34:37]))
		Cores++
		X0 := Xmenos
		Y0 := Ymais

		zoom := 100
		if Largura > Altura && Largura != 0 {
			zoom = 500 * 79 / Largura
		} else if Altura != 0 {
			zoom = 500 * 79 / Altura
		}

		fmt.Printf("           Header: %d %d %d %d %d %d %d %d zoom:%d\n", Xmais, Xmenos, Ymais, Ymenos, Largura, Altura, NrPontos, Cores, zoom)

		X = X0
		Y = Y0

		for i, e := range data {
			if i >= 512 && i < binary.Size(data) {
				b := (i - 2) % 3
				reg[b] = int(e)
				//fmt.Printf("\nbyte :%d  %d (%d)", i, b, reg[b])
				if b == 2 { //terceiro byte

					if (reg[2] & 64) == 64 { //troca de cor

						//fmt.Printf("\nTroca de cor :%d    =   %d", i, reg[2]&0x40)
					}
					salto = false
					if (reg[2] & 128) == 128 {
						salto = true
					}
					if (reg[2] & 4) == 4 {
						X += 81
					}
					if (reg[2] & 8) == 8 {
						X -= 81
					}
					if (reg[2] & 16) == 16 {
						Y += 81
					}
					if (reg[2] & 32) == 32 {
						Y -= 81
					}

					if (reg[1] & 1) == 1 {
						X += 3
					}
					if (reg[1] & 2) == 2 {
						X -= 3
					}
					if (reg[1] & 4) == 4 {
						X += 27
					}
					if (reg[1] & 8) == 8 {
						X -= 27
					}
					if (reg[1] & 16) == 16 {
						Y += 27
					}
					if (reg[1] & 32) == 32 {
						Y -= 27
					}
					if (reg[1] & 64) == 64 {
						Y += 3
					}
					if (reg[1] & 128) == 128 {
						Y -= 3
					}

					if (reg[0] & 1) == 1 {
						X += 1
					}
					if (reg[0] & 2) == 2 {
						X -= 1
					}
					if (reg[0] & 4) == 4 {
						X += 9
					}
					if (reg[0] & 8) == 8 {
						X -= 9
					}
					if (reg[0] & 16) == 16 {
						Y += 9
					}
					if (reg[0] & 32) == 32 {
						Y -= 9
					}
					if (reg[1] & 64) == 64 {
						Y += 1
					}
					if (reg[1] & 128) == 128 {
						Y -= 1
					}

					XX0 := X0 * zoom / 100
					YY0 := Y0 * zoom / 100
					XX := X * zoom / 100
					YY := Y * zoom / 100

					if XX-XX0 > 64|YY-YY0 {
						salto = true
					}

					if !salto {
						bresenham.DrawLine(img, XX0, YY0, XX, YY, cor1)
					}
					//fmt.Printf("\nbytes  %d : %d %d %d %d", i%3, X0, Y0, X, Y)

					X0 = X
					Y0 = Y
				}
			}
		}

		// Codifica a imagem em base64
		var buf bytes.Buffer
		png.Encode(&buf, img)
		b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

		c.Data["json"] = corHex + b64

		c.ServeJSON()
	}

	fmt.Println("LerDst:", cod_linha)
}

func normalize(colorStr string) (string, error) {
	// left as an exercise for the reader
	return colorStr, nil
}

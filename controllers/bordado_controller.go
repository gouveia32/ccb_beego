package controllers

import (
	"bytes"
	"ccb_beego/enums"
	"ccb_beego/models"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	//"github.com/go-playground/colors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
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
		item.Linha.CorRGB, err = models.ParseHexColor (linha.CorHex)
		if err != nil {
			item.Linha.CorRGB = color.RGBA{255,1,1,255}
		}

		linhas = append(linhas, item)
	}

	lps := models.GetLp()

	//fmt.Println("lps:", lps)

	if len(linhas) != int(m.Cores) {
		for i := len(m.LinhaBordadoRel); i < int(m.Cores); i++ {
			linha, err := models.LinhaOne(lps[i])
			if err != nil {
				linha, _ = models.LinhaOne("5311")
			}
			fmt.Println("linha:", linha.Codigo, " ", linha.Nome)

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
// ***************** Corres **************************
func CarregaCoresO(id int, nc int) []color.Color {
	cores_padrao := []color.Color{}
	data := []color.Color{}
	linhas := models.LinhaBordadoPageList(id)
	if len(linhas) > 0 {
		for _, linha := range linhas {
			l, err := models.LinhaOne(linha.Linha.Codigo)
			if err != nil {
				fmt.Println("linha: ", l.Nome)
			}
			//fmt.Println("linha: ", l.Nome)
			colorStr, err := normalize(l.CorHex)
			if err != nil {
				log.Fatal(err)
			}
			b, err1 := hex.DecodeString(colorStr)
			if err1 != nil {
				log.Fatal(err1)
			}
			//fmt.Println("b: ", b)
			cor := color.RGBA{b[0], b[1], b[2], 255}
			data = append(data, cor)

		}
	} else {
		cores_padrao = append(cores_padrao, color.Black)
		cores_padrao = append(cores_padrao, color.RGBA{255, 200, 3, 255})
		cores_padrao = append(cores_padrao, color.RGBA{160, 1, 34, 255})
		cores_padrao = append(cores_padrao, color.RGBA{85, 165, 34, 255})
		cores_padrao = append(cores_padrao, color.RGBA{1, 150, 255, 255})
		cores_padrao = append(cores_padrao, color.RGBA{85, 2, 255, 255})
		cores_padrao = append(cores_padrao, color.White)
		cores_padrao = append(cores_padrao, color.RGBA{100, 80, 15, 255})
		cores_padrao = append(cores_padrao, color.RGBA{10, 200, 180, 255})
		cores_padrao = append(cores_padrao, color.RGBA{160, 1, 34, 255})

		for i, cor := range cores_padrao {
			if i < nc {
				data = append(data, cor)
			}
		}
	}
	//linhas := models.LinhaBordadoPageList(id)
	return data
}

// *
// *
// ***************** Corres **************************
func CarregaCores(id int, nc int) []*models.Linha {
	cores_padrao := []string{"5075","5208","5151","5115","5310","5311","5158","5058","5027","5198"}
	data := []*models.Linha{}

	linhas := models.LinhaBordadoPageList(id)
	if len(linhas) > 0 {
		for _, linha := range linhas {
			l, err := models.LinhaOne(linha.Linha.Codigo)
			if err != nil {
				log.Fatal(err)
			}
			
			colorStr, err := normalize(l.CorHex)
			if err != nil {
				log.Fatal(err)
			}
			b, err := hex.DecodeString(colorStr)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("b: ", b)
			cor := color.RGBA{b[0], b[1], b[2], 255}
			l.CorRGB = cor
			data = append(data, l)
		}
	
	} else {
		for i, cor := range cores_padrao {
			l, err := models.LinhaOne(cor)
			if i < nc && err == nil {
				colorStr, err := normalize(l.CorHex)
				if err != nil {
					log.Fatal(err)
				}
				b, err1 := hex.DecodeString(colorStr)
				if err1 != nil {
					log.Fatal(err1)
				}
				//fmt.Println("b: ", b)
				cor := color.RGBA{b[0], b[1], b[2], 255}
				l.CorRGB = cor
				data = append(data, l)
			}
		}
	}

	return data

}

// *
// *
// ***************** LerDst **************************
func (c *BordadoController) LerDst() {
	cod_linha := c.GetString("cor")
	id, _ := c.GetInt("id", 0) //id do bordado
	seq := c.GetString("seq")
	arq := c.GetString("file")

	fmt.Printf("\n\nPARAMS: cor:%s id:%d seq:%s arq:%s\n", cod_linha,id,seq,arq)

	mCor := 0
	var corHex = ""

	//fileName := c.GetString("arq")

	var imgRect = image.Rect(0, 0, 300, 200)
	var img = image.NewRGBA(imgRect)

	fmt.Printf("\n\nAntes\n")

	if corHex == "" {
		if l, err := models.LinhaOne(cod_linha); err == nil {
			corHex = l.CorHex
		} else {
			corHex =  "f6f6f6"
		}
	}

	fmt.Printf("\n\ncorHex:%s\n", corHex)
	
/* 	if arq == "" {
		
		if nSeq, err := strconv.Atoi(seq); err == nil {
			cores_utilizada[nSeq],_ = models.ParseHexColor(corHex)
			fmt.Printf("\n\nCOR: seq=%d %s", nSeq, cores_utilizada[nSeq])
		}
		//ajustar cor da seq
		
	} 
 */

	data, err := ioutil.ReadFile(arq)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	//fmt.Println("data: ", data)

	salto := false

	if id > 0 {

		/* fmt.Printf("     Header1: %s %s %s %s %s %s\n", 
			strings.TrimSpace(string(data[41:46])), 
			strings.TrimSpace(string(data[50:55])), 
			strings.TrimSpace(string(data[59:64])),
			strings.TrimSpace(string(data[68:73])),
			strings.TrimSpace(string(data[23:30])),
			strings.TrimSpace(string(data[34:37])) ) */

		Xmais, _ := strconv.Atoi(strings.TrimSpace(string(data[41:46])))
		Xmenos, _ := strconv.Atoi(strings.TrimSpace(string(data[50:55])))
		Ymais, _ := strconv.Atoi(strings.TrimSpace(string(data[59:64])))
		Ymenos, _ := strconv.Atoi(strings.TrimSpace(string(data[68:73])))

		Largura := Xmais + Xmenos
		Altura := Ymais + Ymenos

		NrPontos, _ := strconv.Atoi(strings.TrimSpace(string(data[23:30])))
		Cores, _ := strconv.Atoi(strings.TrimSpace(string(data[34:37])))
		Cores++

		cores_padrao := CarregaCores(0, Cores + 1)
		cores_utilizada := cores_padrao
		if (id > 0) {
			cores_utilizada = CarregaCores(id, Cores)
		}
	
		if len(cores_utilizada) < 1 {
			cores_utilizada = cores_padrao
		}
	
		fmt.Printf("\n\ncores_utilizada: %d\n", len(cores_utilizada))
	
		X0 := Xmenos
		Y0 := Ymenos + 20

		zoom := 40
		if (Largura > Altura) && Largura != 0 {
			zoom = 20000 / Largura
		} else if Altura != 0 {
			zoom = 20000 / Altura
		}

		fmt.Printf("     Header2: %d %d %d %d %d %d %d %d zoom:%d len(cores):%d\n", Xmais, Xmenos, Ymais, Ymenos, Largura, Altura, NrPontos, Cores, 
					zoom, len(cores_utilizada))

		X := X0
		Y := Y0

		for i := 512; i < binary.Size(data)-3; i += 3 {
			r1 := data[i]
			r2 := data[i+1]
			r3 := data[i+2]
			//fmt.Printf("\nbyte :%d  (%d %d %d)", i, r1, r2, r3)
			if (r3 & 64) == 64 { //troca de cor
				mCor += 1
				if mCor >= len(cores_utilizada) {
					fmt.Printf("\nMuitas trocas de cor :%d    =   %d", mCor, len(cores_utilizada))
					cores_utilizada = append(cores_utilizada, cores_padrao[mCor])
				}
				fmt.Printf("\nTroca de cor :%d    =   %d", mCor, len(cores_utilizada))
			}
			salto = false
			if (r3 & 128) == 128 {
				salto = true
			}
			if (r3 & 4) == 4 {
				X += 81
			}
			if (r3 & 8) == 8 {
				X -= 81
			}
			if (r3 & 16) == 16 {
				Y += 81
			}
			if (r3 & 32) == 32 {
				Y -= 81
			}

			if (r2 & 1) == 1 {
				X += 3
			}
			if (r2 & 2) == 2 {
				X -= 3
			}
			if (r2 & 4) == 4 {
				X += 27
			}
			if (r2 & 8) == 8 {
				X -= 27
			}
			if (r2 & 16) == 16 {
				Y += 27
			}
			if (r2 & 32) == 32 {
				Y -= 27
			}
			if (r2 & 64) == 64 {
				Y += 3
			}
			if (r2 & 128) == 128 {
				Y -= 3
			}

			if (r1 & 1) == 1 {
				X += 1
			}
			if (r1 & 2) == 2 {
				X -= 1
			}
			if (r1 & 4) == 4 {
				X += 9
			}
			if (r1 & 8) == 8 {
				X -= 9
			}
			if (r1 & 16) == 16 {
				Y += 9
			}
			if (r1 & 32) == 32 {
				Y -= 9
			}
			if (r1 & 64) == 64 {
				Y += 1
			}
			if (r1 & 128) == 128 {
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
				bresenham.DrawLine(img, XX0, YY0, XX, YY, cores_utilizada[mCor].CorRGB)
			}
			X0 = X
			Y0 = Y
		}

		//Ajusta as linhas
/* 		linhas := make([]*models.LinhaBordadoRel, 0)
		for i := len(m.LinhaBordadoRel); i < int(Cores); i++ {
			linha, err := models.LinhaOne(lps[i])
			if err != nil {
				linha, _ = models.LinhaOne("5311")
			}
			fmt.Println("linha:", linha.Codigo, " ", linha.Nome)

			item := models.LinhaBordadoRel{Bordado: m, Linha: linha, Seq: i + 1}
			linhas = append(linhas, &item)
		}
 */		// Codifica a imagem em base64
		var buf bytes.Buffer
		png.Encode(&buf, img)
		b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

		row := make(map[string]interface{})

		row["linhas"] = cores_utilizada
		row["largura"] = Largura
		row["altura"] = Altura
		row["cores"] = Cores
		row["pontos"] = NrPontos
		row["CorHex"] = corHex
		row["img"] = b64
		c.Data["json"] = row

		c.ServeJSON()
	}

	//fmt.Println("LerDst:", cod_linha)
}

func normalize(colorStr string) (string, error) {
	// left as an exercise for the reader
	return colorStr, nil
}

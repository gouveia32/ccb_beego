package controllers

import (
	"bytes"
	"ccb_beego/enums"
	"ccb_beego/models"
	"encoding/base64"
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

	"encoding/hex"
	"log"

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
	seq, _ := c.GetInt("seq", 0)

	data, err := ioutil.ReadFile("C:/BORDADOS/flower15.DST")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	// Imprime o design
	//fmt.Println(data)
	type Stitch struct {
		flags int
		x     int
		y     int
		color int
	}
	type Color struct {
		r           uint8
		g           uint8
		b           uint8
		description string
	}
	
	type Pattern struct {
		colors            []Color
		stitches          []Stitch
		hoop              map[string]int
		lastX             int
		lastY             int
		top               int
		bottom            int
		left              int
		right             int
		currentColorIndex int
	}
	
	p := Pattern{}
		p.colors = []Color{}
		p.stitches = []Stitch{}
		p.hoop = make(map[string]int)
	
		func (p *Pattern) addColorRgb(r uint8, g uint8, b uint8, description string) {
			p.colors = append(p.colors, Color{r, g, b, description})
		}
		func (p *Pattern) addColor(c Color) {
			p.colors = append(p.colors, c)
		}
		func (p *Pattern) addStitchAbs(x int, y int, flags int, isAutoColorIndex bool) {
			if (flags & stitchTypes["end"]) == stitchTypes["end"] {
				p.calculateBoundingBox()
				p.fixColorCount()
			}
			if (flags & stitchTypes["stop"]) == stitchTypes["stop"] && len(p.stitches) == 0 {
				return
			}
			if (flags & stitchTypes["stop"]) == stitchTypes["stop"] && isAutoColorIndex {
				p.currentColorIndex += 1
			}
			p.stitches = append(p.stitches, Stitch{x, y, flags, p.currentColorIndex})
		}
		func (p *Pattern) addStitchRel(dx int, dy int, flags int, isAutoColorIndex bool) {
			p.addStitchAbs(p.lastX+dx, p.lastY+dy, flags, isAutoColorIndex)
		}
		func (p *Pattern) calculateBoundingBox() {
			for _, s := range p.stitches {
				if s.x < p.left {
					p.left = s.x
				}
				if s.x > p.right {
					p.right = s.x
				}
				if s.y < p.top {
					p.top = s.y
				}
				if s.y > p.bottom {
					p.bottom = s.y
				}
			}
		}
		func (p *Pattern) fixColorCount() {
			if p.currentColorIndex >= len(p.colors) {
				p.currentColorIndex = len(p.colors) - 1
			}
		}
		func rgbToColor(r uint8, g uint8, b uint8, description string) Color {
			return Color{r, g, b, description}
		}
		func shadeColor(c color.RGBA, percent float64) color.RGBA {
			f := 1 - percent
			return color.RGBA{uint8(float64(c.R) * f), uint8(float64(c.G) * f), uint8(float64(c.B) * f), c.A}
		}

	var imgRect = image.Rect(0, 0, 300, 300)
	var img = image.NewRGBA(imgRect)

	if id > 0 {

		linhas := models.LinhaBordadoPageList(id)

		fmt.Println("Bord id:", id)
		fmt.Println("linhas: ", linhas)

		var pos = 0
		var cod = ""
		var corHex = ""
		for _, linha := range linhas {
			if linha.Seq == seq {
				cod = cod_linha
			} else {
				cod = linha.Linha.Codigo
			}
			fmt.Println("L:", cod)
			l, err := models.LinhaOne(cod)
			if err != nil {
				c.pageError("Linha inexistente!!e")
			}
			//fmt.Println("linha: ", l.Nome)
			colorStr, err := normalize(l.CorHex)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("colorStr: ", colorStr)

			b, err1 := hex.DecodeString(colorStr)
			if err1 != nil {
				log.Fatal(err1)
			}
			//fmt.Println("b: ", b)
			cor := color.RGBA{b[0], b[1], b[2], 255}
			//fmt.Println("hex: ", cor)

			if linha.Seq == seq {
				corHex = colorStr
				fmt.Println("corHex: ", colorStr)
			}
			// draw line
			bresenham.DrawLine(img, 14, 21, 241+pos, 117+pos, cor)
			pos = pos + 40
		}

		// Codifica a imagem em base64
		var buf bytes.Buffer
		png.Encode(&buf, img)
		b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

		if corHex == "" {
			l, _ := models.LinhaOne(cod_linha)
			corHex = l.CorHex
			fmt.Println("NOVA:", corHex)
		}

		c.Data["json"] = corHex + b64

		c.ServeJSON()
	}

	fmt.Println("LerDst:", cod_linha)
}

func normalize(colorStr string) (string, error) {
	// left as an exercise for the reader
	return colorStr, nil
}

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

	//"image/draw"
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

// *
// *
// ***************** Index **************************
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

// *
// *
// ***************** DataGrid **************************
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

// *
// *
// ***************** DataList **************************
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
		item.Linha.CorRGB, err = models.ParseHexColor(linha.CorHex)
		if err != nil {
			item.Linha.CorRGB = color.RGBA{255, 1, 1, 255}
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

// *
// *
// ***************** Save **************************
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

	fmt.Println("Linhas: ", relslin)

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

// *
// *
// ***************** Delete **************************
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
// ***************** Corres **************************
func Hex2RGB(cHex string) color.Color {
	colorStr, err := normalize(cHex)
	if err != nil {
		log.Fatal(err)
	}
	b, err := hex.DecodeString(colorStr)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("b: ", b)
	return color.RGBA{b[0], b[1], b[2], 255}
}

// *
// *
// ***************** Corres **************************
func CarregaCores(CoresInicial []string, nc int) []*models.Linha {
	cores_padrao := []string{"5208", "5311", "5151", "5075", "5310", "5158", "5045", "5058", "5208", "5311", "5151", "5075", "5310", "5158", "5045", "5058", "5208", "5311", "5151", "5075", "5310", "5158", "5045", "5058", "5208", "5311", "5151", "5075", "5310", "5158", "5045", "5058"}
	data := []*models.Linha{}

	if CoresInicial != nil {
		for i, cor := range CoresInicial {
			l, err := models.LinhaOne(cor)
			if i < nc && err == nil {
				l.CorRGB = Hex2RGB(l.CorHex)
				data = append(data, l)
			}
		}
	} else {
		for i, cor := range cores_padrao {
			l, err := models.LinhaOne(cor)
			if i < nc && err == nil {
				l.CorRGB = Hex2RGB(l.CorHex)
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
	seq, _ := c.GetInt("seq", 0)
	arq := c.GetString("file")
	sLinhas := c.GetString("linhas")
	aLinhas := strings.Split(sLinhas, ",")

	fmt.Printf("\n\nPARAMS: cor:%s id:%d seq:%d arq:%s \nLinhas:%s\n", cod_linha, id, seq, arq, aLinhas)

	mCor := 1
	var corHex = ""

	var imgRect = image.Rect(0, 0, 300, 200)
	var img = image.NewRGBA(imgRect)

	//fmt.Printf("\n\nAntes\n")

	if corHex == "" {
		if l, err := models.LinhaOne(cod_linha); err == nil {
			corHex = l.CorHex
		} else {
			corHex = "f6f6f6"
		}
	}

	//fmt.Printf("\n\ncorHex:%s\n", corHex)

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

		cores_padrao := CarregaCores(nil, Cores+1)
		cores_utilizada := CarregaCores(aLinhas, Cores+1)
		if len(cores_utilizada) < 1 {
			cores_utilizada = cores_padrao
		}

		for _, linha := range cores_utilizada {
			fmt.Printf("\nAntes:%s  %s ", linha.Codigo, linha.Nome)

		}

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
				if mCor > len(cores_utilizada) {
					//fmt.Printf("\nbyte :%d  (%d %d %d)\n", i, r1, r2, r3)
					cores_utilizada = append(cores_utilizada, cores_padrao[mCor-1])
					//fmt.Printf("\nMuitas trocas de cor :%d    =   %d", mCor, len(cores_utilizada))
				}
				//fmt.Printf("\nTroca de cor :%d    =   %d\n", mCor, len(cores_utilizada))
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
				bresenham.DrawLine(img, XX0, YY0, XX, YY, cores_utilizada[mCor-1].CorRGB)
			}
			X0 = X
			Y0 = Y
		}

		for _, linha := range cores_utilizada {
			fmt.Printf("\n%s  %s", linha.Codigo, linha.Nome)
		}

		// Codifica a imagem em base64
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

/* func DrawLine(g *image.RGBA, x1, y1, x2, y2, scaleXY int, currentColor color.Color, flags int){


	const fNORMAL = 0
	const fJUMP = 1
	const fTRIM1 = 2
	const fTRIM = 2
	const fSTOP = 4
	const fEND = 8
	const colorchnge = 16
	const colordst = 5
	const SEQUIN_MODE = 6
	const SEQUIN_MODE = 6
	const SEQUIN_EJECT = 7
	const SEQUIN_MODE = 6
	const SLOW = &HB
	const NEEDLE_SET = 9
	const FAST = &HC

	r1 := 0
	g1 := 0
	b1 := 0
	//  Color color = Color.FromArgb(255, 0, 0);
	//  int dark = adjust(Convert.ToInt32(color), -60);
	somecolor := currentColor
	amt := Math.Round(2.55 * 2)
	str := somecolor.R.ToString() & "," & somecolor.G.ToString() & "," & somecolor.B.ToString()
	r1 = Convert.ToInt16(somecolor.R.ToString())
	g1 = Convert.ToInt16(somecolor.G.ToString())
	b1 = Convert.ToInt16(somecolor.B.ToString())
	r1 = r1 + CInt(amt)
	g1 = g1 + CInt(amt)
	b1 = b1 + CInt(amt)
	if r1 > 250 {
		r1 = 248
	}
	if g1 > 250 {
		g1 = 248
	}
	if b1 > 250 {
		b1 = 250
	}
	if r1 == 5 && g1 == 5 && b1 == 5 {
		r1 = 5
		g1 = 30
		b1 = 10
	}

	currentColor1 := color.RGBA{r1, g1, b1, 255}

	x1Scaled := int(x1 / scaleXY)
	y1Scaled := int(y1 / scaleXY)
	x2Scaled := int(x2 / scaleXY)
	y2Scaled := int(y2 / scaleXY)
	if x1Scaled == x2Scaled && y1Scaled == y2Scaled {
		Return
	}
	corR := 0
	corg := 0
	corb := 0
	if somecolor.R == 0 && somecolor.G == 0 && somecolor.B == 0 {
		corR = adjust(35, 0)
		corg = adjust(35, 0)
		corb = adjust(35, 0)
	} else {
		corR = adjust(somecolor.R, -60)
		corg = adjust(somecolor.G, -60)
		corb = adjust(somecolor.B, -60)
	}

	coresmodificadas := Color.FromArgb(corR, corg, corb, 255)

	Dim ppp As LinearGradientBrush = New LinearGradientBrush(New PointF(x2Scaled, y2Scaled), New PointF(x1Scaled, y1Scaled), coresmodificadas, currentColor)
	g.SmoothingMode = SmoothingMode.AntiAlias
	Dim myColors As Color() = {currentColor, coresmodificadas, currentColor, coresmodificadas, currentColor}
	Dim posArray As Single() = New Single() {0, 0.05F, 0.5F, 0.9F, 1.0F}
	Dim myBlend As ColorBlend = New ColorBlend()
	myBlend.Colors = myColors
	myBlend.Positions = posArray
	ppp.InterpolationColors = myBlend
	Dim pt As Pen = New Pen(ppp, 3.2F)
	pt.DashCap = DashCap.Round
	If flags = CShort(StitchType.SEQUIN_MODE) Then
		sequin_mode = True
	End If
	If flags = CShort(StitchType.fJUMP) Then
		pt = New Pen(Color.Transparent, 1.0F)
		g.DrawLine(pt, x2Scaled, y2Scaled, x1Scaled, y1Scaled)
		sequin_mode = False
	End If
	If sequin_mode = True Then
		If flags = CShort(StitchType.SEQUIN_EJECT) Then
			g.FillEllipse(Brushes.Red, Convert.ToInt32(x2Scaled) - 20, Convert.ToInt32(y2Scaled) - 20, 40, 40)
			g.FillEllipse(Brushes.White, Convert.ToInt32(x2Scaled) - 6, Convert.ToInt32(y2Scaled) - 6, 12, 12)
		End If
	End If
	If flags = CShort(StitchType.fTRIM) Then
		pt = New Pen(Color.Transparent, 1.0F)
		g.DrawLine(pt, x2Scaled, y2Scaled, x1Scaled, y1Scaled)
	End If
	If contar1 >= 10 Then
		g.DrawLine(pt, x2Scaled, y2Scaled, x1Scaled, y1Scaled)
	End If
	contar1 += 1

}

func adjust(color int, amount int) (ret int) {
	a := color >> 24 && &HFF
	r := color >> 16 && &HFF
	g := color >> 8 && &HFF
	b := color && &HFF
	r = clamp(r + amount)
	g = clamp(g + amount)
	b = clamp(b + amount)
	return a << 24 || r << 16 || g << 8 || b
}

fun clamp(v int) (ret int) {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return v
}

*/

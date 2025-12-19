package services

import (
	"Go-API-T/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type Request struct {
	Name           string
	Logo           string
	DescriptionApp string
}

func CreatePDF(RequestData Request, testData []models.Saveendpointresult, c *gin.Context) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)
	pdf.AliasNbPages("")

	//header
	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Arial", "B", 12)
		pdf.MultiCell(0, 10, "Por ahora texto, despues logo", "", "L", false)
		y := pdf.GetY()
		pdf.Line(10, y, 200, y)

	})

	// Footer with number page
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		y := pdf.GetY()
		pdf.Line(10, y, 200, y)
		pdf.AddUTF8Font("Arial", "", "./font/Arial.ttf")
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.SetFont("Arial", "", 14)

	type Section struct {
		Title string
		Page  int
		Level int
	}
	sections := []Section{}
	//Introductory page
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 24)
	pdf.Ln(10)
	//pdf.Image(RequestData.Logo, 10, 100, 50, 0, false, "", 0, "") //para agregar una imagen subida a algun almacenamiento en la nube tengo que hacer un buffer y descargar la imagen, para luego colocarla... mas adelante
	pdf.MultiCell(0, 12, RequestData.Name, "", "L", false)
	pdf.Ln(10)
	y := pdf.GetY()
	pdf.Line(10, y, 200, y)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8, RequestData.DescriptionApp, "", "L", false)
	pdf.Ln(15)
	pdf.SetPage(1)

	// Page of index
	pdf.AddPage()

	// pages with contains (tests)
	for i, item := range testData {
		title := item.Backendtests.Name

		pdf.AddPage()
		pdf.Bookmark(title, 0, 0)

		// Title
		pdf.SetFont("Arial", "B", 16)
		pdf.MultiCell(0, 10, title, "", "L", false)

		// Contain
		pdf.SetFont("Arial", "", 12)
		text := fmt.Sprintf(
			"Description: %s\n\nURL: %s",
			item.Testcasedescription,
			item.Backendtests.Urlapi,
		)
		pdf.MultiCell(0, 8, text, "", "L", false)
		//table

		pdf.SetFont("Arial", "", 12)

		data := map[string]string{
			"ID":          fmt.Sprintf("TCP-%d", i),
			"HTTP Type":   item.Backendtests.Httptype,
			"Request":     formatJSON(item.Backendtests.Request.String()),
			"Header":      formatJSON(item.Backendtests.Header.String()),
			"Response":    formatJSON(item.Backendtests.Response.String()),
			"Token":       item.Backendtests.Token,
			"Http Status": item.Backendtests.Requesttype,
		}

		col1Width := 50.0
		col2Width := 130.0
		lineHeight := 8.0
		pdf.SetFillColor(230, 230, 230)

		keys := []string{"ID", "HTTP Type", "Request", "Header", "Response", "Token", "Http Status"}

		for _, key := range keys {
			value := data[key]

			x := pdf.GetX()
			y := pdf.GetY()

			// Calcular número de líneas de cada columna
			nLinesKey := maxLines(key, col1Width)
			nLinesValue := maxLines(value, col2Width)

			// Altura máxima de la fila
			height := float64(max(nLinesKey, nLinesValue)) * lineHeight

			// Key (nombre) con altura máxima
			pdf.SetXY(x, y)
			pdf.SetFont("Arial", "B", 12)
			pdf.MultiCell(col1Width, height/float64(nLinesKey), key, "1", "L", true)

			// Value con altura máxima
			pdf.SetXY(x+col1Width, y)
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(col2Width, height/float64(nLinesValue), value, "1", "L", false)

			// Mover cursor a la siguiente fila
			pdf.SetY(y + height)
		}

		sections = append(sections, Section{
			Title: title,
			Page:  pdf.PageNo(),
			Level: 0,
		})
	}

	// Save the last page for security
	last := pdf.PageNo()

	// Write index in the first page
	pdf.SetPage(2)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(10, 30)
	pdf.MultiCell(0, 10, "Index", "", "L", false)
	pdf.SetFont("Arial", "", 12)
	for _, sec := range sections {
		entry := fmt.Sprintf("%s ....................................... %d", sec.Title, sec.Page)
		pdf.MultiCell(0, 8, entry, "", "L", false)
		pdf.Ln(8)
	}

	// Return to the last page to close the document with the proper context.
	pdf.SetPage(last)

	if err := pdf.OutputFileAndClose("documento_con_indice.pdf"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func maxLines(text string, lineHeight float64) int {
	return len(strings.Split(text, "\n"))
}

func formatJSON(raw string) string {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(raw), "", "  ")
	if err != nil {
		return raw
	}
	return pretty.String()
}

package services

import (
	"Go-API-T/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)


func CreatePDF(jsonData []models.Saveendpointresult, c *gin.Context) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AliasNbPages("")

	// Footer with number page
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
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
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "BANN")
	pdf.Ln(15)
	pdf.SetPage(1)

	// Page of index
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Index")
	pdf.Ln(15)

	// pages with contains (tests)
	for _, item := range jsonData {
		title := item.Backendtests.Name

		pdf.AddPage()
		pdf.Bookmark(title, 0, 0)

		// Title
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, title)
		pdf.Ln(12)

		// Contain
		pdf.SetFont("Arial", "", 12)
		text := fmt.Sprintf(
			"Name of the test: %s\n\nDescription: %s\n\nURL: %s",
			item.Backendtests.Name,
			item.Testcasedescription,
			item.Backendtests.Urlapi,
		)
		pdf.MultiCell(0, 8, text, "", "L", false)

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
	pdf.SetXY(15, 20)
	pdf.Cell(0, 10, "Index")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 12)
	for _, sec := range sections {
		entry := fmt.Sprintf("%s .......... %d", sec.Title, sec.Page)
		pdf.Cell(0, 8, entry)
		pdf.Ln(8)
	}

	// Return to the last page to close the document with the proper context.
	pdf.SetPage(last)

	if err := pdf.OutputFileAndClose("documento_con_indice.pdf"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

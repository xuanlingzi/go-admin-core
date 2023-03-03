package utils

import "github.com/jung-kurt/gofpdf"

func WritePDF(filename string, content interface{}) gofpdf.Pdf {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "in",
		Size:    gofpdf.SizeType{Wd: 6, Ht: 6},
	})
	return pdf
}

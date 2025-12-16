package excel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AgamariFF/TenderMessage.git/internal/logger"
	"github.com/AgamariFF/TenderMessage.git/internal/models"

	"github.com/xuri/excelize/v2"
)

func ToExcel(Tenders *[]models.Tender) (string, error) {
	f := excelize.NewFile()

	if err := addTendersAndSheet(f, *Tenders, "Двери"); err != nil {
		logger.SugaredLogger.Warn(err)
	}

	filename := fmt.Sprintf("/tmp/tenders_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		logger.SugaredLogger.Warnf("ошибка сохранения Excel: %w", err)
	}

	return filename, nil
}

func addTendersAndSheet(f *excelize.File, tendersZakupkiGovRu []models.Tender, sheet string) error {
	f.NewSheet(sheet)

	index, _ := f.GetSheetIndex("Sheet1")
	if index != -1 {
		f.DeleteSheet("Sheet1")
	}

	CreateHeader(f, sheet)

	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 18,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	err = f.MergeCell(sheet, "A2", "E2")
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheet, "A2", "E2", titleStyle)
	if err != nil {
		return err
	}

	f.SetCellValue(sheet, "A2", "Zakupki.Gov.ru")

	index = 3

	setTenderInf(f, sheet, tendersZakupkiGovRu, &index)

	return nil
}

func CreateHeader(f *excelize.File, sheet string) error {
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:     "center",
			Indent:         1,
			ReadingOrder:   0,
			RelativeIndent: 1,
			ShrinkToFit:    false,
			TextRotation:   0,
			Vertical:       "",
			WrapText:       true,
		},
		Font: &excelize.Font{
			Bold:      true,
			Italic:    false,
			Underline: "",
			Family:    "",
			Size:      12,
			Strike:    false,
		},
	})

	if err != nil {
		return err
	}

	f.SetColWidth(sheet, "A", "B", 16)
	f.SetColWidth(sheet, "C", "C", 34)
	f.SetColWidth(sheet, "D", "D", 40)
	f.SetColWidth(sheet, "E", "E", 100)
	f.SetColWidth(sheet, "F", "F", 20)
	f.SetCellValue(sheet, "A1", "Дата размещения")
	f.SetCellValue(sheet, "B1", "Дата окончания")
	f.SetCellValue(sheet, "C1", "Расположение")
	f.SetCellValue(sheet, "D1", "Заказчик")
	f.SetCellValue(sheet, "E1", "Объект закупки + ссылка")
	f.SetCellValue(sheet, "F1", "Начальная цена")
	f.SetCellStyle(sheet, "A1", "F1", style)
	f.SetCellValue(sheet, "G1", "Дата создания таблицы: "+time.Now().UTC().Format("02.01.2006"))

	return nil
}

func setTenderInf(f *excelize.File, sheet string, tender []models.Tender, index *int) {
	for _, value := range tender {
		f.SetCellValue(sheet, "A"+strconv.Itoa(*index), value.PublishDate)
		f.SetCellValue(sheet, "B"+strconv.Itoa(*index), value.EndDate)
		f.SetCellValue(sheet, "C"+strconv.Itoa(*index), value.Region)
		f.SetCellValue(sheet, "D"+strconv.Itoa(*index), value.Customer)
		f.SetCellValue(sheet, "E"+strconv.Itoa(*index), value.Title)
		f.SetCellValue(sheet, "F"+strconv.Itoa(*index), value.Price)

		f.SetCellHyperLink(sheet, "E"+strconv.Itoa(*index), value.Link, "External")
		*index++
	}
}

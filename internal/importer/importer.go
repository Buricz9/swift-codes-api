package importer

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
	"swift-codes-api/internal/service"
)

func ImportSwiftCodesFromXLSX(filePath string, swiftSvc service.SwiftService) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening xlsx file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("could not read rows from sheet: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 8 {
			fmt.Printf("Skipping row with fewer columns than expected: %v\n", row)
			continue
		}

		countryISO2 := row[0]
		swiftCode := row[1]
		bankName := row[3]
		address := row[4]
		townName := row[5]
		countryName := row[6]

		countryISO2 = strings.ToUpper(countryISO2)
		countryName = strings.ToUpper(countryName)

		combinedAddress := fmt.Sprintf("%s, %s", address, townName)

		isHeadquarter := false
		if strings.HasSuffix(strings.ToUpper(swiftCode), "XXX") {
			isHeadquarter = true
		}

		var headquarterSwiftCode *string
		if !isHeadquarter && len(swiftCode) >= 8 {
			hq := swiftCode[:8] + "XXX"
			headquarterSwiftCode = &hq
		}

		err := swiftSvc.CreateSwiftCode(service.CreateSwiftCodeInput{
			SwiftCode:            swiftCode,
			BankName:             bankName,
			Address:              combinedAddress,
			CountryISO2:          countryISO2,
			CountryName:          countryName,
			IsHeadquarter:        isHeadquarter,
			HeadquarterSwiftCode: headquarterSwiftCode,
		})
		if err != nil {
			return fmt.Errorf("could not import row %d, swiftCode=%s: %w", i+1, swiftCode, err)
		}
	}

	return nil
}

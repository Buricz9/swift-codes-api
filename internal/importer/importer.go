package importer

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
	"swift-codes-api/internal/service"
)

// ImportSwiftCodesFromXLSX odczytuje plik XLSX i importuje dane do bazy
// za pomocą serwisu SwiftService.
func ImportSwiftCodesFromXLSX(filePath string, swiftSvc service.SwiftService) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("error opening xlsx file: %w", err)
	}
	defer f.Close()

	// Domyślnie chcemy czytać wiersze z "Sheet1" (dostosuj, jeśli się różni)
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("could not read rows from sheet: %w", err)
	}

	// rows[0] to zwykle nagłówek, więc iterujemy od drugiego wiersza
	for i, row := range rows {
		if i == 0 {
			// pomijamy nagłówek
			continue
		}

		// Zgodnie z Twoją strukturą:
		//  0 -> COUNTRY ISO2 CODE
		//  1 -> SWIFT CODE
		//  2 -> CODE TYPE (opcjonalnie)
		//  3 -> NAME (bankName)
		//  4 -> ADDRESS
		//  5 -> TOWN NAME
		//  6 -> COUNTRY NAME
		//  7 -> TIME ZONE (opcjonalne)

		if len(row) < 8 {
			// Przykładowa akcja: pomijamy wiersz
			fmt.Printf("Skipping row with fewer columns than expected: %v\n", row)
			continue
		}

		countryISO2 := row[0]
		swiftCode := row[1]
		// codeType := row[2] // pomijamy
		bankName := row[3]
		address := row[4]
		townName := row[5]
		countryName := row[6]
		// timeZone := row[7] // pomijamy

		// Wymaganie: uppercase na polach kraju
		countryISO2 = strings.ToUpper(countryISO2)
		countryName = strings.ToUpper(countryName)

		// Możesz też złączyć address + townName, jeśli chcesz
		combinedAddress := fmt.Sprintf("%s, %s", address, townName)

		// Sprawdzamy, czy kończy się na "XXX" -> HQ
		isHeadquarter := false
		if strings.HasSuffix(strings.ToUpper(swiftCode), "XXX") {
			isHeadquarter = true
		}

		// Dla branch ustawiamy headquarterSwiftCode
		var headquarterSwiftCode *string
		if !isHeadquarter && len(swiftCode) >= 8 {
			hq := swiftCode[:8] + "XXX"
			headquarterSwiftCode = &hq
		}

		// Tworzymy rekord w bazie
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

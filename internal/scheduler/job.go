package scheduler

import (
	"regexp"

	"github.com/AgamariFF/TenderMessage.git/internal/models"
	"github.com/AgamariFF/TenderMessage.git/internal/parsergovru"
)

func Search(re *regexp.Regexp) ([]models.Tender, error) {
	config := setupConfig()
	tenders, err := parsergovru.ParseGovRu("doors", &config, re)
	if err != nil {
		return []models.Tender{}, err
	}
	return tenders, nil
}

func setupConfig() models.Config {
	config := models.Config{
		SearchVent:        false,
		SearchDoors:       true,
		SearchBuild:       false,
		SearchMetal:       false,
		MinPriceVent:      0,
		MinPriceDoors:     0,
		MinPriceBuild:     0,
		MinPriceMetal:     0,
		ProcurementType:   "active",
		VentCustomerPlace: []string{"OKER37", "OKER33", "OKER30"},
	}
	return config
}

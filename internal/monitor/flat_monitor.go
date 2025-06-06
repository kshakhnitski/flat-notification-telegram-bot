package monitor

import (
	"flat_bot/internal/model"
	"flat_bot/internal/parser"
	"flat_bot/internal/repository"
	"log"
	"time"
)

type FlatsHandler func(flats []model.Flat)

type FlatMonitor struct {
	sources         []string
	checkInterval   time.Duration
	newFlatsHandler FlatsHandler
	flatRepository  repository.FlatRepository
}

func NewFlatMonitor(
	checkInterval time.Duration,
	flatRepository repository.FlatRepository,
	newFlatsHandler FlatsHandler,
) FlatMonitor {
	return FlatMonitor{
		checkInterval:   checkInterval,
		flatRepository:  flatRepository,
		newFlatsHandler: newFlatsHandler,
	}
}

func (l FlatMonitor) Start() {
	ticker := time.NewTicker(l.checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Checking for new flats...")
		newFlats := l.checkForNewFlats()

		l.newFlatsHandler(newFlats)
	}
}

func (l FlatMonitor) checkForNewFlats() []model.Flat {
	var loadedFlats []model.Flat
	var newFlats []model.Flat

	loadedFlats = append(loadKufarFlats())

	for _, flat := range loadedFlats {
		exists, err := l.flatRepository.ExistsByID(flat.ID)
		if err != nil {
			log.Fatalf("Error while checking if flat exists: %v", err)
		}

		if exists {
			continue
		}

		newFlats = append(newFlats, flat)
		_, err = l.flatRepository.Create(flat)
		if err != nil {
			log.Fatalf("Error while saving flat: %v", err)
		}
	}

	return newFlats
}

func loadKufarFlats() []model.Flat {
	url := "https://re.kufar.by/l/minsk/snyat/kvartiru-dolgosrochno/bez-posrednikov?cur=USD&gbx=b%3A27.150276810254205%2C53.344712700318226%2C28.932808548535448%2C54.058539029097716&prc=r%3A0%2C350&rms=v.or%3A1%2C2%2C3&size=30"
	kufarFlatParser := parser.NewKufarFlatParser(url)

	flats, err := kufarFlatParser.Parse()
	if err != nil {
		log.Fatalf("Error while parsing kufar flats: %v", err)
	}

	return flats
}

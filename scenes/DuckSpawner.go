package scenes

import (
	"duckgame/models"
	"math/rand"
	"sync"
	"time"
)

type DuckSpawner struct {
	ducksChannel chan []*models.Duck
	ducks        []*models.Duck
	wg           *sync.WaitGroup
}

var newDuck *models.Duck

func NewDuckSpawner(ducksChannel chan []*models.Duck, wg *sync.WaitGroup) *DuckSpawner {
	return &DuckSpawner{
		ducksChannel: ducksChannel,
		wg:           wg,
	}
}

func (ds *DuckSpawner) Start() {
	defer ds.wg.Done()

	spawnDucks := func(count int) []*models.Duck {
		var ducks []*models.Duck
		for j := 0; j < count; j++ {
			newDuck = models.NewDuck(rand.Float64()*800, rand.Float64()*600, (rand.Float64()-0.3)*0.2, (rand.Float64()-0.3)*0.2)
			ducks = append(ducks, newDuck)
		}
		return ducks
	}

	ds.ducksChannel <- spawnDucks(5)
	time.Sleep(time.Second * 2)
	for {
		ds.ducksChannel <- spawnDucks(5)
		time.Sleep(time.Second * 4)
	}
}

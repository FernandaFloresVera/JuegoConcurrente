package views

import (
	"duckgame/models"
	"duckgame/scenes"
	"github.com/faiface/pixel/pixelgl"
	"sync"
)

type MainView struct{}

var duckSpawnerScene *scenes.DuckSpawner
var rendererScene *scenes.Renderer

func NewMainView() *MainView {
	return &MainView{}
}

func (m *MainView) Render() {
	pixelgl.Run(func() {
		var wg sync.WaitGroup

		ducksChannel := make(chan []*models.Duck, 1)
		scoreChan := make(chan int, 1)
		wg.Add(3)

		duckSpawnerScene = scenes.NewDuckSpawner(ducksChannel, &wg)
		rendererScene = scenes.NewRenderer(ducksChannel, scoreChan, &wg)

		go duckSpawnerScene.Start()
		go rendererScene.Render()
		go rendererScene.ScoreTracker()

		wg.Wait()
	})
}

package scenes

import (
	"duckgame/models"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/png"
	"math"
	"os"
	"sync"
)

type Renderer struct {
	ducksChan  chan []*models.Duck
	scoreChan  chan int
	score      int
	totalDucks int
	wg         *sync.WaitGroup
	mu         *sync.Mutex
}

func NewRenderer(ducksChan chan []*models.Duck, scoreChan chan int, wg *sync.WaitGroup) *Renderer {
	return &Renderer{
		ducksChan:  ducksChan,
		scoreChan:  scoreChan,
		score:      0,
		totalDucks: 0,
		wg:         wg,
		mu:         &sync.Mutex{},
	}
}

var ducks []*models.Duck

func (r *Renderer) Render() {
	defer r.wg.Done()

	bgPic, err := loadPicture("assets/background.png")
	if err != nil {
		panic(fmt.Sprintf("Error al cargar background.png: %v", err))
	}
	duckPic, err := loadPicture("assets/duck.png")
	if err != nil {
		panic(fmt.Sprintf("Error al cargar duck.png: %v", err))
	}

	bgSprite := pixel.NewSprite(bgPic, bgPic.Bounds())
	duckSprite := pixel.NewSprite(duckPic, duckPic.Bounds())

	cfg := pixelgl.WindowConfig{
		Title:  "Juego de patitos",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(10, 580), atlas)

	pauseButton := models.NewButton("Pausar", "assets/pause.png", 690, 540, 130, 55)
	restartButton := models.NewButton("Reiniciar", "assets/reload.png", 690, 480, 130, 55)

	paused := false

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		bgSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		pauseButton.Draw(win, atlas, 0.1)
		restartButton.Draw(win, atlas, 0.17)

		if !paused {
			select {
			case newDucks := <-r.ducksChan:
				ducks = append(ducks, newDucks...)
				r.totalDucks += len(newDucks)
			default:
			}

			for i := 0; i < len(ducks); i++ {
				ducks[i].X += ducks[i].VelocityX
				ducks[i].Y += ducks[i].VelocityY

				if ducks[i].X > win.Bounds().Max.X || ducks[i].X < win.Bounds().Min.X {
					ducks[i].VelocityX = -ducks[i].VelocityX
				}

				if ducks[i].Y > win.Bounds().Max.Y || ducks[i].Y < win.Bounds().Min.Y {
					ducks[i].VelocityY = -ducks[i].VelocityY
				}
				duckSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(pixel.V(ducks[i].X, ducks[i].Y)))
			}
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {

			if pauseButton.Contains(win.MousePosition()) {
				paused = !paused

				if paused {
					pauseButton.Text = "Reanudar"
				} else {
					pauseButton.Text = "Pausar"
				}
			} else if restartButton.Contains(win.MousePosition()) {
				ducks = nil
				r.score = 0
				r.totalDucks = 0
			} else if !paused {
				mousePos := win.MousePosition()
				r.mu.Lock()

				for i := len(ducks) - 1; i >= 0; i-- {
					duck := ducks[i]

					if Distance(mousePos.X, mousePos.Y, duck.X, duck.Y) < 30 {

						ducks = append(ducks[:i], ducks[i+1:]...)
						r.score++
						r.scoreChan <- r.score
					}
				}
				r.mu.Unlock()
			}
		}

		txt.Clear()

		if paused {
			fmt.Fprint(txt, "Juego Pausado")
		} else {
			fmt.Fprintf(txt, "Puntos: %d / %d", r.score, r.totalDucks)
		}
		txt.Draw(win, pixel.IM)
		win.Update()
	}
}

func (r *Renderer) ScoreTracker() {
	defer r.wg.Done()
	for {
		currentScore := <-r.scoreChan
		r.mu.Lock()
		r.score = currentScore
		r.mu.Unlock()
	}
}

func loadPicture(path string) (pixel.Picture, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

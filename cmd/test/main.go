package main

import (
	"github.com/therecipe/qt/widgets"
	"os"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/goqtreactor/reactor"
	"time"
	"math/rand"
	"github.com/therecipe/qt/core"
	"fmt"
)

var globalScene *widgets.QGraphicsScene

func runInMainThread(args ...interface{}) {
	i, j := args[0].(int), args[1].(int)
	fmt.Println(i, j)

	x1 := float64(rand.Int() % 400)
	x2 := float64(rand.Int() % 400)
	y1 := float64(rand.Int() % 400)
	y2 := float64(rand.Int() % 400)

	r := rand.Int() % 255
	g := rand.Int() % 255
	b := rand.Int() % 255

	pen := gui.NewQPen3(gui.NewQColor3(r, g, b, 255))
	globalScene.AddLine2(x1, y1, x2, y2, pen)
	globalScene.Update(core.NewQRectF())
}

func run(args ...interface{}) {
	fmt.Println(args)
	for {
		reactor.CallFromThread(runInMainThread, 1, 2)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	w := widgets.NewQGraphicsView(nil)
	scene := widgets.NewQGraphicsScene(w)

	w.SetGeometry2(0, 0, 400, 400)
	scene.SetSceneRect2(0, 0, 400, 400)

	w.SetScene(scene)
	w.Scale(1, -1)

	globalScene = scene

	w.Show()

	reactor.Initialize()
	reactor.CallInThread(run, "test")

	os.Exit(app.Exec())
}

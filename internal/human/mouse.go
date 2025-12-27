package human

import (
	"math/rand"
	"time"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func BezierCurve(t float64, p0, p1, p2, p3 struct{ X, Y float64 }) struct{ X, Y float64 } {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t

	x := uuu*p0.X + 3*uu*t*p1.X + 3*u*tt*p2.X + ttt*p3.X
	y := uuu*p0.Y + 3*uu*t*p1.Y + 3*u*tt*p2.Y + ttt*p3.Y

	return struct{ X, Y float64 }{X: x, Y: y}
}

func MoveMouseNaturally(page *rod.Page, element *rod.Element) {
	startPoint := struct{ X, Y float64 }{X: 100, Y: 100}

	// Ensure element is on screen
	element.MustScrollIntoView()
	
	box := element.MustShape().Box()
	targetX := box.X + (box.Width / 2)
	targetY := box.Y + (box.Height / 2)

	control1 := struct{ X, Y float64 }{
		X: startPoint.X + (rand.Float64()*100 - 50),
		Y: startPoint.Y + (rand.Float64()*100 - 50),
	}
	control2 := struct{ X, Y float64 }{
		X: targetX + (rand.Float64()*100 - 50),
		Y: targetY + (rand.Float64()*100 - 50),
	}

	target := struct{ X, Y float64 }{X: targetX, Y: targetY}

	steps := 20
	// ...existing code...
    for i := 0; i <= steps; i++ {
        t := float64(i) / float64(steps)
        point := BezierCurve(t, startPoint, control1, control2, target)
        
        page.Mouse.MoveTo(proto.Point{X: point.X, Y: point.Y})

        time.Sleep(time.Duration(rand.Intn(10)+5) * time.Millisecond)
    }
// ...existing code...
}

func TypeStealthily(element *rod.Element, text string) {
	element.MustClick()
	for _, char := range text {
		delay := time.Duration(rand.Intn(100)+50) * time.Millisecond
		time.Sleep(delay)
		element.MustInput(string(char))
	}
}
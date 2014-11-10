// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"log"
	"math"
	"os"
	"time"

	_ "image/jpeg"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/sprite"
	"golang.org/x/mobile/sprite/clock"
	"golang.org/x/mobile/sprite/glsprite"
)

var (
	start     = time.Now()
	lastClock = clock.Time(-1)

	eng   = glsprite.Engine()
	scene *sprite.Node
)

func main() {
	app.Run(app.Callbacks{
		Draw:  draw,
		Touch: touch,
	})
}

func draw() {
	if scene == nil {
		loadScene()
	}

	now := clock.Time(time.Since(start) * 60 / time.Second)
	if now == lastClock {
		// TODO: figure out how to limit draw callbacks to 60Hz instead of
		// burning the CPU as fast as possible.
		// TODO: (relatedly??) sync to vblank?
		return
	}
	lastClock = now

	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	eng.Render(scene, now)
	debug.DrawFPS()
}

func touch(t event.Touch) {
}

func newNode() *sprite.Node {
	n := &sprite.Node{}
	eng.Register(n)
	scene.AppendChild(n)
	return n
}

func loadScene() {
	texs := loadTextures()
	scene = &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, 0, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	var n *sprite.Node

	n = newNode()
	eng.SetTexture(n, 0, texs[texBooks])
	eng.SetTransform(n, 0, f32.Affine{
		{36, 0, 0},
		{0, 36, 0},
	})

	n = newNode()
	eng.SetTexture(n, 0, texs[texFire])
	eng.SetTransform(n, 0, f32.Affine{
		{72, 0, 144},
		{0, 72, 144},
	})

	n = newNode()
	n.Arranger = arrangerFunc(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
		// TODO: use a tweening library instead of manually arranging.
		t0 := uint32(t) % 120
		if t0 < 60 {
			eng.SetTexture(n, t, texs[texGopherR])
		} else {
			eng.SetTexture(n, t, texs[texGopherL])
		}

		u := float32(t0) / 120
		u = (1 - f32.Cos(u*2*math.Pi)) / 2

		tx := 18 + u*48
		ty := 36 + u*108
		sx := 36 + u*36
		sy := 36 + u*36
		eng.SetTransform(n, t, f32.Affine{
			{sx, 0, tx},
			{0, sy, ty},
		})
	})
}

const (
	texBooks = iota
	texFire
	texGopherR
	texGopherL
)

func loadTextures() (ret []sprite.Texture) {
	f, err := os.Open("../../testdata/waza-gophers.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	sheet, err := eng.LoadSheet(img)
	if err != nil {
		log.Fatal(err)
	}

	bounds := []image.Rectangle{
		texBooks:   image.Rect(4, 71, 132, 182),
		texFire:    image.Rect(330, 56, 440, 155),
		texGopherR: image.Rect(152, 10, 152+140, 10+90),
		texGopherL: image.Rect(162, 120, 162+140, 120+90),
	}

	for _, b := range bounds {
		tex, err := eng.LoadTexture(sheet, b)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, tex)
	}
	return ret
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }
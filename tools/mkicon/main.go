// Command mkicon generates the Windows app icon from the ship sprite.
//
//	go run ./tools/mkicon -out build/windows
//
// It writes icon.ico (sizes 16–256, PNG-encoded entries) and icon-256.png.
// Each size is drawn with whole-number pixel scaling so the art stays crisp.
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"retro-shooter/assets"
)

var sizes = []int{16, 24, 32, 48, 64, 128, 256}

func main() {
	out := flag.String("out", "build/windows", "output directory")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	ship := rotateLeft(assets.SpriteRGBA(assets.ShipSprite)) // nose up
	var pngs [][]byte
	for _, n := range sizes {
		var buf bytes.Buffer
		img := drawIcon(ship, n)
		if err := png.Encode(&buf, img); err != nil {
			log.Fatal(err)
		}
		pngs = append(pngs, buf.Bytes())
		if n == 256 {
			write(filepath.Join(*out, "icon-256.png"), buf.Bytes())
		}
	}
	write(filepath.Join(*out, "icon.ico"), encodeICO(sizes, pngs))
}

// drawIcon draws the sprite centered on a dark rounded tile of size n×n.
func drawIcon(sprite *image.RGBA, n int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, n, n))
	bg := color.RGBA{0x0b, 0x0d, 0x1a, 0xff}
	border := assets.Color('k')
	radius := max(2, n/6)
	bw := max(1, n/32) // border width
	for y := range n {
		for x := range n {
			if !insideRounded(x, y, n, radius) {
				continue
			}
			c := bg
			if !insideRounded(x-bw, y-bw, n-2*bw, max(1, radius-bw)) || x < bw || y < bw {
				c = border
			}
			img.SetRGBA(x, y, c)
		}
	}

	sw, sh := sprite.Bounds().Dx(), sprite.Bounds().Dy()
	scale := max(1, int(float64(n)*0.72)/max(sw, sh))
	ox := (n - sw*scale) / 2
	oy := (n - sh*scale) / 2
	for y := range sh {
		for x := range sw {
			c := sprite.RGBAAt(x, y)
			if c.A == 0 {
				continue
			}
			for dy := range scale {
				for dx := range scale {
					img.SetRGBA(ox+x*scale+dx, oy+y*scale+dy, c)
				}
			}
		}
	}
	return img
}

// insideRounded reports whether (x, y) is inside an n×n square with corner radius r.
func insideRounded(x, y, n, r int) bool {
	if x < 0 || y < 0 || x >= n || y >= n {
		return false
	}
	cx := clampInt(x, r, n-1-r)
	cy := clampInt(y, r, n-1-r)
	dx, dy := x-cx, y-cy
	return dx*dx+dy*dy <= r*r
}

func clampInt(v, lo, hi int) int { return max(lo, min(hi, v)) }

// rotateLeft rotates an image 90° counter-clockwise (right-pointing → up).
func rotateLeft(src *image.RGBA) *image.RGBA {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, h, w))
	for y := range h {
		for x := range w {
			dst.SetRGBA(y, w-1-x, src.RGBAAt(x, y))
		}
	}
	return dst
}

// encodeICO packs PNG images into a .ico file (supported since Windows Vista).
func encodeICO(sizes []int, pngs [][]byte) []byte {
	var b bytes.Buffer
	le := binary.LittleEndian
	binary.Write(&b, le, [3]uint16{0, 1, uint16(len(pngs))}) // reserved, type=icon, count
	offset := 6 + 16*len(pngs)
	for i, data := range pngs {
		dim := byte(sizes[i]) // 256 is stored as 0
		binary.Write(&b, le, struct {
			W, H, Colors, Reserved byte
			Planes, BPP            uint16
			Size, Offset           uint32
		}{dim, dim, 0, 0, 1, 32, uint32(len(data)), uint32(offset)})
		offset += len(data)
	}
	for _, data := range pngs {
		b.Write(data)
	}
	return b.Bytes()
}

func write(path string, data []byte) {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s (%d bytes)", path, len(data))
}

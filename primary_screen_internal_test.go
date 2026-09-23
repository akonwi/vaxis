package vaxis

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrimaryFrameNativeGraphicsAndMouseShape(t *testing.T) {
	vx, out := pixelTestVaxis(t)
	vx.winSize.Rows = 6
	vx.primaryScreen = &primaryScreen{regionHeight: 1, positioned: true, origin: 3}
	img, err := vx.NewPixelImage(2, 1, []byte{255, 0, 0, 255, 0, 255, 0, 255})
	if err != nil {
		t.Fatal(err)
	}
	img.Draw(placementForTest())
	vx.SetMouseShape(MouseShapeClickable)
	if err := vx.RenderFrame(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "\x1b[4;1H\x1b_Ga=t") {
		t.Fatalf("image was not placed at physical row 4: %q", out.String())
	}
	if !strings.Contains(out.String(), tparm(mouseShape, MouseShapeClickable)) {
		t.Fatalf("missing mouse-shape update: %q", out.String())
	}
	out.Reset()
	vx.exitPrimaryScreen()
	if !strings.Contains(out.String(), "\x1b_Ga=d,d=i,") || len(vx.graphicsLast) != 0 {
		t.Fatalf("image placement was not released: %q", out.String())
	}
}

func TestPrimaryFrameSixelEraseUsesSurfaceOrigin(t *testing.T) {
	vx, out := pixelTestVaxis(t)
	vx.primaryScreen = &primaryScreen{checked: true, positioned: true, origin: 3}
	sixel := &Sixel{vx: vx, w: 1, h: 1, buf: bytes.NewBufferString("encoded sixel")}
	sixel.Draw(vx.Window())
	vx.graphicsNext[0].deleteFn(out)
	if got := out.String(); got != "\x1b[4;1H " {
		t.Fatalf("sixel erase=%q, want physical row 4", got)
	}
	out.Reset()
	vx.primaryScreen.positioned = false
	vx.graphicsNext[0].deleteFn(out)
	if out.Len() != 0 {
		t.Fatalf("unpositioned sixel erased terminal contents: %q", out.String())
	}
}

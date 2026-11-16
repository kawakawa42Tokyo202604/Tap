package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestDecodeOKAndRejectBadResponses(t *testing.T) {
	var value statusView
	if err := decodeOK(`OK {"hp":8,"max_hp":10,"status":"alive"}`, &value); err != nil {
		t.Fatal(err)
	}
	if value.HP != 8 || value.MaxHP != 10 || value.Status != "alive" {
		t.Fatalf("unexpected response: %+v", value)
	}
	for _, input := range []string{"ERR unavailable", "OK {", ""} {
		if decodeOK(input, &value) == nil {
			t.Errorf("accepted malformed response %q", input)
		}
	}
}

func TestCatalogAndEmbeddedArt(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Rooms) < 40 || len(catalog.NPCs) == 0 {
		t.Fatal("world catalog is empty or incomplete")
	}
	image := loadArt("rooms", "unknown")
	if image == nil || image.Bounds().Dx() == 0 || image.Bounds().Dy() == 0 {
		t.Fatal("embedded fallback art cannot be decoded")
	}
}

func TestStatBarOnHeadlessCanvas(t *testing.T) {
	application := test.NewApp()
	defer application.Quit()
	bar := newStatBar()
	for _, sample := range []struct {
		current, max int
		want         float32
	}{{5, 10, .5}, {-1, 10, 0}, {15, 10, 1}, {1, 0, 0}} {
		bar.Set(sample.current, sample.max)
		if bar.layout.ratio != sample.want {
			t.Errorf("bar(%d,%d): %v", sample.current, sample.max, bar.layout.ratio)
		}
	}
}

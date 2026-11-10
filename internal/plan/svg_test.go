package plan

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestSVGDrawsWarehouseOutlineWithPlanDimensions(t *testing.T) {
	svg := Plan{WarehouseCode: "LYS-01", Width: 400, Height: 200}.SVG()

	for _, want := range []string{`viewBox="0 0 400 200"`, `<rect x="0" y="0" width="400" height="200"`} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG() does not contain %q:\n%s", want, svg)
		}
	}
}

func TestSVGDrawsOneLabelledRectanglePerZone(t *testing.T) {
	tests := []struct {
		name  string
		zones []Zone
	}{
		{"no zone", nil},
		{"one zone", []Zone{{Label: "Quai", X: 10, Y: 20, Width: 100, Height: 50}}},
		{"two zones", []Zone{
			{Label: "Quai", X: 10, Y: 20, Width: 100, Height: 50},
			{Label: "Stock", X: 120, Y: 20, Width: 200, Height: 150},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svg := Plan{Width: 400, Height: 200, Zones: tt.zones}.SVG()

			if got, want := strings.Count(svg, "<rect"), 1+len(tt.zones); got != want {
				t.Errorf("rectangles = %d, want %d", got, want)
			}
			for _, z := range tt.zones {
				if !strings.Contains(svg, ">"+z.Label+"</text>") {
					t.Errorf("label %q missing:\n%s", z.Label, svg)
				}
			}
		})
	}
}

func TestSVGCentersLabelInsideItsZone(t *testing.T) {
	svg := Plan{Width: 400, Height: 200, Zones: []Zone{{Label: "Quai", X: 10, Y: 20, Width: 100, Height: 50}}}.SVG()

	if want := `<text x="60" y="45"`; !strings.Contains(svg, want) {
		t.Errorf("SVG() does not contain %q:\n%s", want, svg)
	}
}

func TestSVGEscapesMarkupInZoneLabels(t *testing.T) {
	svg := Plan{Width: 100, Height: 100, Zones: []Zone{{Label: `<b>&"Quai"</b>`, Width: 10, Height: 10}}}.SVG()

	if strings.Contains(svg, "<b>") {
		t.Errorf("label is not escaped:\n%s", svg)
	}
	if err := xml.Unmarshal([]byte(svg), new(struct{ XMLName xml.Name })); err != nil {
		t.Errorf("SVG() is not well-formed XML: %v\n%s", err, svg)
	}
}

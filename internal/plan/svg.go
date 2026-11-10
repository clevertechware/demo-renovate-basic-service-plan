package plan

import (
	"fmt"
	"html"
	"strings"
)

func (p Plan) SVG() string {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`+"\n",
		p.Width, p.Height, p.Width, p.Height)
	fmt.Fprintf(&b, `  <rect x="0" y="0" width="%d" height="%d" fill="#f4f4f5" stroke="#18181b" stroke-width="2"/>`+"\n",
		p.Width, p.Height)
	for _, z := range p.Zones {
		fmt.Fprintf(&b, `  <rect x="%d" y="%d" width="%d" height="%d" fill="#dbeafe" stroke="#1d4ed8"/>`+"\n",
			z.X, z.Y, z.Width, z.Height)
		fmt.Fprintf(&b, `  <text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" font-size="14">%s</text>`+"\n",
			z.X+z.Width/2, z.Y+z.Height/2, html.EscapeString(z.Label))
	}
	b.WriteString("</svg>\n")
	return b.String()
}

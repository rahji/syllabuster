package scale

import (
	"bytes"
	"fmt"
	"math"

	"github.com/rahji/syllabuster/config"
)

/*
returns a markdown table containing a letter grade scale that is
based on an existing scale and a new maximum point value
*/
func Rescale(scale []config.Scale, points float64) string {
	buffer := &bytes.Buffer{}

	buffer.WriteString(`
# Letter Grade Scale

| Letter | Low | High |
| :----- | :-- | :--- |
`)
	for i, g := range scale {

		thismin := math.Round((float64(points) * g.Min) / 100)

		if i == 0 {
			fmt.Fprintf(buffer, "| %s | ≥ %.0f | ≤ %.0f |\n",
				g.Letter,
				thismin,
				points,
			)
			continue
		}

		thismax := math.Round((float64(points) * scale[i-1].Min) / 100)

		fmt.Fprintf(buffer, "| %s | ≥ %.0f | < %.0f |\n",
			g.Letter,
			thismin,
			thismax,
		)
	}

	return buffer.String()
}

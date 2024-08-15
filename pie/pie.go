package pie

import (
	"os"

	"github.com/wcharczuk/go-chart/v2"
)

// takes a filename and a slice of chart.Value to create a pie chart PNG
func Draw(fn string, vals []chart.Value) error {
	pie := chart.PieChart{
		Width:  700,
		Height: 700,
		Values: vals,
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	pie.Render(chart.PNG, f)
	return nil
}

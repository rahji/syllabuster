package assignment

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2"
)

type Assignment struct {
	TotalPoints float64 // eg: 300
	PointString string  // eg: "100 x 3"
	Name        string  // eg: midterm exam
}

// includes a slice of Assignment and the buffer := &bytes.Buffertotal points for the semester,
// which is used when creating a letter grade scale
type AssignmentList struct {
	Assignments    []Assignment
	SemesterPoints float64
}

/*
returns a new AssignmentList based on a slice of strings

ignores any lines that aren't formatted in one of these ways:
1) 400 x 2 assignment1
2) 300 x 1 assignment2
3) 140 assignment3
*/
func NewAssignmentList(input []string) *AssignmentList {

	var semesterPoints float64
	var assignments []Assignment

	// loop through the slice of strings, silently skipping invalid ones
	for _, str := range input {

		var pointString, name string

		// split the string by spaces
		fields := strings.Fields(str)
		if len(fields) < 1 {
			continue
		}

		// the first part is always the points
		points, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			continue
		}

		// check if second part is "x" indicating a multiplier
		if len(fields) > 2 && fields[1] == "x" {
			multiplier, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				continue
			}
			points = points * multiplier
			pointString = fmt.Sprintf("%s x %s", fields[0], fields[2])
			name = strings.Join(fields[3:], " ")
		} else {
			pointString = fields[0]
			name = strings.Join(fields[1:], " ")
		}

		semesterPoints += points

		assignments = append(assignments, Assignment{
			TotalPoints: points,
			PointString: pointString,
			Name:        name,
		})
	}

	return &AssignmentList{Assignments: assignments, SemesterPoints: semesterPoints}
}

// returns a slice of chart.Value based on the AssignmentList
func (al *AssignmentList) ChartVals() (vals []chart.Value) {
	for _, a := range al.Assignments {
		vals = append(vals, chart.Value{Label: a.Name, Value: a.TotalPoints})
	}
	return
}

// returns a markdown table given a slice of Assignment
func (al *AssignmentList) Markdown() string {
	buffer := &bytes.Buffer{}

	buffer.WriteString(`
| Points  | Category            |
| ------- | ------------------- |
`)
	for _, a := range al.Assignments {
		buffer.WriteString("| ")
		buffer.WriteString(a.PointString)
		buffer.WriteString(" | ")
		buffer.WriteString(a.Name)
		buffer.WriteString(" |\n")
	}
	return buffer.String()
}

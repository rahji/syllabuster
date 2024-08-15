package assignment

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2"
)

/*
an Assignment has TotalPoints (including the optional multiplier),
a PointString (which includes the optional multiplier),
and a Name that includes an optional short-name in parenthesis at its end.
The short-name is used for the pie chart label currently
*/
type Assignment struct {
	TotalPoints float64 // eg: 300
	PointString string  // eg: "100 x 3"
	Name        string  // eg: midterm exam
	ShortName   string  // eg: (midterm)
}

// includes a slice of Assignment and the total points for the semester,
// which is used when creating a letter grade scale
type AssignmentList struct {
	Assignments    []Assignment
	SemesterPoints float64
}

/*
returns a new AssignmentList based on a slice of strings

ignores any lines that aren't formatted in one of these ways:
1) 400 x 2 major projects (projects)
2) 300 x 1 midterm
3) 140 final exam (final)
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

		longname, shortname := extractLongAndShort(name)

		semesterPoints += points

		assignments = append(assignments, Assignment{
			TotalPoints: points,
			PointString: pointString,
			Name:        longname,
			ShortName:   shortname,
		})
	}

	return &AssignmentList{Assignments: assignments, SemesterPoints: semesterPoints}
}

// returns a slice of chart.Value based on the AssignmentList.
// it uses the ShortName as the chart's label
func (al *AssignmentList) ChartVals() (vals []chart.Value) {
	for _, a := range al.Assignments {
		vals = append(vals, chart.Value{Label: a.ShortName, Value: a.TotalPoints})
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

/*
finds a match for a shortname in parenthesis at the end of the string
returns the original string (without the parenthesized part) and
the string in parens as the second return value.
if there is not a shortname in the string then the shortname is set to the original
*/
func extractLongAndShort(input string) (string, string) {
	re := regexp.MustCompile(`\(([^)]+)\)$`)
	matches := re.FindStringSubmatch(input)

	// more than 1 match means that there were parens
	if len(matches) > 1 {
		insideParentheses := matches[1]
		// remove the parens from the original string
		result := re.ReplaceAllString(input, "")
		return strings.TrimSpace(result), insideParentheses
	}

	// if no match is found, return the original string twice
	return input, input
}

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rahji/syllabuster/assignment"
	"github.com/rahji/syllabuster/config"
	"github.com/rahji/syllabuster/pie"
	"github.com/rahji/syllabuster/scale"
)

func main() {

	conf, err := config.ReadConfig("mock.yaml")
	if err != nil {
		log.Fatalf("readconfig %s", err)
	}

	s := `
325 x 3 Reading Responses
475 x 6 Technical Exercises
400 x 1 Presentation/Demo
550 x 2 Project Drafts
920 x 2 Major Projects
1750 Participation
`
	lines := strings.Split(s, "\n")
	// for windows...
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	al := assignment.NewAssignmentList(lines)
	fmt.Printf("%s\n\n", al.Markdown())

	fmt.Printf("points: %.0f\n\n", al.SemesterPoints)

	str := scale.Rescale(conf.Scale, al.SemesterPoints)

	fmt.Println("\n\n" + str)

	err = pie.Draw("output.png", al.ChartVals())

}

# syllabuster

This is a command-line tool that helps me generate the grade-related page of a syllabus.

I don't like grading based on percentages. It's easier for me to come up 
with an integer number of points for each assignment/task and sum them to arrive
at a total number of available points for the semester. 

With this tool, I can enter a list of assignments with their corresponding point values 
and get two files as output:

1. A markdown file containing tables listing the assignments and the letter grade scale
2. A PNG file containing a pie chart showing the distribution of points for the whole semester

The letter grade scale is based on an existing scale, which is specified in the configuration file.
While the reference scale may use fractional numbers, all other points will be whole numbers.

This tool is a little weird, but it does exactly what I need.

## Input Format

When entering the list of assignments, be sure to look at the examples shown. 
The requirements are:

* One assignment per line
* Each line starts with the number of points for the assignment
* An optional multiplier can be added next (ie: "300 x 4" if there are four 300-point assignments)
* The name of the assignment is next
* You can add an optional short-name for the project in parenthesis. This is used for the pie chart labels.

## Installation

1. Clone the repo
2. Run `make build` from your terminal

## Usage

Run the command by entering `syllabuster` from your terminal 

## Screenshot

![image](https://github.com/user-attachments/assets/a72a0ec8-8817-452a-9e11-cb59a51b9335)

## Sample PNG Output

![smallerchart](https://github.com/user-attachments/assets/eeabef41-aefd-4aab-b580-91a6137c8f50)

## Notes

* The configuration file also contains the assignments and points but that's unused at the moment.
* The PNG file is 700 pixels x 700 pixels

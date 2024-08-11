# syllabuster

This is a command-line tool that helps me generate the grade-related page of a syllabus.

I don't like grading based on percentages. It's easier for me to come up 
with an integer number of points for each assignment/task and sum them to arrive
at a total number of available points for the semester. 

With this tool, I can enter a list of assignments with their corresponding point values 
and get a markdown page that contains the following:

1. A markdown table that lists the assignments and their corresponding point values
2. A pie chart showing the distribution of points across assignments
3. A markdown table that shows how a student's end-of-semester point total translates to a letter grade

The letter grade scale is based on an existing scale for your school, which 
can be specified in the configuration file. While the reference scale may use fractional 
numbers, all other points will be whole numbers.

The configuration file contains the percentage-to-letter-grade scale and whatever assignments/points
were entered the last time the program was used.

## Installation

The easiest way to install is to download the appropriate archive file from the [Releases](https://github.com/rahji/syllabuster/releases/latest) page, place the `syllabuster` binary [somewhere in your path](https://zwbetz.com/how-to-add-a-binary-to-your-path-on-macos-linux-windows/), and run it from your terminal (eg: Terminal.app in MacOS or [Windows Terminal](https://apps.microsoft.com/store/detail/windows-terminal/9N0DX20HK701?hl=en-us&gl=us&rtc=1))

**OR** If you have `go` installed you can clone this repo and run `make build`

You might want to install [glow](https://github.com/charmbracelet/glow), too.

## Usage

Run the command by entering `syllabuster` from your terminal. 

## Todo

* enter filename for markdown and for chart
* get button to work for generating markdown and image
* get --help and --version to work without cobra

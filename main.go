package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

// global vars
// declare flag vars
var wait, timeVar float64
var message, base string
var live, infinite bool
var ascii = [11]string{" ", "/", "\\", "~", "|", "-", "_", "@", "&", "#", "$"}

func printHelp() {
	fmt.Println("Usage: gobonsai [OPTION]...")
	fmt.Println("")
	fmt.Println("gobonsai is a beautiful randomly generated bonsai tree generator written in Golang.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("   -l, --live              live mode: show each step of growth")
	fmt.Println("   -t, --time=TIME         in live mode, wait TIME secs between")
	fmt.Println("                               steps of growth (must be larger than 0) [default: 0.03]")
	fmt.Println("   -i, --infinite          infinite mode: keep growing trees")
	fmt.Println("   -w, --wait=TIME         in infinite mode, wait TIME between each tree")
	fmt.Println("                           generation [default: 4.00]")
	fmt.Println("   -s, --screensaver       screensaver mode: equivilant of -li and")
	fmt.Println("                               q to quit")
	fmt.Println("   -m, --message=STR       attach message next to the tree")
	fmt.Println("   -b, --base=INT          ascii-art plant base to use (1 or 2), 0 is none")
	os.Exit(0)
}

func clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func genTree(wh ...int) {
	// generate grid
	w := wh[0]
	h := wh[1]

	var widthSlice []int
	for i := w; i > 0; i-- {
		widthSlice = append(widthSlice, 1)
	}

	var grid [][]int
	for i := h; i > 0; i-- {
		grid = append(grid, widthSlice)
	}

	// loop infinitely
	// int values:
	// 0 empty
	// 1 early wood
	// 2 mid wood
	// 3 late wood
	// 4 dead wood
	// 5 leaf
	// 6 outer leaf
	for {
		clear()
		if live {
			showTree(grid)
			fmt.Println("GOING TO SLEEP")
			time.Sleep(time.Duration(timeVar * float64(time.Second)))
			fmt.Println("SHOULD HAVE SLEPT")
		}

		treeDead := true
		// tree growth logic

		// if not infinite and no more growing to do, break the for loop
		if !infinite && treeDead {
			break
		}
	}
	showTree(grid)
}

func showTree(grid [][]int) {
	// ascii values
	// 0 empty
	// 1-6 wood
	// 7-10 leaf
	for r := range grid {
		for c := range grid[r] {
			rand.Seed(time.Now().UnixNano())
			var num int
			if grid[r][c] == 0 {
				num = 0
			} else if grid[r][c] < 7 {
				min, max := 1, 6
				num = rand.Intn(max-min) + min
			} else {
				min, max := 7, 10
				num = rand.Intn(max-min) + min
			}

			fmt.Print(ascii[num])
		}
		// new line for new row
		fmt.Print("\n")
	}
}

func main() {
	// map decleration because consts arent allowed for maps
	var bases = map[int]string{
		0: "",
		1: " \\                           / \n  \\_________________________/ \n  (_)                     (_)",
		2: " (           ) \n  (_________)  ",
	}

	// declare flag vars that should be global
	var baseInt int
	var screensaver bool

	// set help command
	flag.Usage = printHelp

	// set flags + default values
	flag.Float64Var(&timeVar, "time", 0.03, "wait TIME secs between steps of growth (must be larger than 0) [default: 0.03]")
	flag.Float64Var(&wait, "wait", 4.00, "in infinite mode, wait TIME between each tree generation [default: 4.00]")
	flag.StringVar(&message, "message", "", "attach message next to the tree")
	flag.IntVar(&baseInt, "base", 1, "ascii-art plant base to use (1 or 2), 0 is none")
	flag.BoolVar(&live, "live", false, "live mode: show each step of growth")
	flag.BoolVar(&infinite, "infinite", false, "infinite mode: keep growing trees")
	flag.BoolVar(&screensaver, "screensaver", false, "screensaver mode: equivilant of -li and q to quit")
	// shortened version of var names
	flag.Float64Var(&timeVar, "t", 0.03, "wait TIME secs between steps of growth (must be larger than 0) [default: 0.03]")
	flag.Float64Var(&wait, "w", 4.00, "in infinite mode, wait TIME between each tree generation [default: 4.00]")
	flag.StringVar(&message, "m", "", "attach message next to the tree")
	flag.IntVar(&baseInt, "b", 1, "ascii-art plant base to use (1 or 2), 0 is none")
	flag.BoolVar(&live, "l", false, "live mode: show each step of growth")
	flag.BoolVar(&infinite, "i", false, "infinite mode: keep growing trees")
	flag.BoolVar(&screensaver, "s", false, "screensaver mode: equivilant of -li and q to quit")

	// parse the flags
	flag.Parse()

	// if screensaver, override -li to true
	if screensaver {
		live = true
		infinite = true
	}

	// set base
	if !(0 <= baseInt && baseInt <= 2) {
		fmt.Println("Base must be between 0 and 2.")
		fmt.Println("Quitting...")
		return
	}
	base = bases[baseInt]

	// get terminal size
	if !term.IsTerminal(0) {
		return
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}

	// DEBUG use preset width + height for debugging
	genTree(width, height)

	for infinite {
		// regen width + height incase terminal was resized
		width, height, err = term.GetSize(0)
		if err != nil {
			return
		}
		genTree(width, height)
	}
}

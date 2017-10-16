package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

/*================================= types =========================*/
type SelpgArgs struct {
	StartPage  int
	EndPage    int
	PageLen    int
	Pagetype   int //false:l. true:f
	inFilename string
	printDest  string
}

const Inbufsiz = 16 * 1024
const IntMax = 1000000000000

/*=================================  ls =======================*/

var progname string /* program name, for error messages */

/*================================= main()=== =====================*/
func main() {
	progname = "selpg"
	selpgargs := SelpgArgs{StartPage: -1, EndPage: -1, PageLen: 72, Pagetype: 'l', inFilename: ""}
	processArgs(&selpgargs)
	processInput(selpgargs)
}

/*================================= process_args() ================*/
func testStat(filePath string, rw int) int {
	res := [4]int{0, 0, 0, 0}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		//fmt.Printf(filePath)
		return -1
	}
	res[0] = 1
	//fmt.Println(fileInfo.Mode())
	mode := int(fileInfo.Mode())
	i := 3
	for i > 0 {
		res[i] = mode % 8
		mode = mode / 8
		i = i - 1
	}
	for i = 0; i < 4; i++ {
		//fmt.Println(res[i])
	}
	k := 3
	right := [2]int{0, 0}
	right[0] = res[k] / 4                //read
	right[1] = (res[k] - right[0]*4) / 2 //write
	//fmt.Println(right)
	if rw == 1 {
		return right[0]
	} else {
		return right[1]
	}
}

func processArgs(selpgargs *SelpgArgs) {
	startPage := flag.Int("s", 0, "the start page")
	endPage := flag.Int("e", 0, "the start page")
	pageTypeF := flag.Bool("f", false, "sepearte pages use EOF")
	pageLen := flag.Int("l", 72, "the length of a page")
	printdest := flag.String("d", "", "the dest sub program")
	flag.Parse()

	if *startPage <= 0 || *endPage <= 0 || *startPage > *endPage {
		fmt.Println("illegal start or end page number")
		os.Exit(1)
	}
	selpgargs.StartPage = *startPage
	selpgargs.EndPage = *endPage
	selpgargs.PageLen = *pageLen
	selpgargs.printDest = *printdest
	if *pageTypeF {
		selpgargs.Pagetype = 'f'
	}
	//fmt.Println(flag.Args())
	if len(flag.Args()) > 0 {
		selpgargs.inFilename = (flag.Args())[0]
		if testStat(selpgargs.inFilename, 1) != 1 {
			fmt.Println("you can't  read the file")
			os.Exit(2)
		}
	}
	//fmt.Println(*selpgargs)
}

func processInput(selpgargs SelpgArgs) {
	var reader io.Reader
	var writer io.Writer
	var err error
	var EOF rune = 26

	if selpgargs.inFilename == "" {
		reader = os.Stdin
	} else {
		reader, err = os.Open(selpgargs.inFilename)
		if err != nil {
			fmt.Println(progname, ": could not open input file ", selpgargs.inFilename)
			os.Exit(16)
		}
	}

	/* set the output destination */
	if selpgargs.printDest == "" {
		writer = os.Stdout
	} else {
		//fmt.Println("pipe")
		//subproc := exec.Command("cat", "-n", selpgargs.printDest)
		subproc := exec.Command("cat", "-n")
		writer, err = subproc.StdinPipe()
		if err != nil {
			fmt.Println("could not open pipe to ", progname)
			os.Exit(17)
		}
		subproc.Stdout = os.Stdout
		subproc.Stderr = os.Stderr
		subproc.Start()
	}

	Reader := bufio.NewReader(reader)
	var crc string
	var lineCtr int
	var pageCtr int
	/* begin one of two main loops based on page type */
	if selpgargs.Pagetype == 'l' {
		lineCtr = 0
		pageCtr = 1

		for true {
			crc, _ = Reader.ReadString('\n')
			if crc == "" {
				break
			}
			lineCtr++
			if lineCtr > selpgargs.PageLen {
				pageCtr++
				lineCtr = 1
			}
			if (pageCtr >= selpgargs.StartPage) && (pageCtr <= selpgargs.EndPage) {
				io.WriteString(writer, crc)
			}
		}
	} else {
		//fmt.Println("f")
		Writer := bufio.NewWriter(writer)
		pageCtr = 1
		for true {
			r, _, _ := Reader.ReadRune()
			if r == EOF {
				break
			}
			if r == '\f' {
				pageCtr++
			}
			if (pageCtr >= selpgargs.StartPage) && (pageCtr <= selpgargs.EndPage) {
				Writer.WriteRune(r)
			}
		}
	}

	/* end main loop */
	if pageCtr < selpgargs.StartPage {
		fmt.Println(progname, ": start_page greater than total pages, no output written")
	} else if pageCtr < selpgargs.EndPage {
		fmt.Println(progname, ": end_page greater than total pages, less output than expected")
	}
}

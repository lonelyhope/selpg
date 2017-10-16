package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	args := os.Args
	progname = args[0]
	//fmt.Println(args)
	selpgargs := SelpgArgs{StartPage: -1, EndPage: -1, PageLen: 72, Pagetype: 'l', inFilename: ""}
	processArgs(args, &selpgargs)
	processInput(selpgargs)
}

/*================================= process_args() ================*/
func truncateToInt(s string, start, end int) int {
	st := s[start:end]
	t, _ := strconv.Atoi(st)
	return t
}

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

func processArgs(args []string, selpgargs *SelpgArgs) {
	//fmt.Println("processArgs")
	var s1, s2 string
	var argno int
	var i int

	//fmt.Println(args)
	argLen := len(args)
	if argLen < 3 {
		fmt.Println(progname, ": not enough arguments")
		os.Exit(1)
	}

	/* handle 1st arg - start page */
	s1 = args[1]
	if len(s1) < 3 || strings.Compare(s1[0:2], "-s") != 0 {
		fmt.Println(progname, ": 1st arg should be start_page")
		os.Exit(2)
	}
	i = truncateToInt(s1, 2, len(s1))
	if i < 1 || i > (IntMax-1) {
		fmt.Println(progname, ": invalid start page")
		os.Exit(3)
	}
	selpgargs.StartPage = i

	/* handle 2nd arg - end page */
	s1 = args[2]
	if len(s1) < 3 || strings.Compare(s1[0:2], "-e") != 0 {
		//fmt.Println(s1)
		fmt.Println(progname, ": 2nd arg should be -eend_page")
		os.Exit(4)
	}
	i = truncateToInt(s1, 2, len(s1))
	if (i < 1) || (i > (IntMax-1) || i < selpgargs.StartPage) {
		fmt.Println(progname, ": invalid end page")
		os.Exit(5)
	}
	selpgargs.EndPage = i

	/* now handle optional args */
	argno = 3
	for argno <= (argLen-1) && args[argno][0] == '-' {
		/* while there more args and they start with a '-' */
		s1 = args[argno]
		switch s1[1] {
		case 'l':
			s2 = s1[2:len(s1)]
			i, _ = strconv.Atoi(s2)
			if i < 1 || i > (IntMax-1) {
				fmt.Println(progname, ": invalid page length ", s2)
				os.Exit(6)
			}
			selpgargs.PageLen = i
			selpgargs.Pagetype = 'l'
			argno = argno + 1

		case 'f':
			/* check if just "-f" or something more */
			if strings.Compare(s1, "-f") != 0 {
				fmt.Println(progname, ": option should be \"-f\"")
				os.Exit(7)
			}
			selpgargs.Pagetype = 'f'
			argno = argno + 1

		case 'd':
			s2 = s1[2:len(s1)]
			if len(s2) < 1 {
				fmt.Println(progname, ": -d option requires a printer destination")
				os.Exit(8)
			}
			selpgargs.printDest = s2
			argno = argno + 1

		default:
			fmt.Println(progname, ": unknown option ", s1)
			os.Exit(9)
		} /* end switch */
	} /* end for */

	/*++argno;*/
	if argno <= argLen-1 { /* there is one more arg */
		selpgargs.inFilename = args[argno]
		//fileInfo, err := os.Stat(selpgargs.inFilename)
		/* check if file exists */
		if testStat(selpgargs.inFilename, 1) == -1 {
			fmt.Println(progname, ": input file  does not exist")
			os.Exit(10)
		}
		/* check if file is readable */
		if testStat(selpgargs.inFilename, 1) == 0 {
			fmt.Println(progname, ": input file exists but cannot be read")
			os.Exit(11)
		}
	}
	//fmt.Println(selpgargs)
	if selpgargs.StartPage < 0 {
		fmt.Println("The Start Page should bigger to 0")
		os.Exit(12)
	}
	if selpgargs.EndPage < 0 || selpgargs.EndPage < selpgargs.StartPage {
		fmt.Println("the end page should be bigger to start page")
		os.Exit(13)
	}
	if selpgargs.PageLen < 1 {
		fmt.Println("the page length should be bigger than 0")
		os.Exit(14)
	}
	if selpgargs.Pagetype != 'l' && selpgargs.Pagetype != 'f' {
		fmt.Println(selpgargs.Pagetype)
		fmt.Println("illegal page type")
		os.Exit(15)
	}
}

func processInput(selpgargs SelpgArgs) {
	//fmt.Println("processInput")
	//fmt.Println(selpgargs)
	var reader io.Reader
	var writer io.Writer
	var err error
	var EOF rune = 26
	//Writer := bufio.NewWriter(writer)
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

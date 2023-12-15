package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func problemParser(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i, line := range lines {
		r[i] = problem{q: line[0], a: line[1]}
	}
	return r
}

func problemPuller(filename string) ([]problem, error) {
	if fObj, err := os.Open(filename); err == nil {
		csvR := csv.NewReader(fObj)
		if cLines, err := csvR.ReadAll(); err == nil {
			return problemParser(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data from csv format from  %s file; %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", filename, err.Error())
	}
}

func main() {
	fName := flag.String("f", "quiz.csv", "The name of the file to be parsed")
	timer := flag.Int("t", 30, "The time limit for the quiz in seconds")
	flag.Parse()
	problems, err := problemPuller(*fName)
	if err != nil {
		exit(fmt.Sprintf("something went wrong %s", err.Error()))
	}
	correctAns := 0
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)
programLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		go func() {
			fmt.Scanf("%s\n", &answer)
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println("")
			break programLoop
		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}

	}
	fmt.Printf("Your result is %d out of %d", correctAns, len(problems))
	fmt.Println("Press enter two exit")
	<-ansC
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Node struct {
	name     string
	children []Node
}

func NewNode() *Node {
	return &Node{}
}

type Element struct {
	key   string
	value string
}

type Parser struct {
	root        Node
	openBraces  int
	closeBraces int
	lines       int
}

func NewParser() *Parser {
	return &Parser{
		root:        *NewNode(),
		openBraces:  0,
		closeBraces: 0,
		lines:       0,
	}
}

func Normalize(t string) string {
	return strings.Trim(strings.TrimSpace(t), `\t`)
}

func (p *Parser) Parse(text string) error {
	text = Normalize(text)
	r, _ := regexp.Compile(`\"([A-Za-z0-9\\\:\-\(\)\ \_]*)\"`)
	if text == "" {
		return nil
	}
	p.lines++
	if text == "{" {
		p.openBraces++
	} else if text == "}" {
		p.closeBraces++
	} else if r.MatchString(text) {
		fmt.Println(r.FindAllString(text, -1))
	} else {
		return errors.New("Parser Error: Contains the wrong string.")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Println("need args as input file")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	parser := NewParser()
	for scanner.Scan() {
		line := scanner.Text()
		err = parser.Parse(line)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if parser.openBraces != parser.closeBraces {
		log.Println("Parser Error: Unmatched number of parentheses.")
		return
	}
	if err = scanner.Err(); err != nil {
		log.Println(err)
		return
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Node struct {
	Name     string
	Parent   *Node
	Children []*Node
	List     map[string]string
}

func NewNode(name string) *Node {
	return &Node{Name: name, List: make(map[string]string)}
}

func (n *Node) AddChild(name string) *Node {
	childNode := NewNode(name)
	childNode.Parent = n
	n.Children = append(n.Children, childNode)
	return childNode
}

func (n *Node) ReturnParent() *Node {
	if n.Parent != nil {
		return n.Parent
	}
	return nil
}

func (n *Node) SearchTree(nodeLevel int) {
	if nodeLevel == 0 {
		fmt.Printf("Level\tTag\tElements\n")
		fmt.Printf("----------------------------------------\n")
	}
	fmt.Printf("%d\t%s\n", nodeLevel, n.Name)
	if n.List != nil {
		for key, value := range n.List {
			fmt.Printf("\t\t%-10s: %s\n", key, value)
		}
	}
	if n.Children != nil {
		for _, child := range n.Children {
			child.SearchTree(nodeLevel + 1)
		}
	}
}

type Parser struct {
	Root   *Node
	Cursor *Node
}

func NewParser() *Parser {
	return &Parser{}
}

func Normalize(t string) string {
	return strings.Trim(strings.TrimSpace(t), `\t`)
}

func (p *Parser) Parse(lines []string) error {
	// Normalize
	for i, _ := range lines {
		lines[i] = Normalize(lines[i])
	}

	// Parse Lines
	index := 0
	nodeLevel := 0
	r, _ := regexp.Compile(`\"([A-Za-z0-9\\\:\-\(\)\ \_\.]*)\"`)

loop:
	for {
		if index == len(lines) {
			if nodeLevel == 0 {
				break loop
			}
			return errors.New("Parser Error: Unmatched number of braces.")
		}

		line := lines[index]
		index++

		if line == "" {
			continue loop
		} else if r.MatchString(line) {
			matches := r.FindAllString(line, -1)
			if len(matches) == 2 {
				key := strings.Trim(matches[0], `"`)
				value := strings.Trim(matches[1], `"`)
				p.Cursor.List[key] = value
				continue loop
			} else if len(matches) == 1 && lines[index] == "{" {
				index++
				tag := strings.Trim(matches[0], `"`)
				if nodeLevel == 0 {
					p.Root = NewNode(tag)
					p.Cursor = p.Root
					nodeLevel++
				} else {
					p.Cursor = p.Cursor.AddChild(tag)
					nodeLevel++
				}
				continue loop
			} else {
				return errors.New("Parser Error: Contains the wrong string.")
			}
		} else if line == "}" {
			if p.Cursor.Parent != nil {
				p.Cursor = p.Cursor.ReturnParent()
			}
			nodeLevel--
			continue loop
		} else {
			return errors.New("Parser Error: Contains the wrong string.")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Println("need args as input file")
		return
	}
	path := os.Args[1]
	ext := filepath.Ext(path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	if ext != ".acf" && ext != ".vdf" {
		log.Println("Extention Error: This file is not '.acf' or 'vdf' file.")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	parser := NewParser()
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Println(err)
		return
	}
	err = parser.Parse(lines)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Parse Success!")

	// Test
	parser.Root.SearchTree(0)
}

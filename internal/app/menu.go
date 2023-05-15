package app

import (
	. "battleships/internal/utils"
	"fmt"
	"github.com/buger/goterm"
	"github.com/pkg/term"
	"log"
)

const (
	up     byte = 65
	down   byte = 66
	escape byte = 27
	enter  byte = 13
)

var keys = map[byte]bool{
	up:   true,
	down: true,
}

type Menu struct {
	prompt    string
	cursor    int
	menuItems []*struct {
		Text string
		Id   string
	}
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		prompt:    prompt,
		menuItems: make([]*struct{ Text, Id string }, 0),
	}
}

func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &struct {
		Text string
		Id   string
	}{option, id}

	m.menuItems = append(m.menuItems, menuItem)
	return m
}

func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		fmt.Printf("\033[%dA", len(m.menuItems)-1)
	}

	for index, menuItem := range m.menuItems {
		var newline = If(index == len(m.menuItems)-1, "", "\n")

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.cursor {
			cursor = goterm.Color("> ", goterm.YELLOW)
			menuItemText = goterm.Color(menuItemText, goterm.YELLOW)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

func (m *Menu) Display() string {
	defer func() {
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("%s\n", goterm.Color(goterm.Bold(m.prompt)+":", goterm.CYAN))
	m.renderMenuItems(false)

	fmt.Printf("\033[?25l")

	for {
		keyCode := getInput()
		switch keyCode {
		case up:
			m.cursor = (m.cursor + len(m.menuItems) - 1) % len(m.menuItems)
			m.renderMenuItems(true)
		case down:
			m.cursor = (m.cursor + 1) % len(m.menuItems)
			m.renderMenuItems(true)
		case enter:
			menuItem := m.menuItems[m.cursor]
			fmt.Println("\r")
			return menuItem.Id
		case escape:
			return ""
		}
	}
}

func getInput() byte {
	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)

	t.Restore()
	t.Close()

	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

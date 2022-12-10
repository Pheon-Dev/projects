package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle      = lipgloss.NewStyle().Padding(1, 2)
	quitTextStyle = lipgloss.NewStyle().Padding(1, 2)
)

type item struct {
	title       string
	path        string
	description string
}

type model struct {
	list     list.Model
	choice   string
	path     string
	quitting bool
}

type editorFinishedMsg struct{ err error }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func openEditor(path string) tea.Cmd {
	home := os.Getenv("HOME")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command("bash", "-c", "clear && cd "+home+"/"+path+" && "+editor)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			m.quitting = true
			return m, tea.Quit
		case "l":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				m.path = string(i.path)
			}
			return m, openEditor(m.path)
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				m.path = string(i.path)
			}
			return m, tea.Quit
		}

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

func main() {
	items := []list.Item{
		item{title: "nvim", description: "NEOVIM Config", path: ".config/nvim"},
		item{title: "dwm", description: "DWM Config", path: "cd $HOME/.config/dwm && nvim"},
		item{title: "zsh", description: "ZSH Config", path: "cd $HOME/.config/zsh && nvim"},
		item{title: "tmux", description: "TMUX Config", path: "cd $HOME/.tmux && nvim"},
		item{
			title:       "st",
			description: "Simple Terminal (ST) Config",
			path:        "cd $HOME/.config/st && nvim",
		},
		item{
			title:       "lazygit",
			description: "Lazygit Config",
			path:        "cd $HOME/.config/lazygit && nvim",
		},
		item{
			title:       "ranger",
			description: "Ranger Config",
			path:        "cd $HOME/.config/ranger && nvim",
		},
		item{
			title:       "fm",
			description: "File Manager (FM) Config",
			path:        "cd $HOME/.config/fm && nvim",
		},
		item{title: "moc", description: "MOCP Config", path: "cd $HOME/.moc && nvim && nvim"},
		item{
			title:       "go",
			description: "GO Projects",
			path:        "cd $HOME/Documents/go/src/github.com/Pheon-Dev && nvim",
		},
		item{
			title:       "go-git",
			description: "GO Git Projects",
			path:        "cd $HOME/Documents/go/git && nvim",
		},
		item{
			title:       "typescript",
			description: "TypeScript Projects",
			path:        "cd $HOME/Documents/NextJS/App && nvim",
		},
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error Running Program : ", err)
		os.Exit(1)
	}
}

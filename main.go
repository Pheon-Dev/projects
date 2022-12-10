package main

import (
	"fmt"
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

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
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
		case "h", "q", "escape":
			m.quitting = true
			return m, tea.Quit
		case "l", "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				m.path = string(i.path)
			}
			return m, openEditor(m.path)
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
		item{title: "dwm", description: "DWM Config", path: ".config/arco-dwm"},
		item{title: "zsh", description: "ZSH Config", path: ".config/zsh"},
		item{title: "tmux", description: "TMUX Config", path: ".tmux"},
		item{
			title:       "st",
			description: "Simple Terminal (ST) Config",
			path:        ".config/arco-st",
		},
		item{
			title:       "lazygit",
			description: "Lazygit Config",
			path:        ".config/lazygit",
		},
		item{
			title:       "ranger",
			description: "Ranger Config",
			path:        ".config/ranger",
		},
		item{
			title:       "fm",
			description: "File Manager (FM) Config",
			path:        ".config/fm",
		},
		item{title: "moc", description: "MOCP Config", path: ".moc"},
		item{
			title:       "projects",
			description: "p App",
			path:        "Documents/go/src/github.com/Pheon-Dev/p",
		},
		item{
			title:       "go",
			description: "GO Projects",
			path:        "Documents/go/src/github.com/Pheon-Dev",
		},
		item{
			title:       "go-git",
			description: "GO Git Projects",
			path:        "Documents/go/git",
		},
		item{
			title:       "typescript",
			description: "TypeScript Projects",
			path:        "Documents/NextJS/App",
		},
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error Running Program : ", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

var (
	docStyle      = lipgloss.NewStyle().Padding(1, 2)
	quitTextStyle = lipgloss.NewStyle().Padding(1, 2)
	titleStyle    = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5")).
			Background(lipgloss.Color("#536c9e")).
			Padding(0, 1)
	itemStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type item struct {
	title       string
	description string
}

type model struct {
	list        list.Model
	choice      string
	description string
}

type editorFinishedMsg struct{ err error }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func openEditor(description string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command("bash", "-c", "clear && cd "+description+" && "+editor)
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
		case "q", "escape":
			return m, tea.Quit
		case " ", "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				m.description = string(i.description)
			}
			return m, openEditor(m.description)
		}

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func main() {
	vp := viper.New()

	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("$HOME/.config/p")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	title := vp.GetString("title")
	statusbar := vp.GetBool("status-bar")
	filtering := vp.GetBool("filtering")
	// prjcts := vp.Get("projects")

	projects := []list.Item{
		item{title: "nvim", description: "$HOME/.config/nvim"},
		item{title: "dwm", description: "$HOME/.config/arco-dwm"},
		item{title: "zsh", description: "$HOME/.config/zsh"},
		item{title: "dmenu", description: "$HOME/.config/dmenu"},
		item{title: "btop", description: "$HOME/.config/btop"},
		item{title: "tmux", description: "$HOME/.tmux"},
		item{
			title:       "st Simple Terminal",
			description: "$HOME/.config/arco-st",
		},
		item{
			title:       "lazygit",
			description: "$HOME/.config/lazygit",
		},
		item{
			title:       "ranger",
			description: "$HOME/.config/ranger",
		},
		item{
			title:       "fm file manager",
			description: "$HOME/.config/fm",
		},
		item{title: "moc", description: ".moc"},
		item{
			title:       "p app",
			description: "$HOME/Documents/go/src/github.com/Pheon-Dev/p",
		},
		item{
			title:       "neovim",
			description: "$HOME/Documents/Neovim",
		},
		item{
			title:       "class",
			description: "$HOME/Documents/CMT",
		},
		item{
			title:       "go",
			description: "$HOME/Documents/go/src/github.com/Pheon-Dev",
		},
		item{
			title:       "bubbletea",
			description: "$HOME/Documents/go/git/bubbletea/examples",
		},
		item{
			title:       "go apps",
			description: "$HOME/Documents/go/git",
		},
		item{
			title:       "destiny",
			description: "$HOME/Documents/NextJS/App/destiny-credit",
		},
		item{
			title:       "devlen",
			description: "$HOME/Documents/NextJS/App/devlen",
		},
		item{
			title:       "typescript",
			description: "$HOME/Documents/NextJS/App",
		},
	}

	// vp.Set("title", "Configs")
	// vp.Set("status-bar", true)
	// vp.Set("filtering", true)
	// vp.Set("projects", projects)
	// vp.WriteConfig()

	l := list.New(projects, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(statusbar)
	l.SetFilteringEnabled(filtering)
	l.Styles.Title = titleStyle
	m := model{list: l}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error Running Program : ", err)
		os.Exit(1)
	}
}

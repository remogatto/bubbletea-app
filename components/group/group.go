package group

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	foam "github.com/remogatto/sugarfoam"
	"github.com/remogatto/sugarfoam/layout"
)

type Option func(*Model)

type Model struct {
	foam.Common

	KeyMap KeyMap

	layout *layout.Layout
	items  []foam.Focusable

	focused   bool
	currFocus int
}

func New(opts ...Option) *Model {
	group := new(Model)

	group.KeyMap = DefaultKeyMap()
	group.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(group)
	}

	return group
}

type KeyMap struct {
	FocusNext key.Binding
	FocusPrev key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		FocusNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "focus next"),
		),
		FocusPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "focus prev")),
	}
}

func (g *Model) SetSize(width int, height int) {
	g.Common.SetSize(width, height)
	g.layout.SetSize(g.GetWidth(), g.GetHeight())
}

func (g *Model) Current() foam.Focusable {
	return g.items[g.currFocus]
}

func (g *Model) Blur() {
	g.focused = false
}

func (g *Model) Focus() tea.Cmd {
	g.focused = true
	return g.Current().Focus()
}

func (g *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range g.items {
		cmds = append(cmds, c.Init())
	}

	g.Current().Focus()

	return tea.Batch(cmds...)
}

func (g *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, g.KeyMap.FocusNext):
			cmds = append(cmds, g.nextFocus())
		case key.Matches(msg, g.KeyMap.FocusPrev):
			cmds = append(cmds, g.prevFocus())
		}
	}

	cmds = append(cmds, g.updateComponents(msg)...)

	return g, tea.Batch(cmds...)
}

func (g *Model) View() string {
	return g.layout.View()
}

func (m *Model) CanGrow() bool {
	return true
}

func WithKeyMap(km KeyMap) Option {
	return func(g *Model) {
		g.KeyMap = km
	}
}

func WithItems(items ...foam.Focusable) Option {
	return func(g *Model) {
		g.items = items
	}
}

func WithLayout(layout *layout.Layout) Option {
	return func(m *Model) {
		m.layout = layout
	}
}

func (g *Model) nextFocus() tea.Cmd {
	g.Current().Blur()

	g.currFocus--
	if g.currFocus < 0 {
		g.currFocus = len(g.items) - 1
	}

	return g.Current().Focus()
}

func (g *Model) prevFocus() tea.Cmd {
	g.Current().Blur()

	g.currFocus = (g.currFocus - 1) % len(g.items)

	return g.Current().Focus()
}

func (g *Model) updateComponents(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, item := range g.items {
		_, cmd := item.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

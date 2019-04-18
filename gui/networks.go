package gui

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/skanehira/docui/common"
	"github.com/skanehira/docui/docker"
)

type network struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	containers string
}

type networks struct {
	*tview.Table
}

func newNetworks(g *Gui) *networks {
	networks := &networks{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	networks.SetTitle("network list").SetTitleAlign(tview.AlignLeft)
	networks.SetBorder(true)
	networks.setEntries(g)
	networks.setKeybinding(g)
	return networks
}

func (n *networks) name() string {
	return "networks"
}

func (n *networks) setKeybinding(g *Gui) {
	n.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectNetwork()
		}

		switch event.Rune() {
		case 'd':
			g.removeNetwork()
		}

		return event
	})
}

func (n *networks) entries(g *Gui) {
	networks, err := docker.Client.Networks(types.NetworkListOptions{})
	if err != nil {
		common.Logger.Error(err)
		return
	}

	keys := make([]string, 0, len(networks))
	tmpMap := make(map[string]*network)

	for _, net := range networks {
		var containers string

		net, err := docker.Client.InspectNetwork(net.ID)
		if err != nil {
			common.Logger.Error(err)
			continue
		}

		for _, endpoint := range net.Containers {
			containers += fmt.Sprintf("%s ", endpoint.Name)
		}

		tmpMap[net.ID[:12]] = &network{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			containers: containers,
		}

		keys = append(keys, net.ID[:12])

	}

	g.state.resources.networks = make([]*network, 0)
	for _, key := range common.SortKeys(keys) {
		g.state.resources.networks = append(g.state.resources.networks, tmpMap[key])
	}
}

func (n *networks) setEntries(g *Gui) {
	n.entries(g)
	table := n.Clear()

	headers := []string{
		"ID",
		"Name",
		"Driver",
		"Scope",
		"Containers",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, network := range g.state.resources.networks {
		table.SetCell(i+1, 0, tview.NewTableCell(network.ID).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(network.Name).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(network.Driver).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 3, tview.NewTableCell(network.Scope).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 4, tview.NewTableCell(network.containers).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (n *networks) focus(g *Gui) {
	n.SetSelectable(true, false)
	g.app.SetFocus(n)
}

func (n *networks) unfocus() {
	n.SetSelectable(false, false)
}

func (n *networks) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		n.setEntries(g)
	})
}
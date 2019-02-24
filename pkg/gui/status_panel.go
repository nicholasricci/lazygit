package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) refreshStatus(g *gocui.Gui) error {
	v, err := g.View("status")
	if err != nil {
		panic(err)
	}
	// for some reason if this isn't wrapped in an update the clear seems to
	// be applied after the other things or something like that; the panel's
	// contents end up cleared
	pushables, pullables := gui.GitCommand.GetCurrentBranchUpstreamDifferenceCount()
	if err := gui.updateWorkTreeState(); err != nil {
		return err
	}
	g.Update(func(*gocui.Gui) error {
		v.Clear()
		fmt.Fprint(v, "↑"+pushables+"↓"+pullables)
		branches := gui.State.Branches
		if gui.State.WorkingTreeState != "normal" {
			fmt.Fprint(v, utils.ColoredString(fmt.Sprintf(" (%s)", gui.State.WorkingTreeState), color.FgYellow))
		}

		if len(branches) == 0 {
			return nil
		}
		branch := branches[0]
		name := utils.ColoredString(branch.Name, branch.GetColor())
		repo := utils.GetCurrentRepoName()
		fmt.Fprint(v, " "+repo+" → "+name)
		return nil
	})

	return nil
}

func (gui *Gui) handleCheckForUpdate(g *gocui.Gui, v *gocui.View) error {
	gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)
	return gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("CheckingForUpdates"))
}

func (gui *Gui) handleStatusSelect(g *gocui.Gui, v *gocui.View) error {
	magenta := color.New(color.FgMagenta)

	dashboardString := strings.Join(
		[]string{
			lazygitTitle(),
			"Copyright (c) 2018 Jesse Duffield",
			"Keybindings: https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings.md",
			"Config Options: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md",
			"Tutorial: https://youtu.be/VDXvbHZYeKY",
			"Raise an Issue: https://github.com/jesseduffield/lazygit/issues",
			magenta.Sprint("Buy Jesse a coffee: https://donorbox.org/lazygit"), // caffeine ain't free
		}, "\n\n")

	return gui.renderString(g, "main", dashboardString)
}

func (gui *Gui) handleOpenConfig(g *gocui.Gui, v *gocui.View) error {
	return gui.openFile(gui.Config.GetUserConfig().ConfigFileUsed())
}

func (gui *Gui) handleEditConfig(g *gocui.Gui, v *gocui.View) error {
	filename := gui.Config.GetUserConfig().ConfigFileUsed()
	return gui.editFile(filename)
}

func lazygitTitle() string {
	return `
   _                       _ _
  | |                     (_) |
  | | __ _ _____   _  __ _ _| |_
  | |/ _` + "`" + ` |_  / | | |/ _` + "`" + ` | | __|
  | | (_| |/ /| |_| | (_| | | |_
  |_|\__,_/___|\__, |\__, |_|\__|
                __/ | __/ |
               |___/ |___/       `
}

func (gui *Gui) updateWorkTreeState() error {
	merging, err := gui.GitCommand.IsInMergeState()
	if err != nil {
		return err
	}
	if merging {
		gui.State.WorkingTreeState = "merging"
		return nil
	}
	isRebasing, err := gui.GitCommand.IsInRebasingState()
	if err != nil {
		return err
	}
	if isRebasing {
		gui.State.WorkingTreeState = "rebasing"
	} else {
		gui.State.WorkingTreeState = "normal"
	}
	return nil
}

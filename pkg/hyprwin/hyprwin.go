package hyprwin

import (
	"errors"
	"slices"
)

type manager struct {
	command *command
	ipc     IPC
}

func Dispatch(command *command) (out string, err error) {
	mgr := manager{command, InitIPC()}

	win, err := mgr.ipc.ActiveWindow()
	if err != nil {
		return "", err
	}

	var resp []byte
	switch command.dispatcher {
	case dispatcher("movewindow"):
		resp, err = mgr.moveWindow(win)
	case dispatcher("movefocus"):
		resp, err = mgr.moveFocus(win)
	default:
		err = errors.New("unknown dispatcher")
	}
	return string(resp), err
}

func (m manager) moveWindow(win *winObj) (resp []byte, err error) {
	dir := m.command.direction.Str()

	if m.command.direction.ToMonitor() {
		return m.ipc.Hyprctl("dispatch moveoutofgroup", "dispatch movewindow "+dir)
	}

	pos := slices.Index(win.Grouped, win.Address)
	grpSize := len(win.Grouped)
	edgeL := pos == 0 && dir == "l"
	edgeR := pos == grpSize-1 && dir == "r"
	v := dir == "u" || dir == "d"

	if grpSize == 1 {
		return m.ipc.Hyprctl("dispatch movewindow " + dir)
	}
	if grpSize == 0 || edgeL || edgeR || v {
		return m.ipc.Hyprctl("dispatch movewindoworgroup " + dir)
	}

	dir = map[string]string{"l": "b", "r": "f"}[dir]
	return m.ipc.Hyprctl("dispatch movegroupwindow " + dir)
}

func (m manager) moveFocus(win *winObj) (resp []byte, err error) {
	dir := m.command.direction.Str()

	pos := slices.Index(win.Grouped, win.Address)
	grpSize := len(win.Grouped)
	edgeL := pos == 0 && dir == "l"
	edgeR := pos == grpSize-1 && dir == "r"
	v := dir == "u" || dir == "d"

	if grpSize <= 1 || edgeL || edgeR || v {
		return m.ipc.Hyprctl("dispatch movefocus " + dir)
	}
	dir = map[string]string{"l": "b", "r": "f"}[dir]
	return m.ipc.Hyprctl("dispatch changegroupactive " + dir)
}

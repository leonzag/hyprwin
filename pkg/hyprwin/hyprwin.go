package hyprwin

import (
	"errors"
	"slices"
)

type manager struct {
	cmd *CommandRequest
	ipc IPC
}

func Dispatch(cmd *CommandRequest) (out string, err error) {
	ipc, err := InitIPC()
	if err != nil {
		return "", err
	}
	mgr := manager{cmd, ipc}

	win, err := mgr.ipc.ActiveWindow()
	if err != nil {
		return "", err
	}

	var resp []byte
	switch cmd.dispatcher {
	case DispatcherCmd("movewindow"):
		resp, err = mgr.moveWindow(win)
	case DispatcherCmd("movefocus"):
		resp, err = mgr.moveFocus(win)
	default:
		err = errors.New("unknown dispatcher")
	}
	return string(resp), err
}

func (m manager) moveWindow(win *WinObj) (resp []byte, err error) {
	dir := m.cmd.direction.Str()

	if m.cmd.direction.ToMonitor() {
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

func (m manager) moveFocus(win *WinObj) (resp []byte, err error) {
	dir := m.cmd.direction.Str()

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

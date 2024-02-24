# hyprwin

## Description

Simple utility that trying to make Hyprland window movements similar to i3/sway.

It moves focuses or windows in group with keybinds that you're usually use for
default `move[focus|window]` dispatcher.

For example: if window is on the groups edge and you want to move it outside,
you will no longer need separate keybind for this.

## Installation

```bash
go install github.com/leonzag/hyprwin@latest
```

>Note: Make sure you have `$GOPATH` set and `GOPATH/bin` added to your `$PATH`

## Usage

Using **cli**:

```shell
hyprwin DISPATCHER DIRECTION

Flag:
    --help          Show this message

Dispatchers:
    movefocus       Moves the focus in a direction

    movewindow      Moves the active window in a direction or to a monitor.
                    For floating windows, moves the window to the screen edge
                    in that direction
Directions:
    l,r,u,d         For left, right, up, down
    mon:<monitor>   Only for movefocus dispatcher`
```

Bind it in your `hyprland.conf` like:

```hyprlang
bind = WIN, H, exec, hyprwin movefocus l
bind = WIN, L, exec, hyprwin movefocus r
bind = WIN, K, exec, hyprwin movefocus u
bind = WIN, J, exec, hyprwin movefocus d

bind = WIN SHIFT, H, exec, hyprwin movewindow l
bind = WIN SHIFT, L, exec, hyprwin movewindow r
bind = WIN SHIFT, K, exec, hyprwin movewindow u
bind = WIN SHIFT, J, exec, hyprwin movewindow d
```

## Motivation

>Why not shell/python script?

I don't know **Go**, and writing this _"script"_ is a great excuse
to superficially learn it. \
Besides, binary will be faster and less resource-intensive.

> Why not `hy3` Hyprland plugin?

I'm quite happy with the default behavior of the `dwindle` layout,
except for the groups. \
I don't see the point of messing with a separate plugin.

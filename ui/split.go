package ui

import (
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Split struct {
	// Ratio keeps the current layout.
	// 0 is center, -1 completely to the left, 1 completely to the right.
	Ratio float32
	// Bar is the width for resizing the layout
	Bar unit.Dp

	drag    bool
	dragID  pointer.ID
	dragPos float32
}

const defaultBarWidth = unit.Dp(10)

func (s *Split) HLayout(gtx C, left, right layout.Widget) D {
	bar := gtx.Dp(s.Bar)
	if bar <= 1 {
		bar = gtx.Dp(defaultBarWidth)
	}

	proportion := (s.Ratio + 1) / 2
	leftsize := int(proportion*float32(gtx.Constraints.Max.X) - float32(bar))

	rightoffset := leftsize + bar
	rightsize := gtx.Constraints.Max.X - rightoffset

	{ // handle input
		barRect := image.Rect(leftsize, 0, rightoffset, gtx.Constraints.Max.X)
		area := clip.Rect(barRect).Push(gtx.Ops)

		// register for input
		event.Op(gtx.Ops, s)
		pointer.CursorColResize.Add(gtx.Ops)

		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: s,
				Kinds:  pointer.Press | pointer.Drag | pointer.Release | pointer.Cancel,
			})
			if !ok {
				break
			}

			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Kind {
			case pointer.Press:
				if s.drag {
					break
				}

				s.dragID = e.PointerID
				s.dragPos = e.Position.X
				s.drag = true

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}

				deltaX := e.Position.X - s.dragPos
				s.dragPos = e.Position.X

				deltaRatio := deltaX * 2 / float32(gtx.Constraints.Max.X)
				s.Ratio += deltaRatio

				if e.Priority < pointer.Grabbed {
					gtx.Execute(pointer.GrabCmd{
						Tag: s,
						ID:  s.dragID,
					})
				}

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}

		area.Pop()
	}

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftsize, gtx.Constraints.Max.Y))
		left(gtx)
	}

	{
		off := op.Offset(image.Pt(rightoffset, 0)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightsize, gtx.Constraints.Max.Y))
		right(gtx)
		off.Pop()
	}

	barRect := image.Rect(leftsize+bar/2-1, 0, leftsize+bar/2, gtx.Constraints.Max.Y)
	defer clip.Rect(barRect).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: gtx.Constraints.Max}
}

func (s *Split) VLayout(gtx C, top, bottom layout.Widget) D {
	bar := gtx.Dp(s.Bar)
	if bar <= 1 {
		bar = gtx.Dp(defaultBarWidth)
	}

	proportion := (s.Ratio + 1) / 2
	topSize := int(proportion*float32(gtx.Constraints.Max.Y) - float32(bar))

	bottomOffset := topSize + bar
	bottomSize := gtx.Constraints.Max.Y - bottomOffset

	{ // handle input
		barRect := image.Rect(0, topSize, gtx.Constraints.Max.X, bottomOffset)
		area := clip.Rect(barRect).Push(gtx.Ops)

		// register for input
		event.Op(gtx.Ops, s)
		pointer.CursorColResize.Add(gtx.Ops)
		pointer.CursorRowResize.Add(gtx.Ops)

		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: s,
				Kinds:  pointer.Press | pointer.Drag | pointer.Release | pointer.Cancel,
			})
			if !ok {
				break
			}

			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			switch e.Kind {
			case pointer.Press:
				if s.drag {
					break
				}

				s.dragID = e.PointerID
				s.dragPos = e.Position.Y
				s.drag = true

			case pointer.Drag:
				if s.dragID != e.PointerID {
					break
				}

				deltaY := e.Position.Y - s.dragPos
				s.dragPos = e.Position.Y

				deltaRatio := deltaY * 2 / float32(gtx.Constraints.Max.Y)
				s.Ratio += deltaRatio

				if e.Priority < pointer.Grabbed {
					gtx.Execute(pointer.GrabCmd{
						Tag: s,
						ID:  s.dragID,
					})
				}

			case pointer.Release:
				fallthrough
			case pointer.Cancel:
				s.drag = false
			}
		}

		area.Pop()
	}

	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, topSize))
		top(gtx)
	}

	{
		off := op.Offset(image.Pt(0, bottomOffset)).Push(gtx.Ops)
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, bottomSize))
		bottom(gtx)
		off.Pop()
	}

	barRect := image.Rect(0, topSize+bar/2-1, gtx.Constraints.Max.X, topSize+bar/2)
	defer clip.Rect(barRect).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: gtx.Constraints.Max}
}

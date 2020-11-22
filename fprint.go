package notation

import "strings"

func unwrappable(n node) bool {
	return n.len == n.wrapLen.max &&
		n.len == n.fullWrap.max
}

func initialize(t *int, v int) {
	if *t > 0 {
		return
	}

	*t = v
}

func max(t *int, v int) {
	if *t >= v {
		return
	}

	*t = v
}

func nodeLen(t int, n node) node {
	var w, f int
	for i, p := range n.parts {
		switch part := p.(type) {
		case string:
			n.len += len(part)
			w += len(part)
			f += len(part)
		case str:
			// We assume here that an str is always contained by a node that has only a
			// single str.
			//
			// If this changes in the future, then we need to provide tests for the
			// additional cases. If this doesn't change anytime soon, then we can
			// refactor this part.
			//
			n.len = len(part.val)
			if part.raw == "" {
				w = len(part.val)
				f = len(part.val)
			} else {
				lines := strings.Split(part.raw, "\n")
				part.rawLen.first = len(lines[0])
				for _, line := range lines {
					if len(line) > part.rawLen.max {
						part.rawLen.max = len(line)
					}
				}

				part.rawLen.last = len(lines[len(lines)-1])
				n.parts[i] = part
				n.wrapLen.first = part.rawLen.first
				n.fullWrap.first = part.rawLen.first
				n.wrapLen.max = part.rawLen.max
				n.fullWrap.max = part.rawLen.max
				w = part.rawLen.last
				f = part.rawLen.last
			}
		case node:
			part = nodeLen(t, part)
			n.parts[i] = part
			n.len += part.len
			if unwrappable(part) {
				w += part.len
				f += part.len
				continue
			}

			if part.len == part.wrapLen.max {
				w += part.len
			} else {
				w += part.wrapLen.first
				initialize(&n.wrapLen.first, w)
				max(&n.wrapLen.max, w)
				w = part.wrapLen.last
			}

			f += part.fullWrap.first
			initialize(&n.fullWrap.first, f)
			max(&n.fullWrap.max, f)
			f = part.fullWrap.last
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			initialize(&n.wrapLen.first, w)
			max(&n.wrapLen.max, w)
			initialize(&n.fullWrap.first, f)
			max(&n.fullWrap.max, f)
			w, f = 0, 0
			n.len += (len(part.items) - 1) * len(part.sep)
			if part.mode == line {
				w += (len(part.items) - 1) * len(part.sep)
			}

			for j, item := range part.items {
				item = nodeLen(t, item)
				part.items[j] = item
				n.len += item.len
				switch part.mode {
				case line:
					w += item.len
					max(&f, item.len)
				default:
					wj := t + item.len + len(part.suffix)
					max(&w, wj)
					fj := t + item.fullWrap.max
					max(&f, fj)
					fj = t + item.fullWrap.last + len(part.suffix)
					max(&f, fj)
				}
			}

			max(&n.wrapLen.max, w)
			max(&n.fullWrap.max, f)
			w, f = 0, 0
		}
	}

	initialize(&n.wrapLen.first, w)
	max(&n.wrapLen.max, w)
	n.wrapLen.last = w
	initialize(&n.fullWrap.first, f)
	max(&n.fullWrap.max, f)
	n.fullWrap.last = f
	return n
}

func wrapNode(t, cf0, c0, c1 int, n node) node {
	if n.len <= c0 {
		return n
	}

	if n.wrapLen.max >= n.len && n.fullWrap.max >= n.len {
		return n
	}

	if n.len <= c1 && n.len-c0 <= n.wrapLen.max {
		return n
	}

	n.wrap = true
	cc0, cc1 := c0, c1
	lastWrapperIndex := -1
	var trackBack bool
	for i := 0; i < len(n.parts); i++ {
		p := n.parts[i]
		switch part := p.(type) {
		case string:
			cc0 -= len(part)
			cc1 -= len(part)
			if !trackBack && cc1 < 0 {
				cc0 = 0
				cc1 = 0
				i = lastWrapperIndex
				trackBack = true
			}
		case str:
			// We assume here that an str is always contained by a node that has only a
			// single str. Therefore we don't need to trackback to here, because the
			// decision on wrapping was already made for the node.
			//
			// If this changes in the future, then we need to provide tests for the
			// additional cases. If this doesn't change anytime soon, then we can
			// refactor this part.
			//
			part.useRaw = part.raw != ""
			n.parts[i] = part
		case node:
			part = wrapNode(t, cf0, cc0, cc1, part)
			n.parts[i] = part
			if part.wrap {
				// This is an approximation: sometimes part.fullWrap.first should be applied
				// here, but usually those are the same.
				cc0 -= part.wrapLen.first
				cc1 -= part.wrapLen.first
			} else {
				cc0 -= part.len
				cc1 -= part.len
			}

			if !trackBack && cc1 < 0 {
				cc0 = 0
				cc1 = 0
				i = lastWrapperIndex
				trackBack = true
			}
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			cc0, cc1 = c0, c1
			trackBack = false
			lastWrapperIndex = i
			switch part.mode {
			case line:
				cl := cf0 - t
				var w int
				for j, ni := range part.items {
					if w > 0 && w+len(part.sep)+ni.len > cl {
						w = 0
						part.lineEnds = append(
							part.lineEnds,
							j,
						)
					}

					if w > 0 {
						w += len(part.sep)
					}

					w += ni.len
				}

				part.lineEnds = append(part.lineEnds, len(part.items))
				n.parts[i] = part
			default:
				for j := range part.items {
					part.items[j] = wrapNode(
						t,
						cf0,
						c0-t,
						c1-t,
						part.items[j],
					)
				}
			}
		}
	}

	return n
}

func fprint(w *writer, t int, n node) {
	if w.err != nil {
		return
	}

	for _, p := range n.parts {
		switch part := p.(type) {
		case node:
			fprint(w, t, part)
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			if !n.wrap {
				for i, ni := range part.items {
					if i > 0 {
						w.write(part.sep)
					}

					fprint(w, t, ni)
				}

				continue
			}

			switch part.mode {
			case line:
				var (
					lines [][]node
					last  int
				)

				for _, i := range part.lineEnds {
					lines = append(lines, part.items[last:i])
					last = i
				}

				for _, line := range lines {
					w.blankLine()
					w.tabs(1)
					for i, ni := range line {
						if i > 0 {
							w.write(part.sep)
						}

						fprint(w, 0, ni)
					}
				}
			default:
				t++
				for _, ni := range part.items {
					w.line(t)
					fprint(w, t, ni)
					w.write(part.suffix)
				}

				t--
			}

			w.line(t)
		default:
			w.write(part)
		}
	}
}

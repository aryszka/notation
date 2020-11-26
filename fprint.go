package notation

import (
	"fmt"
	"strings"
)

func ifZero(a, b int) int {
	if a == 0 {
		return b
	}

	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func strLen(s str) str {
	l := strings.Split(s.raw, "\n")
	for j, li := range l {
		if j == 0 {
			s.rawLen.first = len(li)
		}

		if len(li) > s.rawLen.max {
			s.rawLen.max = len(li)
		}

		if j == len(l)-1 {
			s.rawLen.last = len(li)
		}
	}

	return s
}

func nodeLen(t int, n node) node {
	// We assume here that an str is always contained
	// by a node that has only a single str.
	//
	if s, ok := n.parts[0].(str); ok {
		s = strLen(s)
		n.parts[0] = s
		n.len = len(s.val)
		if s.raw == "" {
			wl := wrapLen{
				first: len(s.val),
				max:   len(s.val),
				last:  len(s.val),
			}

			n.wrapLen = wl
			n.fullWrap = wl
			return n
		}

		n.wrapLen = s.rawLen
		n.fullWrap = s.rawLen
		return n
	}

	// measure all parts:
	for i := range n.parts {
		switch pt := n.parts[i].(type) {
		case node:
			n.parts[i] = nodeLen(t, pt)
		case wrapper:
			for j := range pt.items {
				pt.items[j] = nodeLen(t, pt.items[j])
			}
		}
	}

	// measure the unwrapped length:
	for _, p := range n.parts {
		switch pt := p.(type) {
		case node:
			n.len += pt.len
		case wrapper:
			if len(pt.items) == 0 {
				continue
			}

			n.len += (len(pt.items) - 1) * len(pt.sep)
			for _, pti := range pt.items {
				n.len += pti.len
			}
		default:
			n.len += len(fmt.Sprint(p))
		}
	}

	// measure the wrapped and the fully wrapped length:
	var w, f int
	for _, p := range n.parts {
		switch pt := p.(type) {
		case node:
			w += pt.wrapLen.first
			if pt.len != pt.wrapLen.first {
				n.wrapLen.first = ifZero(n.wrapLen.first, w)
				n.wrapLen.max = max(n.wrapLen.max, w)
				n.wrapLen.max = max(n.wrapLen.max, pt.wrapLen.max)
				w = pt.wrapLen.last
			}

			f += pt.fullWrap.first
			if pt.len != pt.fullWrap.first {
				n.fullWrap.first = ifZero(n.fullWrap.first, f)
				n.fullWrap.max = max(n.fullWrap.max, f)
				n.fullWrap.max = max(n.fullWrap.max, pt.fullWrap.max)
				f = pt.fullWrap.last
			}
		case wrapper:
			if len(pt.items) == 0 {
				continue
			}

			n.wrapLen.first = ifZero(n.wrapLen.first, w)
			n.wrapLen.max = max(n.wrapLen.max, w)
			n.fullWrap.first = ifZero(n.fullWrap.first, f)
			n.fullWrap.max = max(n.fullWrap.max, f)
			w = 0
			f = 0
			switch pt.mode {
			case line:
				// line wrapping is flexible, here
				// we measure the longest case
				//
				w = (len(pt.items) - 1) * len(pt.sep)
				for _, pti := range pt.items {
					w += pti.len
				}

				// here me measure the shortest
				// possible case
				//
				for _, pti := range pt.items {
					f = max(f, pti.fullWrap.max)
				}
			default:
				// for non-full wrap, we measure the full
				// length of the items
				//
				for _, pti := range pt.items {
					w = max(w, t+pti.len+len(pt.suffix))
				}

				// for full wrap, we measure the fully
				// wrapped length of the items
				//
				for _, pti := range pt.items {
					f = max(f, t+pti.fullWrap.max)
					f = max(f, t+pti.fullWrap.last+len(pt.suffix))
				}
			}

			n.wrapLen.max = max(n.wrapLen.max, w)
			n.fullWrap.max = max(n.fullWrap.max, f)
			w = 0
			f = 0
		default:
			w += len(fmt.Sprint(p))
			f += len(fmt.Sprint(p))
		}
	}

	n.wrapLen.first = ifZero(n.wrapLen.first, w)
	n.wrapLen.max = max(n.wrapLen.max, w)
	n.wrapLen.last = w
	n.fullWrap.first = ifZero(n.fullWrap.first, f)
	n.fullWrap.max = max(n.fullWrap.max, f)
	n.fullWrap.last = f
	return n
}

func wrapNode(t, cf0, c0, c1 int, n node) node {
	// fits:
	if n.len <= c0 {
		return n
	}

	// we don't want to make it longer:
	if n.wrapLen.max >= n.len && n.fullWrap.max >= n.len {
		return n
	}

	// tolerate below c1 when it's not worth wrapping:
	if n.len <= c1 && n.len-c0 <= c0-n.wrapLen.max {
		return n
	}

	// otherwise, we need to wrap the node:
	n.wrap = true

	// We assume here that an str is always contained
	// by a node that has only a single str.
	//
	if s, ok := n.parts[0].(str); ok {
		s.useRaw = s.raw != ""
		n.parts[0] = s
		return n
	}

	// before iterating over the parts, take a copy of
	// the available column width and modify only the
	// copy, to support trackback.
	//
	cc0, cc1 := c0, c1
	lastWrapperIndex := -1
	var trackBack bool
	for i := 0; i < len(n.parts); i++ {
		p := n.parts[i]
		switch part := p.(type) {
		case node:
			part = wrapNode(t, cf0, cc0, cc1, part)
			n.parts[i] = part
			if part.wrap {
				// This is an approximation: sometimes
				// part.fullWrap.last should be applied
				// here, but usually those are the same.
				//
				cc0 -= part.wrapLen.first
				cc1 -= part.wrapLen.first
			} else {
				cc0 -= part.len
				cc1 -= part.len
			}

			if cc1 >= 0 {
				if part.wrap {
					cc0 = c0 - part.wrapLen.last
					cc1 = c1 - part.wrapLen.last
				}

				continue
			}

			if trackBack {
				continue
			}

			// trackback from after the last wrapper:
			i = lastWrapperIndex
			trackBack = true

			// force wrapping during trackback:
			cc0 = 0
			cc1 = 0
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			cc0, cc1 = c0, c1
			trackBack = false
			lastWrapperIndex = i
			switch part.mode {
			case line:
				// we only set the line endings. We use
				// the full column width:
				//
				cl := cf0 - t
				var w int
				for j, nj := range part.items {
					if w > 0 && w+len(part.sep)+nj.len > cl {
						part.lineEnds = append(part.lineEnds, j)
						w = 0
					}

					if w > 0 {
						w += len(part.sep)
					}

					w += nj.len
				}

				part.lineEnds = append(part.lineEnds, len(part.items))
				n.parts[i] = part
			default:
				for j := range part.items {
					part.items[j] = wrapNode(t, cf0, c0-t, c1-t, part.items[j])
				}
			}
		default:
			s := fmt.Sprint(part)
			cc0 -= len(s)
			cc1 -= len(s)
			if cc1 >= 0 {
				continue
			}

			if trackBack {
				continue
			}

			// trackback from after the last wrapper:
			i = lastWrapperIndex
			trackBack = true

			// force wrapping during trackback:
			cc0 = 0
			cc1 = 0
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
					w.line(1)
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

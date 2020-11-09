package notation

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

func wrapNode(t, c0, c1 int, n node) node {
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
	for i, p := range n.parts {
		switch part := p.(type) {
		case node:
			part = wrapNode(t, cc0, cc1, part)
			n.parts[i] = part
			if part.wrap {
				cc0 -= part.wrapLen.last
				cc1 -= part.wrapLen.last
			} else {
				cc0 -= part.len
				cc1 -= part.len
			}
		case wrapper:
			if len(part.items) > 0 {
				cc0, cc1 = c0, c1
			}

			switch part.mode {
			case line:
				c0, c1 = c0-t, c1-t
				var w int
				for j, ni := range part.items {
					if w > 0 && w+len(part.sep)+ni.len > c0 {
						w = 0
						part.lineWrappers = append(
							part.lineWrappers,
							j,
						)
					}

					if w > 0 {
						w += len(part.sep)
					}

					w += ni.len
				}

				n.parts[i] = part
				c0, c1 = c0+t, c1+t
			default:
				for j := range part.items {
					part.items[j] = wrapNode(
						t,
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
					fprint(w, t, ni)
					if i < len(part.items)-1 {
						w.write(part.sep)
					}
				}

				continue
			}

			t++
			switch part.mode {
			case line:
				var (
					wi          int
					lineStarted bool
				)

				w.line(t)
				for i, ni := range part.items {
					if len(part.lineWrappers) > wi &&
						i == part.lineWrappers[wi] {
						wi++
						w.line(t)
						lineStarted = false
					}

					if lineStarted {
						w.write(part.sep)
					}

					fprint(w, 0, ni)
					lineStarted = true
				}
			default:
				for _, ni := range part.items {
					w.line(t)
					fprint(w, t, ni)
					w.write(part.suffix)
				}
			}

			t--
			w.line(t)
		default:
			w.write(part)
		}
	}
}

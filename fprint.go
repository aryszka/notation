package notation

import (
	"fmt"
	"io"
)

type writer struct {
	w   io.Writer
	n   int
	err error
}

func (w *writer) write(o interface{}) {
	if w.err != nil {
		return
	}

	n, err := fmt.Fprint(w.w, o)
	w.n += n
	w.err = err
}

func nodeLen(t int, n node) node {
	var w int
	for i, p := range n.parts {
		switch part := p.(type) {
		case string:
			n.len += len(part)
			w += len(part)
		case node:
			part = nodeLen(t, part)
			n.parts[i] = part
			n.len += part.len
			if part.len == part.wlen {
				w += part.len
			} else {
				w += part.wlen0
				if w > n.wlen {
					n.wlen = w
				}

				if part.wlen > n.wlen {
					n.wlen = part.wlen
				}

				if n.wlen0 == 0 {
					n.wlen0 = w
				}

				w = part.wlenLast
			}
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			if w > n.wlen {
				n.wlen = w
			}

			if n.wlen0 == 0 {
				n.wlen0 = w
			}

			w = 0
			for j, ni := range part.items {
				ni = nodeLen(t, ni)
				part.items[j] = ni
				n.len += ni.len
				wni := t + ni.len + len(part.suffix)
				if wni > w {
					w = wni
				}
			}

			if len(part.items) > 0 {
				n.len += (len(part.items) - 1) * len(part.sep)
			}

			w = 0
		}
	}

	if w > n.wlen {
		n.wlen = w
	}

	if n.wlen0 == 0 {
		n.wlen0 = w
	}

	n.wlenLast = w
	return n
}

func wrapNode(t, c0, c1 int, n node) node {
	if n.len <= c0 || n.wlen == n.len {
		return n
	}

	if n.len <= c1 && n.len-c0 <= n.wlen {
		return n
	}

	n.wrap = true
	if n.wlen <= c0 {
		return n
	}

	for i, p := range n.parts {
		switch part := p.(type) {
		case node:
			n.parts[i] = wrapNode(t, c0, c1, part)
		case wrapper:
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

	return n
}

func fprint(w *writer, t int, n node) {
	if w.err != nil {
		return
	}

	for i := 0; i < t; i++ {
		w.write("\t")
	}

	for _, p := range n.parts {
		switch part := p.(type) {
		case node:
			fprint(w, 0, part)
		case wrapper:
			if len(part.items) == 0 {
				continue
			}

			if n.wrap {
				w.write("\n")
			}

			for i, ni := range part.items {
				if n.wrap {
					fprint(w, t+1, ni)
					w.write(part.suffix)
					w.write("\n")
				} else {
					fprint(w, 0, ni)
					if i < len(part.items)-1 {
						w.write(part.sep)
					}
				}
			}
		default:
			w.write(part)
		}
	}
}

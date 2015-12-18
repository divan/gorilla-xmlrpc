// Copyright 2015 Tamás Gulácsi
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xml

import (
	"bytes"
	"sync"
)

var bufPool = &bufferPool{
	Pool: sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 1024)) }},
}

type bufferPool struct {
	sync.Pool
}

func (p *bufferPool) Get() *bytes.Buffer {
	return p.Pool.Get().(*bytes.Buffer)
}
func (p *bufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	p.Pool.Put(b)
}

// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"

	"github.com/jlourenc/xgo/xerrors"
)

func ExampleNew() {
	err := xerrors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}
	// Output: emit macho dwarf: elf header corrupted
}

func ExampleNewf() {
	const name, id = "bimmler", 17
	err := xerrors.Newf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}
	// Output: user "bimmler" (id 17) not found
}

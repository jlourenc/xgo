// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"
	"os"

	"github.com/jlourenc/xgo/xerrors"
)

func ExampleAppend() {
	vetOperands := func(a, b int) (errs error) {
		if a < 0 {
			errs = xerrors.Append(errs, xerrors.New("left operand is negative"))
		}
		if b < 0 {
			errs = xerrors.Append(errs, xerrors.New("right operand is negative"))
		}
		return
	}

	var errs error

	if err := vetOperands(-1, -2); err != nil {
		errs = xerrors.Append(errs, xerrors.Wrapf(err, "failed to vet operands %d and %d", 0, -2))
	}

	if err := vetOperands(0, -2); err != nil {
		errs = xerrors.Append(errs, xerrors.Wrapf(err, "failed to vet operands %d and %d", 0, 0))
	}

	if errs != nil {
		fmt.Print(errs)
	}

	// Output:
	// 2 errors occurred:
	//	* failed to vet operands 0 and -2: 2 errors occurred:
	// 		* left operand is negative
	// 		* right operand is negative
	// 	* failed to vet operands 0 and 0: 1 error occurred:
	// 		* right operand is negative
}

func ExampleAs() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *os.PathError
		if xerrors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}

	// Output: Failed at path: non-existing
}

func ExampleIs() {
	if _, err := os.Open("non-existing"); err != nil {
		if xerrors.Is(err, os.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

	// Output: file does not exist
}

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

func ExampleUnwrap() {
	err := xerrors.New("elf header corrupted")
	err = xerrors.Wrap(err, "emit macho dwarf")
	err = xerrors.Unwrap(err)
	fmt.Println(err)

	// Output: elf header corrupted
}

func ExampleWithStack() {
	err := xerrors.WithStack(os.ErrNotExist)
	fmt.Println(err)

	// Output: file does not exist
}

func ExampleWrap() {
	err := xerrors.Wrap(os.ErrNotExist, "failed to open file")
	fmt.Println(err)

	// Output: failed to open file: file does not exist
}

func ExampleWrapf() {
	err := xerrors.Wrapf(os.ErrNotExist, "failed to open file %s", "myfile")
	fmt.Println(err)

	// Output: failed to open file myfile: file does not exist
}

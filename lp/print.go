package lp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (dict *Dict) longestCoeff(format string) int {
	var n int
	for _, ai := range dict.A {
		for _, aij := range ai {
			n = max(n, len(fmt.Sprintf(format, aij)))
		}
	}
	for _, bi := range dict.B {
		n = max(n, len(fmt.Sprintf(format, bi)))
	}
	for _, ci := range dict.C {
		n = max(n, len(fmt.Sprintf(format, ci)))
	}
	n = max(n, len(fmt.Sprintf(format, dict.D)))
	return n
}

func (dict *Dict) longestIndex(format string) int {
	var n int
	for _, ni := range dict.Basic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	for _, ni := range dict.NonBasic {
		n = max(n, len(fmt.Sprintf(format, ni)))
	}
	return n
}

func (dict *Dict) Fprint(w io.Writer) error {
	coeffLen := dict.longestCoeff("%+-.4g")
	coeff := "%+-" + fmt.Sprintf("%d", coeffLen) + ".4g"
	indexLen := dict.longestIndex("%d")
	index := "%" + fmt.Sprintf("%d", indexLen) + "d"

	var b bytes.Buffer
	for i := range dict.Basic {
		fmt.Fprintf(&b, "x"+index+" =", dict.Basic[i])
		fmt.Fprintf(&b, " "+coeff, dict.B[i])
		for j := range dict.NonBasic {
			fmt.Fprintf(&b, "  "+coeff+" x"+index, dict.A[i][j], dict.NonBasic[j])
		}
		b.WriteString("\n")
		if _, err := io.Copy(w, &b); err != nil {
			return err
		}
	}

	spacer := strings.Repeat(" ", indexLen)
	fmt.Fprint(&b, "z"+spacer+" =")
	fmt.Fprintf(&b, " "+coeff, dict.D)
	for j := range dict.NonBasic {
		fmt.Fprintf(&b, "  "+coeff+" x"+index, dict.C[j], dict.NonBasic[j])
	}
	b.WriteString("\n")
	if _, err := io.Copy(w, &b); err != nil {
		return err
	}
	return nil
}

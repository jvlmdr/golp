package lp

// Dual returns the corresponding dual dictionary.
//	A, b, c, d = -A', -c, -b, -d
// Dual variables have the label of their primal complement.
func (p *Dict) Dual() *Dict {
	m, n := len(p.Basic), len(p.NonBasic)
	d := NewDict(n, m)

	// Label dual variables with their primal complement.
	copy(d.Basic, p.NonBasic)
	copy(d.NonBasic, p.Basic)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			d.A[j][i] = -p.A[i][j]
		}
	}
	for j := 0; j < n; j++ {
		d.B[j] = -p.C[j]
	}
	for i := 0; i < m; i++ {
		d.C[i] = -p.B[i]
	}
	d.D = -p.D

	return d
}

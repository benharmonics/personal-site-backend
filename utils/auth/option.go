package auth

import "runtime"

type HashedPasswordOption func(*Argon2Params)

func WithCPUThreads(ps *Argon2Params) { ps.threads = uint8(runtime.NumCPU()) }

func newArgon2Params(opts ...HashedPasswordOption) Argon2Params {
	ps := Argon2Params{
		time:       1,
		memory:     64 * 1024, // 64MB
		threads:    4,
		keyLength:  32,
		saltLength: 32,
	}
	for _, f := range opts {
		f(&ps)
	}
	return ps
}

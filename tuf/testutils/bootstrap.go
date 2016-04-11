package testutils

type TestBootstrapper struct {
	Booted bool
}

func (tb *TestBootstrapper) Bootstrap() error {
	tb.Booted = true
	return nil
}

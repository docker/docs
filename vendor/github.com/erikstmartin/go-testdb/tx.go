package testdb

type Tx struct {
	commitFunc   func() error
	rollbackFunc func() error
}

func (t *Tx) Commit() error {
	if t.commitFunc != nil {
		return t.commitFunc()
	}
	return nil
}

func (t *Tx) Rollback() error {
	if t.rollbackFunc != nil {
		return t.rollbackFunc()
	}
	return nil
}

func (t *Tx) SetCommitFunc(f func() error) {
	t.commitFunc = f
}

func (t *Tx) StubCommitError(err error) {
	t.SetCommitFunc(func() error {
		return err
	})
}

func (t *Tx) SetRollbackFunc(f func() error) {
	t.rollbackFunc = f
}

func (t *Tx) StubRollbackError(err error) {
	t.SetRollbackFunc(func() error {
		return err
	})
}

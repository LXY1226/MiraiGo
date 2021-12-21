package utils

func CoverError(fun func()) error {
	errCh := make(chan error, 1)
	func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- r.(error)
			}
		}()
		fun()
		errCh <- nil
	}()
	return <-errCh
}

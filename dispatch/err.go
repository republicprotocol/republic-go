package dispatch

func MergeErrors(errChChIn <-chan <-chan error) <-chan error {
	errCh := make(chan error)
	return errCh
}

func FilterCancelErrors(errChIn <-chan error) <-chan error {
	errCh := make(chan error)
	return errCh
}

func ConsumeErrors(errCh <-chan error, f ...func(err error)) {
	for err := range errCh {
		if f != nil && len(f) > 0 {
			for i := range f {
				f[i](err)
			}
		}
	}
}

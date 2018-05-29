package ome

func ConfirmOrderMatches(done <-chan struct{}, orderMatches <-chan Computation) (<-chan Computation, <-chan error) {
	confirmedOrderMatches := make(chan Computation)
	errs := make(chan error)

	go func() {
		defer close(confirmedOrderMatches)
		defer close(errs)
	}()

	return confirmedOrderMatches, errs
}

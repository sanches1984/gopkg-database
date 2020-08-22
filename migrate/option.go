package migrate

type Option []OptionFn

type OptionFn func(r *Runner)

func WithDryRun() OptionFn {
	return func(r *Runner) {
		r.dryRun = true
	}
}

func WithClean(scheme ...string) OptionFn {
	return func(r *Runner) {
		r.cleanScheme = scheme
	}
}

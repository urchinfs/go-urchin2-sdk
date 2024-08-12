package options

type DagImportSettings struct {
	PinRoots bool
	Silent   bool
	Stats    bool
}

type DagImportOption func(opts *DagImportSettings) error

func DagImportOptions(opts ...DagImportOption) (*DagImportSettings, error) {
	options := &DagImportSettings{
		PinRoots: false,
		Silent:   false,
		Stats:    false,
	}

	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}

	return options, nil
}

func (dagOpts) PinRoots(pinRoots bool) DagImportOption {
	return func(opts *DagImportSettings) error {
		opts.PinRoots = pinRoots
		return nil
	}
}

func (dagOpts) Silent(silent bool) DagImportOption {
	return func(opts *DagImportSettings) error {
		opts.Silent = silent
		return nil
	}
}

func (dagOpts) Stats(stats bool) DagImportOption {
	return func(opts *DagImportSettings) error {
		opts.Stats = stats
		return nil
	}
}

package spoon

type Option func(*optionParam) error

func TagPrefix(tp string) Option {
	return func(p *optionParam) error {
		p.tagPrefix = tp
		return nil
	}
}

func IgnoreTag(ig string) Option {
	return func(p *optionParam) error {
		p.ignoreTag = ig
		return nil
	}
}

package xtracer

var provider Provider

func WithVendor(p Provider) {
	provider = p
}

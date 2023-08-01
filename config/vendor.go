package config

var (
	vendors          = make(map[string]Vendor)
	propertiesVendor Properties
)

type Vendor interface {
	Read() error
	Init(protocol string) error
}

func Register(scheme string, creator Vendor) {
	vendors[scheme] = creator
}

func WithVendor(v Properties) {
	propertiesVendor = v
}

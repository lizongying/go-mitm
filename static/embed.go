package static

import (
	"embed"
	_ "embed"
)

//go:embed tls/ca.crt
var CaCert []byte

//go:embed tls/ca.key
var CaKey []byte

//go:embed dist
var Dist embed.FS

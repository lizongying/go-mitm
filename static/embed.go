package static

import (
	"embed"
	_ "embed"
)

//go:embed tls/ca_crt.pem
var CaCert []byte

//go:embed tls/ca_key.pem
var CaKey []byte

//go:embed dist
var Dist embed.FS

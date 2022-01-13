package config

import "github.com/spf13/afero"

// Permits overriding in tests so we can mock the file system
var Appfs = afero.NewOsFs()

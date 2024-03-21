package featureflags

import "os"

// Enable fake fills for market maker testing purposes
var FlagEnableFakeFill = os.Getenv("ENABLE_FAKE_FILL")

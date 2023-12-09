package featureflags

import "os"

// Enable/disable signature verification
//
// It is optional whether to verify the signature or not in the order book, as it will later be verified on-chain regardless.
var ShouldVerifySig = os.Getenv("VERIFY_SIGNATURE")

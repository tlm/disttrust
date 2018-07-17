package provider

// Represents a new certificate response from a provider. This type keeps it
// simple until more complex cases can be identified
type Response struct {
	// Certificate authority data
	CA string
	// Certificate data
	Certificate string

	// Certificate and CA bundle combined
	CABundle string

	// Private key data
	PrivateKey string

	// Serial of certificate
	Serial string
}

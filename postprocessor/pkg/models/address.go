package models

// AddressDetails contains data related to an SMC's address.
type AddressDetails struct {
	ShortAddress    int
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}

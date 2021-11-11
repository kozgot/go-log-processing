package models

// AddressDetails contains data related to an SMC's address.
type AddressDetails struct {
	ShortAddress    int
	PhysicalAddress string
	LogicalAddress  string
	URL             string
}

// Equals check equality.
func (a *AddressDetails) Equals(other AddressDetails) bool {
	return a.ShortAddress == other.ShortAddress && a.PhysicalAddress == other.PhysicalAddress &&
		a.LogicalAddress == other.LogicalAddress && a.URL == other.URL
}

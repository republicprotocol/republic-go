package crypto

// A Keystore is an optionally encrypted storage struct for the keys used by
// Republic Protocol. It marshals, and unmarshals, to, and from, JSON.
type Keystore struct {
}

// MarshalJSON implements the json.Marshaler interface.
func (keystore Keystore) MarshalJSON() ([]byte, error) {
	panic("unimplemented")
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (keystore *Keystore) UnmarshalJSON(data []byte) error {
	panic("unimplemented")
}

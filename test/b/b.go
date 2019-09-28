package b

// Fuga is a struct
type Fuga struct {
	Piyo string
}

// GetPiyo returns f.Piyo
func (f *Fuga) GetPiyo() string {
	return f.Piyo
}

// NewFuga returns a new Fuga
func NewFuga() *Fuga {
	return &Fuga{"piyo"}
}

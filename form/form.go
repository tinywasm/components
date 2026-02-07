package form

// InputType defines supported input types
type InputType string

const (
	Text     InputType = "text"
	Password InputType = "password"
	Email    InputType = "email"
	Number   InputType = "number"
	Date     InputType = "date"
	Hidden   InputType = "hidden"
)

// Form represents a form definition
type Form struct {
	ID           string
	Name         string
	Autocomplete bool
	Fieldsets    []*Fieldset
	// Additional attributes can be added here
}

// Fieldset represents a group of fields
type Fieldset struct {
	Legend   string
	TabIndex int
	Fields   []Field
	Class    string
}

// Field interface for different types of form elements
type Field interface {
	Render() string
}

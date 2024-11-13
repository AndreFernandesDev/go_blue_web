package types

type ViewData struct {
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

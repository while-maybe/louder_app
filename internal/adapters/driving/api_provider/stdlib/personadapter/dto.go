package personadapter

// CreatePersonRequest defines the expected JSON payload for creating a person.
// This is a Data Transfer Object (DTO) for the HTTP layer.
type CreatePersonRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// PersonResponse defines the JSON payload for returning a person just created (inc UUID).
// This is also DTO for the HTTP layer.
type PersonResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`

	// TODO - implement later
	// Pets             []domain.Pet `json:"pets,omitempty"`
	// BirthCountry     domain.Country `json:"birth_country,omitempty"`
	// ResidentCountry  domain.Country `json:"residence_country,omitempty"`
	// VisitedCountries []domain.Country `json:"visited_countries,omitempty"`
}

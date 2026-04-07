package model

// TaxonResponse is the response for a successful taxon lookup.
type TaxonResponse struct {
	TaxID string `json:"tax_id" example:"9606"`
	Name  string `json:"name"   example:"Homo sapiens"`
}

// HealthResponse is the response for the health endpoint.
type HealthResponse struct {
	Status        string `json:"status"          example:"ok"`
	TaxonCount int    `json:"taxon_count"  example:"2498561"`
}

// ErrorResponse is returned on error.
type ErrorResponse struct {
	Error string `json:"error" example:"tax_id not found"`
}

package models

// DomainAvailableResponse represent the domain available response
type DomainAvailableResponse struct {
	Available  bool   `json:"available"`
	Currency   string `json:"currency"`
	Definitive bool   `json:"definitive"`
	Domain     string `json:"domain"`
	Period     int32  `json:"period"`
	Price      int32  `json:"price"`
}

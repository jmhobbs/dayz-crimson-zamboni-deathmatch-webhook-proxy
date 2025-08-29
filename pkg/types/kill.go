package types

type Kill struct {
	Victim   string `json:"victim"`
	Killer   string `json:"killer"`
	Weapon   string `json:"weapon"`
	Distance int    `json:"distance"`
}

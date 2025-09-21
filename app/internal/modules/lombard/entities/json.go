package entities

type ResponseDistributionData struct {
	Status    bool          `json:"success"`
	Recipient RecipientData `json:"recipient"`
}

type RecipientData struct {
	Address    string `json:"address"`
	OrgSlug    string `json:"org_slug"`
	Claimable  string `json:"total_claimable"`
	AcceptedAt string `json:"terms_accepted_at"`
	Msg        string `json:"terms_acceptance_message"`
}

type SignConfirmationData struct {
	Address string `json:"address"`
	Slug    string `json:"org_slug"`
	Sign    string `json:"signature"`
}

type SignConfirmationResponse struct {
	Status bool `json:"success"`
}

type HardClaimResponse struct {
	Claims []HardClaimData `json:"claims"`
}

type HardClaimData struct {
	ID          string    `json:"id"`
	Slug        string    `json:"org_slug"`
	Recipient   string    `json:"recipient_address"`
	Amount      string    `json:"amount"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Proofs      ProofData `json:"proof"`
	LeafIndex   string    `json:"leaf_index"`
	MerkleRoot  string    `json:"merkle_root"`
	Phase       int       `json:"phase"`
	Distributor string    `json:"distributor_address"`
}

type ProofData struct {
	Hashes []string `json:"hashes"`
}

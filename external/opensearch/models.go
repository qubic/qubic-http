package opensearch

type Response struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Found  bool        `json:"found"`
	Source interface{} `json:"_source"`
}

type TickDataResponse struct {
	Computor       int           `json:"computor"`
	Epoch          int           `json:"epoch"`
	Tick           uint32        `json:"tick"`
	Time           []int         `json:"time"`
	VarStruct      string        `json:"varStruct"`
	Timelock       string        `json:"timelock"`
	Signature      string        `json:"sig"`
	NumTx          int           `json:"numtx"`
	TransactionIDs []string      `json:"txids"`
	PotentialBxs   []PotentialBx `json:"potentialBx"`
}

type PotentialBx struct {
	Index       int    `json:"index"`
	Destination string `json:"dest"`
	Amount      string `json:"amount"`
}

type TxResponse struct {
	BxID        string `json:"bxid"`
	Utime       string `json:"utime"`
	Epoch       string `json:"epoch"`
	Tick        string `json:"tick"`
	Type        string `json:"type"`
	Source      string `json:"src"`
	Destination string `json:"dest"`
	Amount      string `json:"amount"`
	Extra       string `json:"extra"`
	Signature   string `json:"sig"`
}

type BxResponse struct {
	TxID        string `json:"txid"`
	Utime       string `json:"utime"`
	Epoch       string `json:"epoch"`
	Tick        string `json:"tick"`
	Type        string `json:"type"`
	Source      string `json:"src"`
	Destination string `json:"dest"`
	Amount      string `json:"amount"`
}

type StatusResponse struct {
	Epoch     int   `json:"epoch"`
	BxID      []int `json:"bxid"`
	TxID      []int `json:"txid"`
	Quorum    []int `json:"quorum"`
	Tick      []int `json:"tick"`
	Latest    []int `json:"latest"`
	Validate  []int `json:"validate"`
}

type QuorumResponse struct {
	Computor                      int      `json:"computor"`
	Epoch                         int      `json:"epoch"`
	Tick                          uint32   `json:"tick"`
	Time                          []int    `json:"time"`
	PreviousResourceTestingDigest string   `json:"prevRTD"`
	SaltedResourceTestingDigest   string   `json:"saltRTD"`
	Digests                       []string `json:"digests"`
	Signature                     string   `json:"sig"`
	Diffs                         []Diff   `json:"diffs"`
	NumVotes                      int      `json:"numvotes"`
}

type Diff struct {
	Computor                    int      `json:"computor"`
	SaltedResourceTestingDigest string   `json:"saltRTD"`
	Digests                     []string `json:"digests"`
	Signature                   string   `json:"sig"`
}

type ComputorsResponse struct {
	Epoch      string   `json:"epoch"`
	Identities []string `json:"pubkeys"`
	Signature  string   `json:"sig"`
}

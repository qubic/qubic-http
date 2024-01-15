package opensearch

type Response struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Found  bool        `json:"found"`
	Source interface{} `json:"_source"`
}

type TickDataResponse struct {
	Computor       int      `json:"computor"`
	Epoch          int      `json:"epoch"`
	Tick           uint32   `json:"tick"`
	Time           []int    `json:"time"`
	VarStruct      string   `json:"varStruct"`
	Timelock       string   `json:"timelock"`
	Signature      string   `json:"sig"`
	NumTx          int      `json:"numtx"`
	TransactionIDs []string `json:"txids"`
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
	Utime       string `json:"utime"`
	Epoch       string `json:"epoch"`
	Tick        string `json:"tick"`
	Type        string `json:"type"`
	Source      string `json:"src"`
	Destination string `json:"dest"`
	Amount      string `json:"amount"`
}

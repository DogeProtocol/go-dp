package qcreadapi

type BlockDetails struct {
	// The block number as of which the details were retrieved
	BlockNumber *int64 `json:"blockNumber,omitempty"`
}

// AssertBlockDetailsRequired checks if the required fields are not zero-ed
func AssertBlockDetailsRequired(obj BlockDetails) error {
	return nil
}

// AssertBlockDetailsConstraints checks if the values respects the defined constraints
func AssertBlockDetailsConstraints(obj BlockDetails) error {
	return nil
}

package qcreadapi

type LatestBlockDetailsResponse struct {
	Result BlockDetails `json:"result,omitempty"`
}

// AssertLatestBlockDetailsResponseRequired checks if the required fields are not zero-ed
func AssertLatestBlockDetailsResponseRequired(obj LatestBlockDetailsResponse) error {
	if err := AssertBlockDetailsRequired(obj.Result); err != nil {
		return err
	}
	return nil
}

// AssertLatestBlockDetailsResponseConstraints checks if the values respects the defined constraints
func AssertLatestBlockDetailsResponseConstraints(obj LatestBlockDetailsResponse) error {
	if err := AssertBlockDetailsConstraints(obj.Result); err != nil {
		return err
	}
	return nil
}

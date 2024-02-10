package match

type MockMatcher struct{}

func (m *MockMatcher) Match(_ string, _ string) error {
	return nil
}

func (m *MockMatcher) Suggest(_ string) ([]UserMatch, error) {
	return []UserMatch{}, nil
}

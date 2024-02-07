package match

type MockMatcher struct{}

func (m *MockMatcher) Match(_ string, _ string) error {
	return nil
}

func (m *MockMatcher) Suggest(id string) ([]UserMatch, error) {
	return []UserMatch{}, nil
}

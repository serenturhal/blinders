package match

type MockMatcher struct{}

func (m *MockMatcher) Match(_ string, _ string) error {
	return nil
}

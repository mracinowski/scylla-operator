package selectors

type Matcher interface {
	Match(any, any) bool
}

type matcher[L, R any] struct {
	lambda func(L, R) bool
}

func NewMatcher[L, R any](lambda func(L, R) bool) Matcher {
	return &matcher[L, R]{lambda: lambda}
}

func (m *matcher[L, R]) Match(lhs any, rhs any) bool {
	// TODO: Maybe do not panic when types are mismatched?
	return m.lambda(lhs.(L), rhs.(R))
}

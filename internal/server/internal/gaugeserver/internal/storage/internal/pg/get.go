package pg

import "context"

func (s *storage) GetGauger(ctx context.Context, name string) (float64, bool) {
	row := s.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)
	var currentVal float64

	err := row.Scan(&currentVal)
	if err != nil {
		return 0.0, false
	}
	return currentVal, true
}

func (s *storage) GetCounter(ctx context.Context, name string) (int64, bool) {
	row := s.QueryRowContext(ctx, "SELECT value FROM counter WHERE name = $1", name)
	var currentVal int64

	err := row.Scan(&currentVal)
	if err != nil {
		return 0, false
	}
	return currentVal, true
}

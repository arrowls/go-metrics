package repository

type Metric interface {
	AddGaugeValue(key string, value float64)
	AddCounterValue(key string, value int64)
}

type Repository struct {
	Metric Metric
}

func NewRepository() *Repository {
	return &Repository{
		Metric: NewMetricRepository(),
	}
}

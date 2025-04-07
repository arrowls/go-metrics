package collector

type MetricProvider interface {
	Collect()
	AsMap() *map[string]interface{}
}

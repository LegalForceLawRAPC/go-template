package items

type Service interface{}

type itemsSvc struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &itemsSvc{repo: r}
}

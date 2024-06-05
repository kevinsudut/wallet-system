package usecasetransaction

type usecase struct{}

func Init() UsecaseItf {
	return &usecase{}
}

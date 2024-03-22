package main

type Service interface {
	Serve() error
}

type GroupService struct {
	Services []Service
}

func (s *GroupService) run() {
	for _, service := range s.Services {
		service.Serve()
	}
}

func (s *GroupService) registerService(t Service) {
	s.Services = append(s.Services, t)
}

func main() {
	g := &GroupService{}
	g.registerService(NewOrderService(ConfigOrderService{ServerUrl: "localhost", ServerPort: "3000"}))
	g.registerService(&EstimateCostStreamService{})
	g.run()
}

package service

type Service interface {
	HealthCheck() bool
}

type CommentService struct {}

func (s *CommentService)HealthCheck() bool {
	return true
}
func NewCommentService() *CommentService {
	return &CommentService{}
}
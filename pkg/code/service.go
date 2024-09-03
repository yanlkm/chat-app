package code

import "context"

type CodeService interface {
	CreateCode(ctx context.Context, code *CodeEntity) error
	UpdateCode(ctx context.Context, codeString *string) error
	CheckCode(ctx context.Context, codeString *string) (bool, error)
}

type codeService struct {
	repo CodeRepository
}

func NewCodeService(repo CodeRepository) CodeService {
	return &codeService{repo: repo}
}

func (c *codeService) CreateCode(ctx context.Context, code *CodeEntity) error {
	return c.repo.Create(ctx, code)
}

func (c *codeService) UpdateCode(ctx context.Context, codeString *string) error {
	return c.repo.Update(ctx, codeString)
}

func (c *codeService) CheckCode(ctx context.Context, codeString *string) (bool, error) {
	return c.repo.Check(ctx, codeString)
}

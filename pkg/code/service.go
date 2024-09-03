package code

import "context"

// CodeService : represents the service of code
type CodeService interface {
	CreateCode(ctx context.Context, code *CodeEntity) error
	UpdateCode(ctx context.Context, codeString *string) error
	CheckCode(ctx context.Context, codeString *string) (bool, error)
}

// codeService : represents the service of code
type codeService struct {
	repo CodeRepository
}

// NewCodeService : creates a new code service
func NewCodeService(repo CodeRepository) CodeService {
	return &codeService{repo: repo}
}

// CreateCode : creates a new code
func (c *codeService) CreateCode(ctx context.Context, code *CodeEntity) error {
	return c.repo.Create(ctx, code)
}

// UpdateCode : updates a code
func (c *codeService) UpdateCode(ctx context.Context, codeString *string) error {
	return c.repo.Update(ctx, codeString)
}

// CheckCode : checks if a code exists
func (c *codeService) CheckCode(ctx context.Context, codeString *string) (bool, error) {
	return c.repo.Check(ctx, codeString)
}

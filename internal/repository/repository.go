package repository

import (
	"context"

	"foodsharing-backend/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Repositories struct {
	Users          Users
	Groups         Groups
	Sessions       Sessions
	DonorCompanies DonorCompanies
	Acts           Acts
	ActContents    ActContents
	Files          Files
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Users:          NewUsersRepository(pool),
		Groups:         NewGroupsRepository(pool),
		Sessions:       NewSessionsRepository(pool),
		DonorCompanies: NewDonorCompaniesRepository(pool),
		Acts:           NewActsRepository(pool),
		ActContents:    NewActContentsRepository(pool),
		Files:          NewFilesRepository(pool),
	}
}

type Users interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id domain.ID) error

	GetByID(ctx context.Context, id domain.ID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
}

type Groups interface {
	Create(ctx context.Context, group *domain.Group) error
	Update(ctx context.Context, group domain.Group) error
	Delete(ctx context.Context, id domain.ID) error

	GetByID(ctx context.Context, id domain.ID) (domain.Group, error)
	GetByName(ctx context.Context, name string) ([]domain.Group, error)
	GetByPermissions(ctx context.Context, permission domain.Permission) ([]domain.Group, error)
	GetAll(ctx context.Context) ([]domain.Group, error)

	AddUser(ctx context.Context, groupID domain.ID, userID domain.ID) error
	GetUserGroups(ctx context.Context, userID domain.ID) ([]domain.Group, error)
	RemoveUser(ctx context.Context, groupID domain.ID, userID domain.ID) error
}

type Sessions interface {
	Create(ctx context.Context, session domain.Session) error

	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error)
	GetByUserID(ctx context.Context, userID domain.ID) ([]domain.Session, error)
}

type DonorCompanies interface {
	Create(ctx context.Context, company *domain.DonorCompany) error
	Update(ctx context.Context, company domain.DonorCompany) error
	Delete(ctx context.Context, id domain.ID) error

	GetByID(ctx context.Context, id domain.ID) (domain.DonorCompany, error)
	GetByCity(ctx context.Context, cityID domain.ID) ([]domain.DonorCompany, error)
	GetAll(ctx context.Context) ([]domain.DonorCompany, error)
}

type Acts interface {
	Create(ctx context.Context, act *domain.Act) error
	Update(ctx context.Context, act domain.Act) error
	Delete(ctx context.Context, id domain.ID) error

	GetByID(ctx context.Context, id domain.ID) (domain.Act, error)
	GetByUserID(ctx context.Context, userID domain.ID) ([]domain.Act, error)
	GetByDonorCompanyID(ctx context.Context, donorCompanyID domain.ID) ([]domain.Act, error)
	GetAll(ctx context.Context) ([]domain.Act, error)

	AddFile(ctx context.Context, fileID domain.ID, actID domain.ID) error
	GetActFiles(ctx context.Context, actID domain.ID) ([]domain.File, error)
	RemoveFile(ctx context.Context, fileID domain.ID, actID domain.ID) error
}

type ActContents interface {
	Create(ctx context.Context, contents ...domain.ActContent) error
	Update(ctx context.Context, contents ...domain.ActContent) error
	Delete(ctx context.Context, contentIDs ...domain.ID) error

	GetByID(ctx context.Context, id domain.ID) (domain.ActContent, error)
	GetByActID(ctx context.Context, actID domain.ID) ([]domain.ActContent, error)
}

type Files interface {
	Create(ctx context.Context, file *domain.File) error
	UpdateStatus(ctx context.Context, fileID domain.ID, status domain.FileStatus) error
	GetForUploading(ctx context.Context) (domain.File, error)
	UpdateStatusAndSetURL(ctx context.Context, fileID domain.ID, url string) error
	GetByID(ctx context.Context, fileID domain.ID) (domain.File, error)
}

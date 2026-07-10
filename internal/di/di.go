package di

import (
	"context"
	"net/http"
	"os"

	"github.com/hemra-siirow/literary/internal/domain/repository"
	infraauth "github.com/hemra-siirow/literary/internal/infrastructure/auth"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
	pg "github.com/hemra-siirow/literary/internal/infrastructure/persistence/postgres"
	"github.com/hemra-siirow/literary/internal/infrastructure/storage"
	"github.com/hemra-siirow/literary/internal/presentation/http/handler"
	prouter "github.com/hemra-siirow/literary/internal/presentation/http/router"
	authusecase "github.com/hemra-siirow/literary/internal/usecase/auth"
	"github.com/hemra-siirow/literary/internal/usecase/biography"
	"github.com/hemra-siirow/literary/internal/usecase/book"
	"github.com/hemra-siirow/literary/internal/usecase/broadcast"
	"github.com/hemra-siirow/literary/internal/usecase/category"
	"github.com/hemra-siirow/literary/internal/usecase/film"
	"github.com/hemra-siirow/literary/internal/usecase/personalletter"
	"github.com/hemra-siirow/literary/internal/usecase/photoarchive"
	"github.com/hemra-siirow/literary/internal/usecase/theatre"
	"github.com/hemra-siirow/literary/internal/usecase/translation"
	"github.com/hemra-siirow/literary/internal/usecase/work"
)

type Container struct {
	Router   http.Handler
	Shutdown func(context.Context) error
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// ensure upload base path exists
	if err := os.MkdirAll(cfg.UploadBasePath, 0o755); err != nil {
		return nil, err
	}

	pool, err := pg.NewPool(cfg)
	if err != nil {
		return nil, err
	}

	// adapters
	jwt := infraauth.NewJWTAdapter(cfg)
	hasher := infraauth.NewBcryptHasher(12)

	// repositories
	workRepo := pg.NewWorkRepository(pool)
	userRepo := pg.NewUserRepository(pool)
	refreshTokenRepo := pg.NewRefreshTokenRepository(pool)
	bookRepo := pg.NewBookRepository(pool)
	broadcastRepo := pg.NewBroadcastRepository(pool)
	filmRepo := pg.NewFilmRepository(pool)
	photoRepo := pg.NewPhotoArchiveRepository(pool)
	plRepo := pg.NewPersonalLetterRepository(pool)
	translatedByRepo := pg.NewTranslatedByAuthorRepository(pool)
	translatedIntoRepo := pg.NewTranslatedIntoLanguageRepository(pool)
	biographyRepo := pg.NewBiographyRepository(pool)
	biographyEventRepo := pg.NewBiographyEventRepository(pool)
	theatreRepo := pg.NewTheatreProductionRepository(pool)
	categoryRepo := pg.NewCategoryRepository(pool)
	localStorage := storage.NewLocalStorage(cfg.UploadBasePath)
	bookPhotoRepo := pg.NewBookPhotoRepository(pool)

	// services
	workSvc := work.NewWorkService(workRepo, localStorage)
	bookSvc := book.NewBookService(bookRepo, bookPhotoRepo, localStorage)
	broadcastSvc := broadcast.NewBroadcastService(broadcastRepo, localStorage)
	filmSvc := film.NewFilmService(filmRepo, localStorage)
	photoSvc := photoarchive.NewPhotoArchiveService(photoRepo, localStorage)
	plSvc := personalletter.NewPersonalLetterService(plRepo, localStorage)
	translationSvc := translation.NewTranslationService(translatedByRepo, translatedIntoRepo)
	bioSvc := biography.NewBiographyService(biographyRepo, localStorage)
	biographyEventSvc := biography.NewBiographyEventService(biographyEventRepo, biographyRepo)
	categorySvc := category.NewCategoryService(categoryRepo)
	theatreSvc := theatre.NewTheatreService(theatreRepo)

	// auth use cases
	loginUseCase := authusecase.NewLoginUseCase(userRepo, hasher, jwt, refreshTokenRepo)
	refreshUseCase := authusecase.NewRefreshTokenUseCase(userRepo, jwt, refreshTokenRepo)
	logoutUseCase := authusecase.NewLogoutUseCase(refreshTokenRepo)
	createUserUseCase := authusecase.NewCreateUserUseCase(userRepo, hasher, jwt, refreshTokenRepo)
	changePasswordUseCase := authusecase.NewChangePasswordUseCase(userRepo, hasher)

	// handlers
	authHandler := handler.NewAuthHandler(loginUseCase, refreshUseCase, logoutUseCase, createUserUseCase, changePasswordUseCase)
	workHandler := handler.NewWorkHandler(workSvc, localStorage)
	bookHandler := handler.NewBookHandler(bookSvc, localStorage)
	broadcastHandler := handler.NewBroadcastHandler(broadcastSvc, localStorage)
	filmHandler := handler.NewFilmHandler(filmSvc, localStorage)
	photoHandler := handler.NewPhotoArchiveHandler(photoSvc, localStorage)
	plHandler := handler.NewPersonalLetterHandler(plSvc, localStorage)
	translationHandler := handler.NewTranslationHandler(translationSvc)
	bioHandler := handler.NewBiographyHandler(bioSvc, localStorage)
	biographyEventHandler := handler.NewBiographyEventHandler(biographyEventSvc, bioSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	theatreHandler := handler.NewTheatreHandler(theatreSvc)
	registrars := []prouter.RouteRegistrar{authHandler, workHandler, bookHandler, broadcastHandler, filmHandler, photoHandler, plHandler, translationHandler, bioHandler, biographyEventHandler, categoryHandler, theatreHandler}

	// small provider to satisfy router's tokenGen param
	tg := &jwtProvider{j: jwt}

	r := prouter.NewRouter(registrars, "*", tg)

	shutdown := func(ctx context.Context) error {
		pool.Close()
		return nil
	}

	return &Container{Router: r, Shutdown: shutdown}, nil
}

type jwtProvider struct {
	j *infraauth.JWTAdapter
}

func (p *jwtProvider) TokenGenerator() repository.TokenGenerator {
	return p.j
}

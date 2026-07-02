package di

import (
	"context"
	"net/http"
	"os"

	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/hemra-siirow/literary/internal/infrastructure/auth"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
	pg "github.com/hemra-siirow/literary/internal/infrastructure/persistence/postgres"
	"github.com/hemra-siirow/literary/internal/infrastructure/storage"
	"github.com/hemra-siirow/literary/internal/presentation/http/handler"
	prouter "github.com/hemra-siirow/literary/internal/presentation/http/router"
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
	jwt := auth.NewJWTAdapter(cfg)

	// repositories
	workRepo := pg.NewWorkRepository(pool)
	bookRepo := pg.NewBookRepository(pool)
	broadcastRepo := pg.NewBroadcastRepository(pool)
	filmRepo := pg.NewFilmRepository(pool)
	photoRepo := pg.NewPhotoArchiveRepository(pool)
	plRepo := pg.NewPersonalLetterRepository(pool)
	translatedByRepo := pg.NewTranslatedByAuthorRepository(pool)
	translatedIntoRepo := pg.NewTranslatedIntoLanguageRepository(pool)
	biographyRepo := pg.NewBiographyRepository(pool)
	theatreRepo := pg.NewTheatreProductionRepository(pool)
	categoryRepo := pg.NewCategoryRepository(pool)
	localStorage := storage.NewLocalStorage(cfg.UploadBasePath)

	bookPhotoRepo := pg.NewBookPhotoRepository(pool)

	// usecases for works
	createWorkUC := work.NewCreateWorkUseCase(workRepo)
	getWorkUC := work.NewGetWorkUseCase(workRepo)
	listWorkUC := work.NewListWorksUseCase(workRepo)
	updateWorkUC := work.NewUpdateWorkUseCase(workRepo, localStorage)
	deleteWorkUC := work.NewDeleteWorkUseCase(workRepo, localStorage)

	// usecases for books
	createBookUC := book.NewCreateBookUseCase(bookRepo)
	getBookUC := book.NewGetBookUseCase(bookRepo)
	listBookUC := book.NewListBooksUseCase(bookRepo)
	updateBookUC := book.NewUpdateBookUseCase(bookRepo, localStorage)
	deleteBookUC := book.NewDeleteBookUseCase(bookRepo, bookPhotoRepo, localStorage)

	// broadcasts
	createBroadcastUC := broadcast.NewCreateBroadcastUseCase(broadcastRepo)
	getBroadcastUC := broadcast.NewGetBroadcastUseCase(broadcastRepo)
	listBroadcastUC := broadcast.NewListBroadcastsUseCase(broadcastRepo)
	updateBroadcastUC := broadcast.NewUpdateBroadcastUseCase(broadcastRepo)
	deleteBroadcastUC := broadcast.NewDeleteBroadcastUseCase(broadcastRepo, localStorage)

	// films
	createFilmUC := film.NewCreateFilmUseCase(filmRepo)
	getFilmUC := film.NewGetFilmUseCase(filmRepo)
	listFilmUC := film.NewListFilmsUseCase(filmRepo)
	updateFilmUC := film.NewUpdateFilmUseCase(filmRepo)
	deleteFilmUC := film.NewDeleteFilmUseCase(filmRepo)

	// photo archive
	createPhotoUC := photoarchive.NewCreatePhotoArchiveUseCase(photoRepo)
	getPhotoUC := photoarchive.NewGetPhotoArchiveUseCase(photoRepo)
	listPhotoUC := photoarchive.NewListPhotoArchiveUseCase(photoRepo)
	updatePhotoUC := photoarchive.NewUpdatePhotoArchiveUseCase(photoRepo)
	deletePhotoUC := photoarchive.NewDeletePhotoArchiveUseCase(photoRepo)

	// personal letters
	createPLUC := personalletter.NewCreatePersonalLetterUseCase(plRepo)
	getPLUC := personalletter.NewGetPersonalLetterUseCase(plRepo)
	listPLUC := personalletter.NewListPersonalLettersUseCase(plRepo)
	updatePLUC := personalletter.NewUpdatePersonalLetterUseCase(plRepo)
	deletePLUC := personalletter.NewDeletePersonalLetterUseCase(plRepo)

	// translations
	createByAuthorUC := translation.NewCreateTranslatedByAuthorUseCase(translatedByRepo)
	getByAuthorUC := translation.NewGetTranslatedByAuthorUseCase(translatedByRepo)
	listByAuthorUC := translation.NewListTranslatedByAuthorsUseCase(translatedByRepo)
	deleteByAuthorUC := translation.NewDeleteTranslatedByAuthorUseCase(translatedByRepo)

	createIntoUC := translation.NewCreateTranslatedIntoLanguageUseCase(translatedIntoRepo)
	getIntoUC := translation.NewGetTranslatedIntoLanguageUseCase(translatedIntoRepo)
	listIntoUC := translation.NewListTranslatedIntoLanguagesUseCase(translatedIntoRepo)
	deleteIntoUC := translation.NewDeleteTranslatedIntoLanguageUseCase(translatedIntoRepo)

	// biography
	createBiographyUC := biography.NewCreateBiographyUseCase(biographyRepo)
	getBiographyUC := biography.NewGetBiographyUseCase(biographyRepo)
	updateBiographyUC := biography.NewUpdateBiographyUseCase(biographyRepo)

	// categories
	createCategoryUC := category.NewCreateCategoryUseCase(categoryRepo)
	getCategoryUC := category.NewGetCategoryUseCase(categoryRepo)
	listCategoriesUC := category.NewListCategoriesUseCase(categoryRepo)
	updateCategoryUC := category.NewUpdateCategoryUseCase(categoryRepo)
	deleteCategoryUC := category.NewDeleteCategoryUseCase(categoryRepo)

	// theatre
	createTheatreUC := theatre.NewCreateTheatreProductionUseCase(theatreRepo)
	getTheatreUC := theatre.NewGetTheatreProductionUseCase(theatreRepo)
	listTheatreUC := theatre.NewListTheatreProductionsUseCase(theatreRepo)
	updateTheatreUC := theatre.NewUpdateTheatreProductionUseCase(theatreRepo)
	deleteTheatreUC := theatre.NewDeleteTheatreProductionUseCase(theatreRepo)

	// handlers
	workHandler := handler.NewWorkHandler(createWorkUC, getWorkUC, listWorkUC, updateWorkUC, deleteWorkUC, localStorage)
	bookHandler := handler.NewBookHandler(createBookUC, getBookUC, listBookUC, updateBookUC, deleteBookUC, localStorage)
	broadcastHandler := handler.NewBroadcastHandler(createBroadcastUC, getBroadcastUC, listBroadcastUC, updateBroadcastUC, deleteBroadcastUC, localStorage)
	filmHandler := handler.NewFilmHandler(createFilmUC, getFilmUC, listFilmUC, updateFilmUC, deleteFilmUC, localStorage)
	photoHandler := handler.NewPhotoArchiveHandler(createPhotoUC, getPhotoUC, listPhotoUC, updatePhotoUC, deletePhotoUC, localStorage)
	plHandler := handler.NewPersonalLetterHandler(createPLUC, getPLUC, listPLUC, updatePLUC, deletePLUC, localStorage)
	translationHandler := handler.NewTranslationHandler(createByAuthorUC, getByAuthorUC, listByAuthorUC, deleteByAuthorUC, createIntoUC, getIntoUC, listIntoUC, deleteIntoUC)
	bioHandler := handler.NewBiographyHandler(getBiographyUC, createBiographyUC, updateBiographyUC, localStorage)
	categoryHandler := handler.NewCategoryHandler(createCategoryUC, getCategoryUC, listCategoriesUC, updateCategoryUC, deleteCategoryUC)
	theatreHandler := handler.NewTheatreHandler(createTheatreUC, getTheatreUC, listTheatreUC, updateTheatreUC, deleteTheatreUC)

	registrars := []prouter.RouteRegistrar{workHandler, bookHandler, broadcastHandler, filmHandler, photoHandler, plHandler, translationHandler, bioHandler, categoryHandler, theatreHandler}

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
	j *auth.JWTAdapter
}

func (p *jwtProvider) TokenGenerator() repository.TokenGenerator {
	return p.j
}

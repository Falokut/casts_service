package postgresrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Falokut/casts_service/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CastsRepository struct {
	logger *logrus.Logger
	db     *sqlx.DB
}

const (
	castTableName        = "casts"
	professionsTableName = "professions"
)

func NewCastsRepository(db *sqlx.DB, logger *logrus.Logger) *CastsRepository {
	return &CastsRepository{db: db, logger: logger}
}

func (r *CastsRepository) Shutdown() {
	r.db.Close()
}

func (r *CastsRepository) GetCast(ctx context.Context, id int32, professionsIds []int32) (cast models.Cast, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetCast")

	query := fmt.Sprintf("SELECT person_id, profession_id, COALESCE(%[1]s.name,'') AS profession_name FROM %[2]s "+
		"LEFT JOIN %[1]s ON profession_id=%[1]s.id WHERE movie_id=$1", professionsTableName, castTableName)

	var persons []models.Person
	if len(professionsIds) > 0 {
		query += " AND profession_id=ANY($2)"
		err = r.db.SelectContext(ctx, &persons, query, id, professionsIds)
	} else {
		err = r.db.SelectContext(ctx, &persons, query, id)
	}

	if err != nil {
		return
	}

	cast.Persons = persons
	return
}

func (r *CastsRepository) GetProfessions(ctx context.Context) (professions []models.Profession, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetProfessions")

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", professionsTableName)
	err = r.db.SelectContext(ctx, &professions, query)
	return
}

func (r *CastsRepository) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("movies repository error occurred")
	} else {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("movies repository error occurred")
	}
}

package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type postgreRepository struct {
	db *sqlx.DB
}

const (
	castTableName        = "casts"
	professionsTableName = "professions"
)

func NewCastsRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *postgreRepository) Shutdown() {
	r.db.Close()
}

func (r *postgreRepository) GetCast(ctx context.Context, id int32, professionsIds []int32) (Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetCast")
	defer span.Finish()

	query := fmt.Sprintf("SELECT actor_id, profession_id, COALESCE(%[1]s.name,'') AS profession_name FROM %[2]s "+
		"LEFT JOIN %[1]s ON profession_id=%[1]s.id WHERE movie_id=$1", professionsTableName, castTableName)

	var actors []Actor
	var err error
	if len(professionsIds) > 0 {
		query += " AND profession_id=ANY($2)"
		err = r.db.SelectContext(ctx, &actors, query, id, professionsIds)
	} else {
		err = r.db.SelectContext(ctx, &actors, query, id)
	}

	if err != nil {
		return Cast{}, err
	}
	return Cast{Actors: actors}, nil
}

func (r *postgreRepository) GetProfessions(ctx context.Context) ([]Profession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetCast")
	defer span.Finish()

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", professionsTableName)
	var professions []Profession
	err := r.db.SelectContext(ctx, &professions, query)
	if err != nil {
		return []Profession{}, err
	}
	return professions, nil
}

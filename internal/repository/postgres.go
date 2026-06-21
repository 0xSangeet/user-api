package repository

import (
	"context"
	"errors"

	"github.com/0xSangeet/user-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostGresRepo struct {
	db *pgxpool.Pool
}

func NewPostGresRepo(db *pgxpool.Pool) *PostGresRepo {
	return &PostGresRepo{
		db: db,
	}
}

func (p *PostGresRepo) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := p.db.Exec(ctx, query, u.ID, u.Name, u.Email)

	if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            if pgErr.Code == "23505" {
                return domain.ErrUserAlreadyExists
            }
        }
        return err
    }

    return nil
}

func (p *PostGresRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
    query := `SELECT id, name, email FROM users WHERE id = $1`

    var user domain.User

    err := p.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound
        }

        return nil, err
    }

    return &user, nil
}

func (p *PostGresRepo) GetAll(ctx context.Context) ([]domain.User, error) {

	query := `SELECT id, name, email FROM users`

	rows, err := p.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]domain.User, 0)

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		results = append(results, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *PostGresRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	tag, err := p.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
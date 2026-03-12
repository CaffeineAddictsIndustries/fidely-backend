package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"fidely-backend/internal/auth"
	"fidely-backend/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateAdminUsername = errors.New("duplicate admin username across admin tables")

// AdminAuthRepository persists admin auth state in PostgreSQL.
type AdminAuthRepository struct {
	pool *pgxpool.Pool
}

func NewAdminAuthRepository(pool *pgxpool.Pool) *AdminAuthRepository {
	return &AdminAuthRepository{pool: pool}
}

func (repo *AdminAuthRepository) FindByUsername(ctx context.Context, username string) (*service.AdminPrincipal, error) {
	storeAdmin, err := repo.findStoreAdminByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	fidelyAdmin, err := repo.findFidelyAdminByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if storeAdmin != nil && fidelyAdmin != nil {
		return nil, ErrDuplicateAdminUsername
	}
	if storeAdmin != nil {
		return storeAdmin, nil
	}
	if fidelyAdmin != nil {
		return fidelyAdmin, nil
	}

	return nil, nil
}

func (repo *AdminAuthRepository) FindByTypeAndID(ctx context.Context, adminType auth.AdminType, adminID int) (*service.AdminPrincipal, error) {
	switch adminType {
	case auth.AdminTypeStoreAdmin:
		return repo.findStoreAdminByID(ctx, adminID)
	case auth.AdminTypeFidelyAdmin:
		return repo.findFidelyAdminByID(ctx, adminID)
	default:
		return nil, auth.ErrInvalidAdminType
	}
}

func (repo *AdminAuthRepository) CreateSession(ctx context.Context, session auth.AdminSession) (auth.AdminSession, error) {
	query := `
		INSERT INTO admin_sessions (
			admin_type,
			admin_id,
			session_token_hash,
			expires_at,
			revoked_at,
			created_at,
			last_seen_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	persisted := session
	if err := repo.pool.QueryRow(
		ctx,
		query,
		string(session.AdminType),
		session.AdminID,
		session.SessionTokenHash,
		session.ExpiresAt,
		session.RevokedAt,
		session.CreatedAt,
		session.LastSeenAt,
	).Scan(&persisted.ID); err != nil {
		return auth.AdminSession{}, fmt.Errorf("insert admin session: %w", err)
	}

	return persisted, nil
}

func (repo *AdminAuthRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*auth.AdminSession, error) {
	query := `
		SELECT id, admin_type, admin_id, session_token_hash, expires_at, revoked_at, created_at, last_seen_at
		FROM admin_sessions
		WHERE session_token_hash = $1
		LIMIT 1
	`

	var adminType string
	var session auth.AdminSession
	if err := repo.pool.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&adminType,
		&session.AdminID,
		&session.SessionTokenHash,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
		&session.LastSeenAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query admin session: %w", err)
	}

	session.AdminType = auth.AdminType(adminType)
	return &session, nil
}

func (repo *AdminAuthRepository) UpdateSession(ctx context.Context, session auth.AdminSession) error {
	query := `
		UPDATE admin_sessions
		SET session_token_hash = $2,
			expires_at = $3,
			revoked_at = $4,
			last_seen_at = $5
		WHERE id = $1
	`

	commandTag, err := repo.pool.Exec(
		ctx,
		query,
		session.ID,
		session.SessionTokenHash,
		session.ExpiresAt,
		session.RevokedAt,
		session.LastSeenAt,
	)
	if err != nil {
		return fmt.Errorf("update admin session: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return auth.ErrInvalidSession
	}

	return nil
}

func (repo *AdminAuthRepository) findStoreAdminByUsername(ctx context.Context, username string) (*service.AdminPrincipal, error) {
	query := `
		SELECT id, username, password_hash, role, store_id
		FROM store_admin
		WHERE username = $1
		LIMIT 1
	`

	var principal service.AdminPrincipal
	var storeID sql.NullInt32
	if err := repo.pool.QueryRow(ctx, query, username).Scan(
		&principal.AdminID,
		&principal.Username,
		&principal.PasswordHash,
		&principal.Role,
		&storeID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query store_admin by username: %w", err)
	}

	principal.AdminType = auth.AdminTypeStoreAdmin
	if storeID.Valid {
		value := int(storeID.Int32)
		principal.StoreID = &value
	}
	return &principal, nil
}

func (repo *AdminAuthRepository) findFidelyAdminByUsername(ctx context.Context, username string) (*service.AdminPrincipal, error) {
	query := `
		SELECT id, username, password_hash, role
		FROM fidely_admin
		WHERE username = $1
		LIMIT 1
	`

	var principal service.AdminPrincipal
	var role int
	if err := repo.pool.QueryRow(ctx, query, username).Scan(
		&principal.AdminID,
		&principal.Username,
		&principal.PasswordHash,
		&role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query fidely_admin by username: %w", err)
	}

	principal.AdminType = auth.AdminTypeFidelyAdmin
	principal.Role = strconv.Itoa(role)
	return &principal, nil
}

func (repo *AdminAuthRepository) findStoreAdminByID(ctx context.Context, adminID int) (*service.AdminPrincipal, error) {
	query := `
		SELECT id, username, password_hash, role, store_id
		FROM store_admin
		WHERE id = $1
		LIMIT 1
	`

	var principal service.AdminPrincipal
	var storeID sql.NullInt32
	if err := repo.pool.QueryRow(ctx, query, adminID).Scan(
		&principal.AdminID,
		&principal.Username,
		&principal.PasswordHash,
		&principal.Role,
		&storeID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query store_admin by id: %w", err)
	}

	principal.AdminType = auth.AdminTypeStoreAdmin
	if storeID.Valid {
		value := int(storeID.Int32)
		principal.StoreID = &value
	}
	return &principal, nil
}

func (repo *AdminAuthRepository) findFidelyAdminByID(ctx context.Context, adminID int) (*service.AdminPrincipal, error) {
	query := `
		SELECT id, username, password_hash, role
		FROM fidely_admin
		WHERE id = $1
		LIMIT 1
	`

	var principal service.AdminPrincipal
	var role int
	if err := repo.pool.QueryRow(ctx, query, adminID).Scan(
		&principal.AdminID,
		&principal.Username,
		&principal.PasswordHash,
		&role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query fidely_admin by id: %w", err)
	}

	principal.AdminType = auth.AdminTypeFidelyAdmin
	principal.Role = strconv.Itoa(role)
	return &principal, nil
}

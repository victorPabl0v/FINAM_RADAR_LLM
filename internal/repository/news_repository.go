package repository

import (
	"FINAM_RADAR/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	INewsRepository
	pool *pgxpool.Pool
}

func NewNewsRepository(pool *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{pool: pool}
}

func (repo *NewsRepository) Add(news models.NewsRequest, ctx context.Context) (AddNews, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return AddNews{Status: false}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, n := range news.News {
		var dedupGroupID int64
		err = tx.QueryRow(ctx,
			`
			INSERT INTO news.dedup_groups (dedup_group)
			VALUES ($1)
			ON  CONFLICT (dedup_group) DO UPDATE SET dedup_group = EXCLUDED.dedup_group
			RETURNING id;
		`, n.DedupGroup).Scan(&dedupGroupID)
		if err != nil {
			return AddNews{Status: false}, err
		}

		var draftID int64
		err = tx.QueryRow(ctx,
			`
			INSERT INTO news.draft (title, lead, quote)
			VALUES ($1, $2, $3)
			RETURNING id;
		`, n.Draft.Title, n.Draft.Lead, n.Draft.Quote).Scan(&draftID)
		if err != nil {
			return AddNews{Status: false}, err
		}

		for _, bullet := range n.Draft.Bullets {
			var bulletID int64
			err = tx.QueryRow(ctx, `
			INSERT INTO bullets.bullets (bullet)
			VALUES ($1)
			ON CONFLICT (bullet) DO UPDATE SET bullet = EXCLUDED.bullet
			RETURNING id;
		`, bullet).Scan(&bulletID)
			if err != nil {
				return AddNews{Status: false}, err
			}

			_, err = tx.Exec(ctx, `
			INSERT INTO bullets.draft_bullets (draft_id, bullet_id)
			VALUES ($1, $2)
			ON CONFLICT (draft_id, bullet_id) DO NOTHING;
		`, draftID, bulletID)
			if err != nil {
				return AddNews{Status: false}, err
			}
		}
		var newsID int64
		err = tx.QueryRow(ctx, `
		INSERT INTO news.news (headline, hotness, why_now, timeline, draft_id, dedup_group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`, n.Headline, n.Hotness, n.WhyNow, n.Timeline, draftID, dedupGroupID).Scan(&newsID)
		if err != nil {
			return AddNews{Status: false}, err
		}
		for _, entity := range n.Entities {
			var entityID int64
			err = tx.QueryRow(ctx, `
			INSERT INTO entities.entities (entity)
			VALUES ($1)
			ON CONFLICT (entity) DO UPDATE SET entity = EXCLUDED.entity
			RETURNING id;
		`, entity).Scan(&entityID)
			if err != nil {
				return AddNews{Status: false}, err
			}

			_, err = tx.Exec(ctx, `
			INSERT INTO entities.news_entities (news_id, entity_id)
			VALUES ($1, $2)
			ON CONFLICT (news_id, entity_id) DO NOTHING;
		`, newsID, entityID)
			if err != nil {
				return AddNews{Status: false}, err
			}
		}
		for _, source := range n.Sources {
			var sourceID int64
			err = tx.QueryRow(ctx, `
			INSERT INTO sources.sources (source)
			VALUES ($1)
			ON CONFLICT (source) DO UPDATE SET source = EXCLUDED.source
			RETURNING id;
		`, source).Scan(&sourceID)
			if err != nil {
				return AddNews{Status: false}, err
			}

			_, err = tx.Exec(ctx, `
			INSERT INTO sources.news_sources (news_id, source_id)
			VALUES ($1, $2)
			ON CONFLICT (news_id, source_id) DO NOTHING;
		`, newsID, sourceID)
			if err != nil {
				return AddNews{Status: false}, err
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return AddNews{Status: false}, err
	}
	return AddNews{Status: true}, err
}

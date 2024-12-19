package source

import (
	"context"
	"strings"

	"github.com/SlyMarbo/rss"
	"github.com/didsqq/news_feed_bot/internal/models"
	"github.com/samber/lo"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func NewRSSSourceFromModel(m models.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := s.loadFeed(ctx, s.URL)
	if err != nil {
		return nil, err
	}

	return lo.Map(feed.Items, func(item *rss.Item, _ int) models.Item {
		return models.Item{
			Title:      item.Title,
			Categories: item.Categories,
			Link:       item.Link,
			Date:       item.Date,
			SourceName: s.SourceName,
			Summary:    strings.TrimSpace(item.Summary),
		}
	}), nil
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}
		feedCh <- feed
	}()

	select { // реагирует на первый полученный сигнал
	case <-ctx.Done(): // контекст завершился(отменен или закончилось в время), функция возвращает ошибку контекса
		return nil, ctx.Err()
	case err := <-errCh: // ошибка в горутине
		return nil, err
	case feed := <-feedCh: // успех
		return feed, nil
	}
}
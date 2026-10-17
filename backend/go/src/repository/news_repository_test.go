package repository

import (
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetNews_UsesDayBoundaryOfTargetLocation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.News{}, &models.NewsReaction{}); err != nil {
		t.Fatal(err)
	}

	tokyo := mustLoadLocation(t, "Asia/Tokyo")
	newYork := mustLoadLocation(t, "America/New_York")
	fixtures := map[string]time.Time{
		"tokyo-1013-2359": time.Date(2026, 10, 13, 23, 59, 59, 999999999, tokyo),
		"tokyo-1014-0000": time.Date(2026, 10, 14, 0, 0, 0, 0, tokyo),
		"utc-1013-1810":   time.Date(2026, 10, 13, 18, 10, 0, 123456789, time.UTC),
		"tokyo-1014-2359": time.Date(2026, 10, 14, 23, 59, 0, 0, tokyo),
		"tokyo-1015-0000": time.Date(2026, 10, 15, 0, 0, 0, 0, tokyo),
		"ny-1101-0000":    time.Date(2026, 11, 1, 0, 0, 0, 0, newYork),
		"ny-1101-2330":    time.Date(2026, 11, 1, 23, 30, 0, 0, newYork),
		"ny-1102-0000":    time.Date(2026, 11, 2, 0, 0, 0, 0, newYork),
	}
	for title, createdAt := range fixtures {
		if err := db.Create(&models.News{Title: title, CreatedAt: createdAt}).Error; err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		target time.Time
		want   []string
	}{
		{
			name:   "Asia/Tokyo",
			target: time.Date(2026, 10, 14, 0, 0, 0, 0, tokyo),
			want:   []string{"tokyo-1014-0000", "tokyo-1014-2359", "utc-1013-1810"},
		},
		{
			name:   "UTC",
			target: time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC),
			want:   []string{"tokyo-1013-2359", "tokyo-1014-0000", "utc-1013-1810"},
		},
		{
			name:   "America/New_York on DST end day",
			target: time.Date(2026, 11, 1, 0, 0, 0, 0, newYork),
			want:   []string{"ny-1101-0000", "ny-1101-2330"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			news := NewNewsRepository(db).GetNews(tt.target, c)

			got := make([]string, len(news))
			for i, n := range news {
				got[i] = n.Title
			}
			sort.Strings(got)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func TestGetLatestNews_OrdersByCreatedAtDescAndLimits(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.News{}, &models.NewsReaction{}); err != nil {
		t.Fatal(err)
	}

	tokyo := mustLoadLocation(t, "Asia/Tokyo")
	fixtures := map[string]time.Time{
		"oldest":       time.Date(2026, 10, 15, 3, 0, 0, 0, time.UTC),
		"middle-tokyo": time.Date(2026, 10, 16, 9, 0, 0, 0, tokyo),
		"newest":       time.Date(2026, 10, 16, 1, 0, 0, 0, time.UTC),
	}
	for title, createdAt := range fixtures {
		if err := db.Create(&models.News{Title: title, CreatedAt: createdAt}).Error; err != nil {
			t.Fatal(err)
		}
	}

	news, err := NewNewsRepository(db).GetLatestNews(2)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"newest", "middle-tokyo"}
	if len(news) != len(want) {
		t.Fatalf("expected %v, got %d items", want, len(news))
	}
	for i := range want {
		if news[i].Title != want[i] {
			t.Fatalf("expected %v at %d, got %s", want[i], i, news[i].Title)
		}
	}
}

func mustLoadLocation(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

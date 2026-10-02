package main

import (
	"fmt"
	"log"
	"time"

	"github.com/birdseyeapi/birds_eye_v3/go/src/ai"
	"github.com/birdseyeapi/birds_eye_v3/go/src/db"
	"github.com/birdseyeapi/birds_eye_v3/go/src/models"
	"github.com/birdseyeapi/birds_eye_v3/go/src/scraping/news"
	"gorm.io/gorm"
)

const mockURLPrefix = "https://mock.example.com/"

type mockArticle struct {
	title      string
	summary    string
	sourceBy   string
	category   string
	importance float64
	sentiment  models.ReactionSentiment
	comments   []string
}

var mockArticles = []mockArticle{
	{
		title:      "新型AIモデル「Birds-1」が公開、コーディング性能が大幅向上",
		summary:    "モックAI社が新型モデルBirds-1を公開した。コード生成ベンチマークで従来比30%の改善を達成し、長文コンテキストにも対応する。API価格は据え置き。",
		sourceBy:   news.HatenaSourceName,
		category:   ai.CategoryAIModelsIndustry,
		importance: 0.95,
		sentiment:  models.ReactionSentiment{Positive: 5, Neutral: 2, Negative: 1},
		comments:   []string{"性能向上が楽しみ", "価格据え置きは助かる", "ベンチマークだけでは分からない", "試してみる", "コンテキスト長が実用的"},
	},
	{
		title:      "Claude Codeでテスト駆動開発を回す実践ガイド",
		summary:    "コーディングエージェントを使ったTDDの進め方を解説。テストを先に書かせ、失敗確認から実装、リファクタリングまでをエージェントに任せる手順を紹介する。",
		sourceBy:   news.ZennSourceName,
		category:   ai.CategoryAIDevTools,
		importance: 0.82,
		sentiment:  models.ReactionSentiment{Positive: 3, Neutral: 1, Negative: 0},
		comments:   []string{"参考になった", "プロンプト例がありがたい", "うちでも導入したい"},
	},
	{
		title:      "主要ライブラリに深刻な脆弱性、直ちにアップデートを",
		summary:    "広く使われているOSSライブラリでリモートコード実行の脆弱性が見つかった。修正版が公開済みで、利用者には早急なアップデートが呼びかけられている。",
		sourceBy:   news.ZDNetSourceName,
		category:   ai.CategorySecurity,
		importance: 0.9,
		sentiment:  models.ReactionSentiment{Positive: 0, Neutral: 3, Negative: 4},
		comments:   []string{"またか", "今すぐ確認する", "依存関係の管理が大変", "SBOM整備しておくべきだった"},
	},
	{
		title:      "Go 1.30のリリース予定機能まとめ",
		summary:    "次期Goバージョンで追加予定の機能を整理。標準ライブラリの拡充、ジェネリクスの改善、ビルド時間の短縮などが予定されている。",
		sourceBy:   news.HatenaSourceName,
		category:   ai.CategorySoftwareDev,
		importance: 0.6,
		sentiment:  models.ReactionSentiment{Positive: 2, Neutral: 2, Negative: 0},
		comments:   []string{"ビルドが速くなるのは嬉しい", "ジェネリクス周りが気になる", "待ち遠しい", "標準ライブラリ強化は歓迎"},
	},
	{
		title:      "クラウド障害の振り返り：設定変更がリージョン全体に波及",
		summary:    "大手クラウド事業者が先日の障害の事後報告を公開。誤った設定変更が段階的展開の検証をすり抜け、リージョン全体に影響したと説明している。",
		sourceBy:   news.CloudWatchSourceName,
		category:   ai.CategoryCloudInfra,
		importance: 0.7,
		sentiment:  models.ReactionSentiment{Positive: 1, Neutral: 2, Negative: 2},
		comments:   []string{"事後報告が丁寧", "段階展開の重要性を再認識", "マルチリージョンにしておけば", "他人事ではない", "透明性は評価できる"},
	},
	{
		title:      "モック株式会社、第2四半期は増収増益",
		summary:    "モック株式会社が第2四半期決算を発表。クラウド事業が牽引し売上高は前年同期比12%増。通期見通しも上方修正した。",
		sourceBy:   news.ZDNetSourceName,
		category:   ai.CategoryBusiness,
		importance: 0.4,
		sentiment:  models.ReactionSentiment{},
		comments:   nil,
	},
	{
		title:      "ノート向け新プロセッサ発表、消費電力を半減",
		summary:    "新世代プロセッサが発表された。同等性能で消費電力を約半分に抑え、薄型ノートでのバッテリー駆動時間が大きく伸びる見込み。",
		sourceBy:   news.CloudWatchSourceName,
		category:   ai.CategoryHardware,
		importance: 0.5,
		sentiment:  models.ReactionSentiment{Positive: 2, Neutral: 1, Negative: 0},
		comments:   []string{"バッテリー持ちが気になってた", "価格次第", "実機レビュー待ち"},
	},
	{
		title:      "新プロジェクト管理ツール「Nest」が正式版をリリース",
		summary:    "タスク管理とドキュメント機能を統合した新サービスNestが正式版となった。無料プランに加え、チーム向けの有料プランも提供する。",
		sourceBy:   news.ZennSourceName,
		category:   ai.CategoryProductRelease,
		importance: 0.35,
		sentiment:  models.ReactionSentiment{Positive: 1, Neutral: 1, Negative: 0},
		comments:   []string{"UIが良さそう", "既存ツールからの移行が課題"},
	},
	{
		title:      "リモートワーク縮小の動き、出社回帰を進める企業が増加",
		summary:    "国内IT企業で出社頻度を引き上げる動きが広がっている。一方でリモート継続を採用の強みとする企業も多く、働き方の二極化が進んでいる。",
		sourceBy:   news.HatenaSourceName,
		category:   ai.CategorySociety,
		importance: 0.45,
		sentiment:  models.ReactionSentiment{Positive: 0, Neutral: 1, Negative: 5},
		comments:   []string{"通勤が苦痛", "生産性は変わらないのに", "転職を考える", "会社による", "出社の意義が不明", "対面の方が楽という人もいる"},
	},
	{
		title:      "分類しづらい雑多なニュースのサンプル",
		summary:    "どのカテゴリにも当てはまらない記事の表示確認用モック。長めのタイトルや本文でレイアウトが崩れないかを確かめるために、あえて冗長な文章にしている。",
		sourceBy:   news.ZennSourceName,
		category:   ai.CategoryOther,
		importance: 0.1,
		sentiment:  models.ReactionSentiment{},
		comments:   nil,
	},
}

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	removed, err := deleteMockNews(database)
	if err != nil {
		log.Fatalf("Failed to delete previous mock data: %v", err)
	}

	inserted, err := insertMockNews(database, time.Now())
	if err != nil {
		log.Fatalf("Failed to insert mock data: %v", err)
	}

	fmt.Printf("seed completed: removed=%d inserted=%d\n", removed, inserted)
}

func deleteMockNews(database *gorm.DB) (int, error) {
	var ids []uint
	err := database.Unscoped().
		Model(&models.News{}).
		Where("scraped_url LIKE ?", mockURLPrefix+"%").
		Pluck("id", &ids).Error
	if err != nil || len(ids) == 0 {
		return 0, err
	}

	err = database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("news_id IN ?", ids).Delete(&models.NewsReaction{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Where("id IN ?", ids).Delete(&models.News{}).Error
	})
	return len(ids), err
}

func insertMockNews(database *gorm.DB, now time.Time) (int, error) {
	err := database.Transaction(func(tx *gorm.DB) error {
		for i, a := range mockArticles {
			n := buildNews(i, a, now)
			if err := tx.Create(&n).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return len(mockArticles), err
}

func buildNews(i int, a mockArticle, now time.Time) models.News {
	scrapedAt := now.Add(-time.Duration(i) * time.Minute)
	articleURL := fmt.Sprintf("%sarticles/%d", mockURLPrefix, i+1)

	reactions := make([]models.NewsReaction, len(a.comments))
	for j, c := range a.comments {
		reactions[j] = models.NewsReaction{
			Author:          fmt.Sprintf("mock_user_%d", j+1),
			Comment:         c,
			ScrapedDateTime: scrapedAt,
			CommentUrl:      fmt.Sprintf("%s#comment-%d", articleURL, j+1),
		}
	}

	return models.News{
		CreatedAt:          now,
		Title:              a.title,
		Description:        a.summary,
		SummarizedText:     a.summary,
		SourceBy:           a.sourceBy,
		Category:           a.category,
		CategoryConfidence: 0.9,
		Importance:         a.importance,
		ReactionSentiment:  a.sentiment,
		ScrapedUrl:         mockURLPrefix + "source",
		ScrapedDateTime:    scrapedAt,
		ArticleUrl:         articleURL,
		Reactions:          reactions,
	}
}

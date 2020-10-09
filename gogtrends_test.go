package gogtrends

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locUS  = "US"
	catAll = "all"
	langEN = "EN"

	concurrentGoroutinesNum = 10
	loadTestNum             = 20 // changed for test speed up
	catProgramming          = 31
)

func TestDebug(t *testing.T) {
	Debug(true)
	assert.True(t, client.debug)
	Debug(false)
}

func TestDailyTrending(t *testing.T) {
	_, err := Daily(context.Background(), "unknown", "Kashyyyk")
	assert.Error(t, err)

	resp, err := Daily(context.Background(), langEN, locUS)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Searches[0].Title.Query) > 0)
}

func TestRealtimeTrending(t *testing.T) {
	categories := TrendsCategories()
	assert.True(t, len(categories) > 0)
	_, ok := categories[catAll]
	assert.True(t, ok)

	_, err := Realtime(context.Background(), langEN, locUS, "random")
	assert.Error(t, err)

	resp, err := Realtime(context.Background(), langEN, locUS, catAll)
	assert.NoError(t, err)
	assert.True(t, len(resp[0].Title) > 0)
}

func TestRealtimeTrendingConcurrent(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(concurrentGoroutinesNum)
	for i := 0; i < concurrentGoroutinesNum; i++ {
		go func() {
			defer wg.Done()

			categories := TrendsCategories()
			assert.True(t, len(categories) > 0)
			_, ok := categories[catAll]
			assert.True(t, ok)

			resp, err := Realtime(context.Background(), langEN, locUS, catAll)
			assert.NoError(t, err)
			assert.True(t, len(resp[0].Title) > 0)
		}()
	}
	wg.Wait()
}

func TestExploreCategories(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(concurrentGoroutinesNum)
	for i := 0; i < concurrentGoroutinesNum; i++ {
		go func() {
			defer wg.Done()

			exploreCats, err := ExploreCategories(context.Background())
			assert.NoError(t, err)
			assert.True(t, len(exploreCats.Children) > 0)
		}()
	}
}

func TestExploreLocations(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(concurrentGoroutinesNum)
	for i := 0; i < concurrentGoroutinesNum; i++ {
		go func() {
			defer wg.Done()

			exploreLocs, err := ExploreLocations(context.Background())
			assert.NoError(t, err)
			assert.True(t, len(exploreLocs.Children) > 0)
		}()
	}
}

func TestExplore(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)
}

func TestInterestOverTime(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	overTime, err := InterestOverTime(context.Background(), explore[0], langEN)
	assert.NoError(t, err)
	assert.True(t, len(overTime) > 0)

	explore[0].ID = ""
	_, err = InterestOverTime(context.Background(), explore[0], langEN)
	assert.Error(t, err)
}

func TestInterestByLocation(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	byLoc, err := InterestByLocation(context.Background(), explore[1], langEN)
	assert.NoError(t, err)
	assert.True(t, len(byLoc) > 0)

	explore[1].ID = ""
	_, err = InterestByLocation(context.Background(), explore[1], langEN)
	assert.Error(t, err)
}

func TestInterestByLocationConcurrent(t *testing.T) {
	wg := new(sync.WaitGroup)

	wg.Add(concurrentGoroutinesNum)
	for i := 0; i < concurrentGoroutinesNum; i++ {
		go func() {
			defer wg.Done()

			req := &ExploreRequest{
				ComparisonItems: []*ComparisonItem{
					{
						Keyword: "Golang",
						Time:    "today 12-m",
					},
				},
				Category: 31, // Programming category
				Property: "",
			}

			explore, err := Explore(context.Background(), req, langEN)
			assert.NoError(t, err)
			assert.True(t, len(explore) == 4)

			byLoc, err := InterestByLocation(context.Background(), explore[1], langEN)
			assert.NoError(t, err)
			assert.True(t, len(byLoc) > 0)
		}()
	}

	wg.Wait()
}

func TestRelated(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	relatedTopics, err := Related(context.Background(), explore[2], langEN)
	assert.NoError(t, err)
	assert.True(t, len(relatedTopics) > 0)

	relatedQueries, err := Related(context.Background(), explore[3], langEN)
	assert.NoError(t, err)
	assert.True(t, len(relatedQueries) > 0)

	explore[3].ID = ""
	_, err = Related(context.Background(), explore[3], langEN)
	assert.Error(t, err)
}

func TestLoadDaily(t *testing.T) {
	res := make([][]*TrendingSearch, loadTestNum)
	errors := make([]error, loadTestNum)
	for i := 0; i < loadTestNum; i++ {
		res[i], errors[i] = Daily(context.Background(), langEN, locUS)
	}

	for _, e := range errors {
		assert.NoError(t, e)
	}

	for _, r := range res {
		assert.True(t, len(r[0].Title.Query) > 0)
	}
}

func TestLoadRealtime(t *testing.T) {
	res := make([][]*TrendingStory, loadTestNum)
	errors := make([]error, loadTestNum)
	for i := 0; i < loadTestNum; i++ {
		res[i], errors[i] = Realtime(context.Background(), langEN, locUS, catAll)
	}

	for _, e := range errors {
		assert.NoError(t, e)
	}

	for _, r := range res {
		assert.True(t, len(r[0].Title) > 0)
	}
}

func TestLoadOverTime(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	res := make([][]*Timeline, loadTestNum)
	errors := make([]error, loadTestNum)

	for i := 0; i < loadTestNum; i++ {
		res[i], errors[i] = InterestOverTime(context.Background(), explore[0], langEN)
	}

	for _, e := range errors {
		assert.NoError(t, e)
	}

	for _, r := range res {
		assert.True(t, len(r) > 0)
	}
}

func TestLoadByLocation(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)
	assert.True(t, len(explore) == 4)

	res := make([][]*GeoMap, loadTestNum)
	errors := make([]error, loadTestNum)

	for i := 0; i < loadTestNum; i++ {
		res[i], errors[i] = InterestByLocation(context.Background(), explore[1], langEN)
	}

	for _, e := range errors {
		assert.NoError(t, e)
	}

	for _, r := range res {
		assert.True(t, len(r) > 0)
	}
}

func TestCompareInterest(t *testing.T) {
	req := &ExploreRequest{
		ComparisonItems: []*ComparisonItem{
			{
				Keyword: "Golang",
				Geo:     locUS,
				Time:    "today 12-m",
			},
			{
				Keyword: "Python",
				Geo:     locUS,
				Time:    "today+12-m",
			},
			{
				Keyword: "PHP",
				Geo:     locUS,
				Time:    "today+12-m",
			},
			{
				Keyword: "Паскаль",
				Geo:     locUS,
				Time:    "today 12-m",
			},
		},
		Category: catProgramming,
		Property: "",
	}
	explore, err := Explore(context.Background(), req, langEN)
	assert.NoError(t, err)

	// interest over time for 3 compared items in one chart
	overTime, err := InterestOverTime(context.Background(), explore[0], langEN)
	assert.NoError(t, err)
	assert.True(t, len(overTime) > 0)

	// interest over time for 3 compared items in one map
	byLoc, err := InterestByLocation(context.Background(), explore[1], langEN)
	assert.NoError(t, err)
	assert.True(t, len(byLoc) > 0)
}

func TestCompareInterestConcurrent(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(concurrentGoroutinesNum)
	for i := 0; i < concurrentGoroutinesNum; i++ {
		go func() {
			defer wg.Done()

			req := &ExploreRequest{
				ComparisonItems: []*ComparisonItem{
					{
						Keyword: "Golang",
						Geo:     locUS,
						Time:    "today 12-m",
					},
					{
						Keyword: "Python",
						Geo:     locUS,
						Time:    "today+12-m",
					},
					{
						Keyword: "PHP",
						Geo:     locUS,
						Time:    "today+12-m",
					},
					{
						Keyword: "Паскаль",
						Geo:     locUS,
						Time:    "today 12-m",
					},
				},
				Category: catProgramming,
				Property: "",
			}

			explore, err := Explore(context.Background(), req, langEN)
			assert.NoError(t, err)

			// interest over time for 3 compared items in one chart
			overTime, err := InterestOverTime(context.Background(), explore[0], langEN)
			assert.NoError(t, err)
			assert.True(t, len(overTime) > 0)

			// interest over time for 3 compared items in one map
			byLoc, err := InterestByLocation(context.Background(), explore[1], langEN)
			assert.NoError(t, err)
			assert.True(t, len(byLoc) > 0)
		}()
	}
}

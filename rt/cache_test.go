package rt

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
)

type setTester struct{}

func (c *setTester) Run(job Job, ctx *Ctx) (interface{}, error) {
	data := job.Bytes()

	if err := ctx.Cache.Set("important", data, 1); err != nil {
		return nil, err
	}

	return nil, nil
}

// OnChange runs on worker changes
func (c *setTester) OnChange(_ ChangeEvent) error {
	return nil
}

type getTester struct{}

func (c *getTester) Run(job Job, ctx *Ctx) (interface{}, error) {
	key := job.String()

	val, err := ctx.Cache.Get(key)
	if err != nil {
		return nil, err
	}

	return string(val), nil
}

// OnChange runs on worker changes
func (c *getTester) OnChange(_ ChangeEvent) error {
	return nil
}

func TestCacheGetSet(t *testing.T) {
	r := setupReactrWithCache()
	r.Register("set", &setTester{})
	r.Register("get", &getTester{})

	_, err := r.Do(NewJob("set", "very important information")).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set"))
		return
	}

	val, err := r.Do(NewJob("get", "important")).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "get job failed"))
		return
	}

	if val.(string) != "very important information" {
		t.Error("result did not match expected 'very important information': ", val.(string))
	}
}

func TestCacheGetSetWithTTL(t *testing.T) {
	r := setupReactrWithCache()
	r.Register("set", &setTester{})
	r.Register("get", &getTester{})

	_, err := r.Do(NewJob("set", "very important information")).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to set"))
		return
	}

	<-time.After(time.Second * 2)

	_, err = r.Do(NewJob("get", "important")).Then()
	if err == nil {
		t.Error("should have errored, did not")
		return
	}
}

func setupReactrWithCache() *Reactr {
	var r *Reactr

	if _, exists := os.LookupEnv("REACTR_TEST_REDIS"); exists {
		fmt.Println("using Redis for Reactr cache tests")

		config := rcap.DefaultCapabilityConfig()
		config.Cache = &rcap.CacheConfig{
			Enabled: true,
			RedisConfig: &rcap.RedisConfig{
				ServerAddress: "localhost:6379",
			},
		}

		r = NewWithConfig(config)
	} else {
		r = New()
	}

	return r
}

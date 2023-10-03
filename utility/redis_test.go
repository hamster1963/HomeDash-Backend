package utility

import (
	"context"
	"errors"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/google/uuid"
	"testing"
	"time"
)

var (
	// ErrLockFailed 加锁失败
	ErrLockFailed = errors.New("尝试加锁失败")
	// ErrTimeout 加锁超时
	ErrTimeout = errors.New("timeout")
)

var (
	config = gredis.Config{
		Address: "120.24.211.49:6379",
		Db:      1,
		Pass:    "deny1963",
	}
	group = "cache"
	ctx   = gctx.New()
)

type Locker struct {
	redisClient *gredis.Redis
	ttl         int64
	tryInterval int
}

func NewLocker(redisClient *gredis.Redis, ttl int64, tryInterval int) *Locker {
	return &Locker{redisClient: redisClient, ttl: ttl, tryInterval: tryInterval}
}

type Lock struct {
	redisClient *gredis.Redis
	resource    string
	uuid        string
	ttl         int64
	tryInterval int
}

func (l *Locker) GetLock(resource string) *Lock {
	return &Lock{
		redisClient: l.redisClient,
		resource:    resource,
		uuid:        uuid.NewString(),
		ttl:         l.ttl,
		tryInterval: l.tryInterval,
	}

}

func (l *Lock) TryLock(ctx context.Context) error {
	lockStatus, err := l.redisClient.SetNX(ctx, l.resource, l.uuid)
	if err != nil {
		return err
	}
	if lockStatus {
		err := l.redisClient.SetEX(ctx, l.resource, l.uuid, l.ttl)
		if err != nil {
			return err
		}
	} else {
		return ErrLockFailed
	}
	return nil
}

func (l *Lock) UnLock(ctx context.Context) error {
	lockStatus, err := l.redisClient.Get(ctx, l.resource)
	if err != nil {
		return err
	}
	if lockStatus.String() != l.uuid {
		return ErrLockFailed
	}
	_, err = l.redisClient.Del(ctx, l.resource)
	if err != nil {
		return err
	}
	return nil
}

func (l *Lock) Lock(ctx context.Context) error {
	err := l.TryLock(ctx)
	if err == nil {
		return nil
	}
	if err != ErrLockFailed {
		return err
	}
	glog.Info(ctx, "尝试加锁失败，开始重试")
	// 尝试加锁失败，开始重试
	ticker := time.NewTicker(time.Duration(l.tryInterval) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		case <-ticker.C:
			if err := l.TryLock(ctx); err != nil {
				continue
			}
			return nil
		}
	}
}

// 测试方法注释
func Test12_57_15(t *testing.T) {
	gredis.SetConfig(&config, group)
	redisClient := g.Redis(group)
	locker := NewLocker(redisClient, 6, 200)
	lock := locker.GetLock("hamster")

	err := lock.TryLock(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	err = lock.Lock(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("加锁成功")
	err = lock.UnLock(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("解锁成功")
}

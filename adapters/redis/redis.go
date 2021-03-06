package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

/*type Store struct {
	client *redis.Client
	prefix string
	ctx context.Context
}

func NewRedisStore(ctx context.Context, client *redis.Client, prefix string) *Store {
	if client == nil {
		panic("redis client is nil")
	}

	return &Store{
		client: client,
		prefix: prefix,
		ctx: ctx,
	}
}

func (r *Store) token(token string) string {
	return r.prefix + ":" + token
}

func (r *Store) Find(token string) (b []byte, exists bool, err error) {
	logrus.WithFields(logrus.Fields{
		"token": r.token(token),
	}).Debug("finding session data")
	result, err := r.client.Get(r.ctx, r.token(token)).Result()
	switch {
	case err == redis.Nil:
		logrus.WithFields(logrus.Fields{
			"token": r.token(token),
		}).Warn("session data is nil")
		return nil, false, nil
	case err != nil:
		logrus.WithFields(logrus.Fields{
			"token": r.token(token),
		}).Warnf("unknown error occurred while retrieving session data: %v", err)
		return nil, false, err
	}

	return []byte(result), true, nil
}

func (r *Store) Commit(token string, b []byte, expiry time.Time) error {
	logrus.WithFields(logrus.Fields{
		"token": r.token(token),
		"expiry": expiry,
	}).Debugf("committing data: %+v", string(b))

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logrus.Warnf("cannot load location: %v", err)
	}

	currentTime := time.Now().In(loc)
	expiry = expiry.In(loc)
	duration := expiry.Sub(currentTime)

	logrus.WithFields(logrus.Fields{
		"token": r.token(token),
		"value": string(b),
		"willExpireAt": expiry,
		"duration": duration,
	}).Debug("saving session data")
	err = r.client.Set(r.ctx, r.token(token), b, duration).Err()
	if err != nil {
		logrus.Warnf("saving session data failed: %v", err)
	}
	return err
}

func (r *Store) Delete(token string) error {
	logrus.WithFields(logrus.Fields{
		"token": r.token(token),
	}).Debug("deleting session data")

	err := r.client.Del(r.ctx, r.token(token)).Err()
	if err != nil {
		logrus.Warnf("error occurred while deleting session data: %v", err)
	}
	return err
}*/

type Store struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

func NewRedisStore(client *redis.Client, ctx context.Context) *Store {
	if client == nil {
		panic("redis client is nil")
	}
	return NewRedisStoreWithPrefix(client, ctx, "scs:session")
}

func NewRedisStoreWithPrefix(client *redis.Client, ctx context.Context, prefix string) *Store {
	return &Store{
		client: client,
		ctx:    ctx,
		prefix: prefix,
	}
}

func (r *Store) token(token string) string {
	return r.prefix + ":" + token
}

func (r *Store) Find(token string) (b []byte, exists bool, err error) {
	logrus.WithFields(logrus.Fields{
		"sessionToken": r.token(token),
	}).Debug("Finding session data")
	bytes, err := r.client.Get(r.ctx, r.token(token)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debugf("Finding the value for token: %s caused an error", r.token(token))
		return nil, false, nil
	}

	return bytes, true, nil
}

func (r *Store) Commit(token string, b []byte, expiry time.Time) error {
	logrus.WithFields(logrus.Fields{
		"sessionToken": token,
	}).Debug("Saving session to redis")
	currentTime := time.Now()

	expiryDuration := expiry.Sub(currentTime)
	logrus.WithFields(logrus.Fields{
		"sessionToken": r.token(token),
		"value":        string(b),
		"willExpireAt": expiry,
		"duration":     expiryDuration,
	}).Debug("Saving session data")
	return r.client.Set(r.ctx, r.token(token), b, expiryDuration).Err()
}

func (r *Store) Delete(token string) error {
	logrus.WithFields(logrus.Fields{
		"sessionToken": token,
	}).Debug("Deleting session")

	result, err := r.client.Del(r.ctx, r.token(token)).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"result": result,
			"error":  err,
		}).Debug("Deleting session token caused error")
	}
	return nil
}


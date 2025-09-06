package memcache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	lib_store "github.com/eko/gocache/lib/v4/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMemcache(t *testing.T) {
	// Given
	client := NewMockMemcacheClientInterface(t)

	// When
	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// Then
	assert.IsType(t, new(MemcacheStore), store)
	assert.Equal(t, client, store.client)
	assert.Equal(t, &lib_store.Options{Expiration: 3 * time.Second}, store.options)
}

func TestMemcacheGet(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(&memcache.Item{
		Value: cacheValue,
	}, nil)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, err := store.Get(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestMemcacheGetWithMissingItem(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(nil, memcache.ErrCacheMiss)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, err := store.Get(ctx, cacheKey)

	// Then
	assert.ErrorIs(t, err, lib_store.NotFound{})
	assert.Nil(t, value)
}

func TestMemcacheGetWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	expectedErr := errors.New("an unexpected error occurred")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(nil, expectedErr)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, err := store.Get(ctx, cacheKey)

	// Then
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
}

func TestMemcacheGetWithTTL(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(&memcache.Item{
		Value:      cacheValue,
		Expiration: int32(5),
	}, nil)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, ttl, err := store.GetWithTTL(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
	assert.Equal(t, 5*time.Second, ttl)
}

func TestMemcacheGetWithTTLWhenMissingItem(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(nil, memcache.ErrCacheMiss)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, ttl, err := store.GetWithTTL(ctx, cacheKey)

	// Then
	assert.ErrorIs(t, err, lib_store.NotFound{})
	assert.Nil(t, value)
	assert.Equal(t, 0*time.Second, ttl)
}

func TestMemcacheGetWithTTLWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	expectedErr := errors.New("an unexpected error occurred")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get(cacheKey).Return(nil, expectedErr)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	value, ttl, err := store.GetWithTTL(ctx, cacheKey)

	// Then
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
	assert.Equal(t, 0*time.Second, ttl)
}

func TestMemcacheSet(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Set(&memcache.Item{
		Key:        cacheKey,
		Value:      cacheValue,
		Expiration: int32(5),
	}).Return(nil)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	err := store.Set(ctx, cacheKey, cacheValue, lib_store.WithExpiration(5*time.Second))

	// Then
	assert.Nil(t, err)
}

func TestMemcacheSetWhenNoOptionsGiven(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Set(&memcache.Item{
		Key:        cacheKey,
		Value:      cacheValue,
		Expiration: int32(3),
	}).Return(nil)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	err := store.Set(ctx, cacheKey, cacheValue)

	// Then
	assert.Nil(t, err)
}

func TestMemcacheSetWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	expectedErr := errors.New("an unexpected error occurred")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Set(&memcache.Item{
		Key:        cacheKey,
		Value:      cacheValue,
		Expiration: int32(3),
	}).Return(expectedErr)

	store := NewMemcache(client, lib_store.WithExpiration(3*time.Second))

	// When
	err := store.Set(ctx, cacheKey, cacheValue)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestMemcacheSetWithTags(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	tagKey := "gocache_tag_tag1"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Set(mock.Anything).Return(nil)
	client.EXPECT().Get(tagKey).Return(nil, memcache.ErrCacheMiss)
	client.EXPECT().Add(&memcache.Item{
		Key:        tagKey,
		Value:      []byte(cacheKey),
		Expiration: int32(TagKeyExpiry.Seconds()),
	}).Return(nil)

	store := NewMemcache(client)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, lib_store.WithTags([]string{"tag1"}))

	// Then
	assert.Nil(t, err)
}

func TestMemcacheSetWithTagsWhenAlreadyInserted(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"
	cacheValue := []byte("my-cache-value")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Set(mock.Anything).Return(nil)
	client.EXPECT().Get("gocache_tag_tag1").Return(&memcache.Item{
		Value: []byte("my-key,a-second-key"),
	}, nil)

	store := NewMemcache(client)

	// When
	err := store.Set(ctx, cacheKey, cacheValue, lib_store.WithTags([]string{"tag1"}))

	// Then
	assert.Nil(t, err)
}

func TestMemcacheDelete(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Delete(cacheKey).Return(nil)

	store := NewMemcache(client)

	// When
	err := store.Delete(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
}

func TestMemcacheDeleteWithMissingItem(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKey := "my-key"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Delete(cacheKey).Return(memcache.ErrCacheMiss)

	store := NewMemcache(client)

	// When
	err := store.Delete(ctx, cacheKey)

	// Then
	assert.Nil(t, err)
}

func TestMemcacheDeleteWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	expectedErr := errors.New("unable to delete key")

	cacheKey := "my-key"

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Delete(cacheKey).Return(expectedErr)

	store := NewMemcache(client)

	// When
	err := store.Delete(ctx, cacheKey)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestMemcacheInvalidate(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKeys := &memcache.Item{
		Value: []byte("a23fdf987h2svc23,jHG2372x38hf74"),
	}

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get("gocache_tag_tag1").Return(cacheKeys, nil)
	client.EXPECT().Delete("a23fdf987h2svc23").Return(nil)
	client.EXPECT().Delete("jHG2372x38hf74").Return(nil)

	store := NewMemcache(client)

	// When
	err := store.Invalidate(ctx, lib_store.WithInvalidateTags([]string{"tag1"}))

	// Then
	assert.Nil(t, err)
}

func TestMemcacheInvalidateWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	cacheKeys := &memcache.Item{
		Value: []byte("a23fdf987h2svc23,jHG2372x38hf74"),
	}

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().Get("gocache_tag_tag1").Return(cacheKeys, nil)
	client.EXPECT().Delete("a23fdf987h2svc23").Return(errors.New("unexpected error"))
	client.EXPECT().Delete("jHG2372x38hf74").Return(nil)

	store := NewMemcache(client)

	// When
	err := store.Invalidate(ctx, lib_store.WithInvalidateTags([]string{"tag1"}))

	// Then
	assert.Nil(t, err)
}

func TestMemcacheClear(t *testing.T) {
	// Given
	ctx := context.Background()

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().FlushAll().Return(nil)

	store := NewMemcache(client)

	// When
	err := store.Clear(ctx)

	// Then
	assert.Nil(t, err)
}

func TestMemcacheClearWhenError(t *testing.T) {
	// Given
	ctx := context.Background()

	expectedErr := errors.New("an unexpected error occurred")

	client := NewMockMemcacheClientInterface(t)
	client.EXPECT().FlushAll().Return(expectedErr)

	store := NewMemcache(client)

	// When
	err := store.Clear(ctx)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestMemcacheGetType(t *testing.T) {
	// Given
	client := NewMockMemcacheClientInterface(t)

	store := NewMemcache(client)

	// When - Then
	assert.Equal(t, MemcacheType, store.GetType())
}

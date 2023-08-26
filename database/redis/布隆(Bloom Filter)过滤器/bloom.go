package bloom

import (
	"context"
	"errors"
	"strconv"

	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
// maps as k in the error rate table
// 表示经过多少散列函数计算
const maps = 14

var (
	// ErrTooLargeOffset indicates the offset is too large in bitset.
	ErrTooLargeOffset = errors.New("too large offset")

	// 用lua脚本呢保证整个操作是原子性执行的
	//通过redis bitmap 数据结构实现
	//ARGV:偏移量offset数组
	//KYES[1]: setbit操作的key
	//全部设置为1
	setScript = redis.NewScript(`
		for _, offset in ipairs(ARGV) do
			redis.call("setbit", KEYS[1], offset, 1)
		end
	`)

	//ARGV:偏移量offset数组
	//KYES[1]: setbit操作的key
	//检查是否全部为1
	testScript = redis.NewScript(`
		for _, offset in ipairs(ARGV) do
			if tonumber(redis.call("getbit", KEYS[1], offset)) == 0 then
				return false
			end
		end
		return true
	`)
)

type (
	// A Filter is a bloom filter.
	// 定义布隆过滤器结构体
	Filter struct {
		bits   uint
		bitSet bitSetProvider
	}

	// 位数组操作接口定义
	bitSetProvider interface {
		check(ctx context.Context, offsets []uint) (bool, error)
		set(ctx context.Context, offsets []uint) error
	}
)

// New create a Filter, store is the backed redis, key is the key for the bloom filter,
// bits is how many bits will be used, maps is how many hashes for each addition.
// best practices:
// elements - means how many actual elements
// when maps = 14, formula: 0.7*(bits/maps), bits = 20*elements, the error rate is 0.000067 < 1e-4
// for detailed error rate table, see http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
func New(store *redis.Redis, key string, bits uint) *Filter {
	return &Filter{
		bits:   bits,
		bitSet: newRedisBitSet(store, key, bits),
	}
}

// Add adds data into f.
func (f *Filter) Add(data []byte) error {
	return f.AddCtx(context.Background(), data)
}

// AddCtx adds data into f with context.
func (f *Filter) AddCtx(ctx context.Context, data []byte) error {
	locations := f.getLocations(data)
	return f.bitSet.set(ctx, locations)
}

// Exists checks if data is in f.
func (f *Filter) Exists(data []byte) (bool, error) {
	return f.ExistsCtx(context.Background(), data)
}

// ExistsCtx checks if data is in f with context.
func (f *Filter) ExistsCtx(ctx context.Context, data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.check(ctx, locations)
	if err != nil {
		return false, err
	}

	return isSet, nil
}

// k次散列计算出k个offset
func (f *Filter) getLocations(data []byte) []uint {
	locations := make([]uint, maps)
	for i := uint(0); i < maps; i++ {
		//哈希计算,使用的是"MurmurHash3"算法,并每次追加一个固定的i字节进行计算
		hashValue := hash.Hash(append(data, byte(i)))
		////取下标offset
		locations[i] = uint(hashValue % uint64(f.bits))
	}

	return locations
}

type redisBitSet struct {
	store *redis.Redis
	key   string
	bits  uint
}

func newRedisBitSet(store *redis.Redis, key string, bits uint) *redisBitSet {
	return &redisBitSet{
		store: store,
		key:   key,
		bits:  bits,
	}
}

func (r *redisBitSet) buildOffsetArgs(offsets []uint) ([]string, error) {
	var args []string

	for _, offset := range offsets {
		if offset >= r.bits {
			return nil, ErrTooLargeOffset
		}

		args = append(args, strconv.FormatUint(uint64(offset), 10))
	}

	return args, nil
}

func (r *redisBitSet) check(ctx context.Context, offsets []uint) (bool, error) {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return false, err
	}

	////执行脚本
	resp, err := r.store.ScriptRunCtx(ctx, testScript, []string{r.key}, args)
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	exists, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return exists == 1, nil
}

// del only use for testing.
func (r *redisBitSet) del() error {
	_, err := r.store.Del(r.key)
	return err
}

// expire only use for testing.
func (r *redisBitSet) expire(seconds int) error {
	return r.store.Expire(r.key, seconds)
}

func (r *redisBitSet) set(ctx context.Context, offsets []uint) error {
	args, err := r.buildOffsetArgs(offsets)
	if err != nil {
		return err
	}

	_, err = r.store.ScriptRunCtx(ctx, setScript, []string{r.key}, args)
	if err == redis.Nil {
		return nil
	}

	return err
}

package tools

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	epochRedis    = 1609459200 // 自定义纪元开始时间(2021-01-01 00:00:00 UTC)，可以减少时间戳位数
	timestampBits = 41         // 时间戳位数(41位可用约69年)
	machineIDBits = 10         // 机器ID位数(最多1024台机器)
	sequenceBits  = 12         // 序列号位数(每毫秒4096个ID)

	maxMachineID   = -1 ^ (-1 << machineIDBits)
	maxSequenceNum = -1 ^ (-1 << sequenceBits)

	timestampShift = machineIDBits + sequenceBits
	machineIDShift = sequenceBits
)

type RedisInt64IDGenerator struct {
	client        *redis.Redis
	machineID     int64
	lastTimestamp int64
	sequence      int64
	mu            sync.Mutex
}

func NewRedisInt64IDGenerator(client *redis.Redis, machineID int64) (*RedisInt64IDGenerator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}

	return &RedisInt64IDGenerator{
		client:    client,
		machineID: machineID,
	}, nil
}

func (g *RedisInt64IDGenerator) GenerateInt64ID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	currentTimestamp := time.Now().Unix() - epochRedis

	// 获取当前毫秒时间戳对应的Redis key
	redisKey := fmt.Sprintf("id_gen:%d", currentTimestamp)

	ctx := context.Background()

	// 使用Lua脚本保证原子性操作
	luaScript := `
	local key = KEYS[1]
	local machineID = tonumber(ARGV[1])
	local maxSeq = tonumber(ARGV[2])
	
	local seq = redis.call("INCR", key)
	if seq == 1 then
		redis.call("EXPIRE", key, 5)
	end
	
	if seq > maxSeq then
		return nil
	end
	
	return (tonumber(ARGV[3]) * 2^22 + machineID * 2^12 + seq) - 1
	`

	// 执行Lua脚本
	result, err := g.client.EvalCtx(ctx, luaScript, []string{redisKey},
		g.machineID, maxSequenceNum, currentTimestamp)
	if err != nil {
		return 0, fmt.Errorf("redis lua script failed: %v", err)
	}

	if result == nil {
		return 0, errors.New("sequence number overflow, wait for next millisecond")
	}

	id, ok := result.(int64)
	if !ok {
		return 0, errors.New("invalid id type returned from redis")
	}

	return id, nil
}

func ParseInt64ID(id int64) (timestamp int64, machineID int64, sequence int64) {
	id += 1 // 恢复原始值
	timestamp = (id >> timestampShift) + epochRedis
	machineID = (id >> machineIDShift) & maxMachineID
	sequence = id & maxSequenceNum
	return
}

package tools

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	epochRedis    = 1672531200 // 2023-01-01 00:00:00 UTC，作为自定义纪元起点
	timestampBits = 8          // 时间戳位数(8位=256个时间单位)
	machineIDBits = 2          // 机器ID位数(4台机器)
	sequenceBits  = 10         // 序列号位数(1024个序列号)

	maxMachineID   = -1 ^ (-1 << machineIDBits)
	maxSequenceNum = -1 ^ (-1 << sequenceBits)

	timestampShift = machineIDBits + sequenceBits
	machineIDShift = sequenceBits

	timeUnit = time.Second * 10 // 每个时间戳单位代表10秒
)

type SimpleRedisIDGenerator struct {
	client    *redis.Redis
	machineID int64
	mu        sync.Mutex
}

func NewSimpleRedisIDGenerator(client *redis.Redis, machineID int64) (*SimpleRedisIDGenerator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}

	return &SimpleRedisIDGenerator{
		client:    client,
		machineID: machineID,
	}, nil
}

func (g *SimpleRedisIDGenerator) GenerateID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 获取当前时间戳(以10秒为单位)
	currentTime := time.Now().Unix()
	currentTimestamp := (currentTime - epochRedis) / int64(timeUnit.Seconds())

	if currentTimestamp > (1<<timestampBits)-1 {
		return 0, errors.New("timestamp overflow, epoch needs adjustment")
	}

	// 生成Redis key
	redisKey := fmt.Sprintf("id_gen_session:%d", currentTimestamp)

	// 使用Redis的INCR命令获取自增序列号
	seq, err := g.client.Incr(redisKey)
	if err != nil {
		return 0, fmt.Errorf("redis incr failed: %v", err)
	}

	// 设置key的过期时间(至少保留3个时间单位，即30秒)
	if seq == 1 {
		err := g.client.Expire(redisKey, 3)
		if err != nil {
			return 0, err
		}
	}

	// 检查序列号是否超过最大值
	if seq > maxSequenceNum {
		return 0, errors.New("sequence number overflow, wait for next time unit")
	}

	// 组合ID: 时间戳(8位) | 机器ID(2位) | 序列号(10位)
	id := (currentTimestamp << timestampShift) |
		(g.machineID << machineIDShift) |
		seq

	return id, nil
}

func ParseID(id int64) (timestamp int64, machineID int64, sequence int64) {
	timestamp = (id >> timestampShift) & 0xFF         // 取8位
	machineID = (id >> machineIDShift) & maxMachineID // 取2位
	sequence = id & maxSequenceNum                    // 取10位

	// 转换为实际时间戳
	timestamp = timestamp*int64(timeUnit.Seconds()) + epochRedis
	return
}

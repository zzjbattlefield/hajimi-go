package sync

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/zzjbattlefield/hajimi-go/internal/extractor"
	"github.com/zzjbattlefield/hajimi-go/internal/logger"
)

// SyncItem 表示要同步的项目
type SyncItem struct {
	Secret    extractor.Secret `json:"secret"`
	Timestamp time.Time        `json:"timestamp"`
	Attempts  int              `json:"attempts"`
}

// SyncManager 管理发现密钥的同步
type SyncManager struct {
	queue      chan *SyncItem
	batchQueue chan []*SyncItem
	logger     *logger.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex

	// 配置
	batchSize     int
	flushInterval time.Duration
	maxRetries    int
}

// NewSyncManager 创建一个新的 SyncManager
func NewSyncManager(logger *logger.Logger, batchSize int, flushInterval time.Duration, maxRetries int) *SyncManager {
	ctx, cancel := context.WithCancel(context.Background())

	sm := &SyncManager{
		queue:         make(chan *SyncItem, 1000), // 用于单个项目的缓冲通道
		batchQueue:    make(chan []*SyncItem, 10), // 用于批次的缓冲通道
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		maxRetries:    maxRetries,
	}

	return sm
}

// AddToQueue 将密钥添加到同步队列中
func (sm *SyncManager) AddToQueue(secret extractor.Secret) {
	item := &SyncItem{
		Secret:    secret,
		Timestamp: time.Now(),
		Attempts:  0,
	}

	select {
	case sm.queue <- item:
		sm.logger.Debugf("已将密钥添加到同步队列: %s", secret.Type)
	default:
		sm.logger.Warnf("同步队列已满，丢弃密钥: %s", secret.Type)
	}
}

// Start 启动同步过程
func (sm *SyncManager) Start() {
	sm.wg.Add(2)

	// 启动队列处理器
	go sm.processQueue()

	// 启动批次处理器
	go sm.processBatches()

	sm.logger.Info("同步管理器已启动")
}

// Stop 停止同步过程
func (sm *SyncManager) Stop() {
	sm.logger.Info("正在停止同步管理器...")
	sm.cancel()
	sm.wg.Wait()
	sm.logger.Info("同步管理器已停止")
}

// processQueue 处理队列中的项目并将其分批
func (sm *SyncManager) processQueue() {
	defer sm.wg.Done()

	ticker := time.NewTicker(sm.flushInterval)
	defer ticker.Stop()

	batch := make([]*SyncItem, 0, sm.batchSize)

	for {
		select {
		case item := <-sm.queue:
			batch = append(batch, item)

			// 如果我们达到了批次大小，发送批次
			if len(batch) >= sm.batchSize {
				sm.sendBatch(batch)
				batch = make([]*SyncItem, 0, sm.batchSize)
			}

		case <-ticker.C:
			// 定时刷新批次
			if len(batch) > 0 {
				sm.sendBatch(batch)
				batch = make([]*SyncItem, 0, sm.batchSize)
			}

		case <-sm.ctx.Done():
			// 停止时刷新任何剩余项目
			if len(batch) > 0 {
				sm.sendBatch(batch)
			}
			return
		}
	}
}

// sendBatch 将一批项目发送到批次队列
func (sm *SyncManager) sendBatch(batch []*SyncItem) {
	select {
	case sm.batchQueue <- batch:
		sm.logger.Debugf("已将 %d 个项目批次发送到批次处理器", len(batch))
	default:
		sm.logger.Warnf("批次队列已满，丢弃 %d 个项目批次", len(batch))
		// 在实际实现中，您可能希望以不同的方式处理此情况
		// 例如，您可以重试或保存到备份存储
	}
}

// processBatches 处理项目批次
func (sm *SyncManager) processBatches() {
	defer sm.wg.Done()

	// 需要重试的失败项目
	retryQueue := make([]*SyncItem, 0)

	for {
		select {
		case batch := <-sm.batchQueue:
			// 与需要重试的任何项目合并
			combinedBatch := append(batch, retryQueue...)
			retryQueue = make([]*SyncItem, 0) // 清空重试队列

			// 处理合并后的批次
			failedItems := sm.processBatch(combinedBatch)

			// 处理失败的项目
			for _, item := range failedItems {
				item.Attempts++
				if item.Attempts < sm.maxRetries {
					sm.logger.Debugf("重试项目 (尝试 %d): %s", item.Attempts, item.Secret.Type)
					retryQueue = append(retryQueue, item)
				} else {
					sm.logger.Errorf("项目已超过最大重试次数: %s", item.Secret.Type)
					// 在实际实现中，您可能希望将其保存到死信队列
				}
			}

		case <-sm.ctx.Done():
			// 在停止之前尝试处理重试队列中的任何剩余项目
			if len(retryQueue) > 0 {
				sm.logger.Infof("在停止之前处理 %d 个剩余项目", len(retryQueue))
				sm.processBatch(retryQueue)
			}
			return
		}
	}
}

// processBatch 处理一批项目
// 在实际实现中，这会将数据发送到外部服务
// 返回处理失败的项目切片
func (sm *SyncManager) processBatch(batch []*SyncItem) []*SyncItem {
	sm.logger.Infof("正在处理 %d 个项目的批次", len(batch))

	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	// 在实际实现中，您会在这里将批次发送到外部服务
	// 例如，发送到 webhook、API 或消息队列

	var failedItems []*SyncItem

	// 记录批次内容以进行演示
	for _, item := range batch {
		// 模拟偶尔的失败以进行演示
		if item.Attempts == 0 && item.Secret.Type == "github_token" {
			// 模拟在第一次尝试时 github 令牌的失败
			sm.logger.Debugf("模拟 github 令牌失败 (尝试 %d)", item.Attempts)
			failedItems = append(failedItems, item)
			continue
		}

		// 记录项目（在实际实现中，您会将其发送到外部服务）
		jsonData, err := json.Marshal(item)
		if err != nil {
			sm.logger.Errorf("无法序列化同步项目: %v", err)
			failedItems = append(failedItems, item)
			continue
		}

		sm.logger.Debugf("已处理同步项目: %s", string(jsonData))
	}

	if len(failedItems) > 0 {
		sm.logger.Warnf("未能处理批次中的 %d 个项目", len(failedItems))
	}

	sm.logger.Infof("已完成处理 %d 个项目的批次 (%d 个失败)", len(batch), len(failedItems))

	return failedItems
}

// GetQueueLength 返回队列的当前长度
func (sm *SyncManager) GetQueueLength() int {
	return len(sm.queue)
}

// GetBatchQueueLength 返回批次队列的当前长度
func (sm *SyncManager) GetBatchQueueLength() int {
	return len(sm.batchQueue)
}

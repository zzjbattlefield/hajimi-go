package data

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zzjbattlefield/hajimi-go/internal/config"
	"github.com/zzjbattlefield/hajimi-go/internal/logger"
	"github.com/zzjbattlefield/hajimi-go/internal/validator"
)

var freeFile *os.File
var payFile *os.File

var CacheData *Cache

type Cache struct {
	data   map[string]struct{}
	rw     sync.RWMutex
	signal chan struct{}
}

// Size 返回cache中项目的数量
func (c *Cache) Size() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.data)
}

var once sync.Once

func init() {
	once.Do(func() {
		var err error
		cache := &Cache{
			data:   make(map[string]struct{}),
			rw:     sync.RWMutex{},
			signal: make(chan struct{}),
		}
		CacheData = cache

		// 初始化 freeFile
		freeFile, err = initFile(config.Conf.FreeFileName, cache)
		if err != nil {
			panic(fmt.Errorf("初始化免费文件失败: %w", err))
		}

		// 初始化 payFile
		payFile, err = initFile(config.Conf.PayFileName, cache)
		if err != nil {
			panic(fmt.Errorf("初始化付费文件失败: %w", err))
		}

		// 启动一个 goroutine 定时保存数据
		go func() {
			ticker := time.NewTicker(time.Second * 30)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// 定时保存
					if err := freeFile.Sync(); err != nil {
						logger.Log.Errorf("定时保存免费数据失败: %v", err)
					}
					if err := payFile.Sync(); err != nil {
						logger.Log.Errorf("定时保存付费数据失败: %v", err)
					}
				case <-cache.signal:
					// 关闭时保存
					if err := freeFile.Sync(); err != nil {
						logger.Log.Errorf("保存免费数据失败: %v", err)
					}
					if err := payFile.Sync(); err != nil {
						logger.Log.Errorf("保存付费数据失败: %v", err)
					}
					return
				}
			}
		}()
	})
}

func initFile(fileName string, cache *Cache) (*os.File, error) {
	// 构建文件路径
	filePath := filepath.Join(config.Conf.DataPath, fileName)
	// 确保数据目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	// 先以只读模式打开文件来读取现有数据
	readFile, err := os.Open(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	// 如果文件存在，则读取现有数据到cache中
	if err == nil {
		defer readFile.Close()
		scanner := bufio.NewScanner(readFile)
		cache.rw.Lock()
		for scanner.Scan() {
			content := scanner.Text()
			cache.data[content] = struct{}{}
		}
		cache.rw.Unlock()

		// 检查扫描过程中是否有错误
		if err := scanner.Err(); err != nil {
			logger.Log.Errorf("读取文件时发生错误: %v", err)
		}
	}

	// 以追加写入模式重新打开文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件用于写入失败: %w", err)
	}
	return file, nil
}

func (c *Cache) Add(secretResult *validator.ValidationResult) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, exists := c.data[secretResult.Value]; !exists {
		c.data[secretResult.Value] = struct{}{}
		if !secretResult.Pay {
			freeFile.WriteString(secretResult.Value + "\n")
		} else {
			payFile.WriteString(secretResult.Value + "\n")
		}
	}
}

func (c *Cache) Check(secret string) bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	_, ok := c.data[secret]
	return ok
}

func (c *Cache) Clean() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.data = make(map[string]struct{})
	freeFile.Truncate(0)
	freeFile.Seek(0, 0)
	payFile.Truncate(0)
	payFile.Seek(0, 0)
}

func (c *Cache) Close() {
	close(c.signal)
	freeFile.Close()
	payFile.Close()
}

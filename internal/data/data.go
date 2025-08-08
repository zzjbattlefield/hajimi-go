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
)

var saveFile *os.File

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
		cache := &Cache{
			data:   make(map[string]struct{}),
			rw:     sync.RWMutex{},
			signal: make(chan struct{}),
		}
		CacheData = cache
		// 构建文件路径
		filePath := filepath.Join(config.Conf.DataPath, config.Conf.SaveFileName)
		// 确保数据目录存在
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Errorf("创建数据目录失败: %w", err))
		}

		// 先以只读模式打开文件来读取现有数据
		readFile, err := os.Open(filePath)
		if err != nil && !os.IsNotExist(err) {
			panic(fmt.Errorf("打开文件失败: %w", err))
		}

		// 如果文件存在，则读取现有数据到cache中
		if err == nil {
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

			// 关闭只读文件
			readFile.Close()
		}

		// 以追加写入模式重新打开文件
		saveFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("打开文件用于写入失败: %w", err))
		}

		// 启动一个 goroutine 定时保存数据
		go func() {
			ticker := time.NewTicker(time.Second * 30)
			select {
			case <-ticker.C:
				// 定时保存
				if err := saveFile.Sync(); err != nil {
					logger.Log.Errorf("定时保存数据失败: %v", err)
				}
			case <-cache.signal:
				// 关闭时保存
				if err := saveFile.Sync(); err != nil {
					logger.Log.Errorf("保存数据失败: %v", err)
				}
				return
			}
		}()
	})
}

func (c *Cache) Add(secret string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if _, exists := c.data[secret]; !exists {
		c.data[secret] = struct{}{}
		saveFile.WriteString(secret + "\n")
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
	saveFile.Truncate(0)
	saveFile.Seek(0, 0)
}

func (c *Cache) Close() {
	close(c.signal)
	saveFile.Close()
}

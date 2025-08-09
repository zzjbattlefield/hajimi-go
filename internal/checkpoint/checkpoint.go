package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zzjbattlefield/hajimi-go/internal/config"
)

// Checkpoint 结构体定义了扫描过程中的状态快照。
// 这个结构体用于保存和恢复扫描进度，以便在扫描任务意外中断（例如程序崩溃或手动停止）后，
// 能够从上次记录的位置精确地继续执行，避免重复工作和数据丢失。
type Checkpoint struct {
	// Query 是当前正在处理的GitHub搜索查询语句。
	// 这个字段记录了扫描任务所使用的具体查询条件。
	Query Query `json:"query"`

	// LastPage 是上次成功获取并处理完的搜索结果的页码。
	// GitHub API 的搜索结果是分页的，这个字段帮助我们从下一页继续。
	LastPage int `json:"last_page"`

	// LastRepo 是在当前结果页面中，上次成功处理的最后一个仓库的完整名称（例如 "owner/repo"）。
	// 当需要从特定仓库继续处理时，此字段非常关键。
	LastRepo string `json:"last_repo"`

	// LastFile 是在特定仓库中，上次成功处理的最后一个文件的路径。
	// 这使得我们可以从仓库中的特定文件继续扫描。
	LastFile string `json:"last_file"`

	// LastCommit 是与上次处理的文件或仓库关联的最后一次提交的SHA哈希值。
	// 这可以用于验证文件自上次扫描以来是否发生了变化。
	LastCommit string `json:"last_commit"`

	// LastUpdated 是该检查点最后一次被更新的时间戳。
	// 这个信息对于监控扫描进度和调试非常有用。
	LastUpdated time.Time `json:"last_updated"`

	// TotalResults 是当前查询返回的总结果数。
	// 这个字段提供了对扫描范围的整体了解。
	TotalResults int `json:"total_results"`

	// Processed 是已经成功处理的结果数量。
	// 通过与 TotalResults 比较，可以计算出扫描的完成百分比。
	Processed int `json:"processed"`
}

// Manager 结构体负责管理检查点的所有操作，包括加载、保存和更新。
// 它封装了与检查点文件交互的底层逻辑，并通过一个读写互斥锁（RWMutex）
// 来确保对检查点数据的并发访问安全，防止在多线程环境下出现数据竞争问题。
type Manager struct {
	dataPath    string
	checkpoints checkpoints
	mu          sync.RWMutex
	querys      []Query
	queryIndex  int
}

type checkpoints map[Query]*Checkpoint

// NewManager 创建并返回一个新的检查点管理器（Manager）实例。
// 它需要一个配置对象（*config.Config）作为参数，该配置对象应包含
// 检查点文件的存储路径等信息。
func NewManager() *Manager {
	dataPath := config.Conf.DataPath
	querys := ReadQueryFile(config.Conf.QueriesFile)
	manager := &Manager{
		dataPath:   dataPath,
		querys:     querys,
		queryIndex: -1,
	}
	manager.load(querys)
	return manager
}

func newCheckpoints(querys []Query) checkpoints {
	checkpoints := make(checkpoints, len(querys))
	for _, query := range querys {
		checkpoints[query] = &Checkpoint{
			Query: query,
		}
	}
	return checkpoints
}

func (m *Manager) QueryNext() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queryIndex++
	return m.queryIndex < len(m.querys)
}

func (m *Manager) Query() Query {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.queryIndex < 0 || m.queryIndex >= len(m.querys) {
		return ""
	}
	return m.querys[m.queryIndex]
}

// load 方法从文件中加载检查点状态。
// 如果检查点文件（默认为 "checkpoint.json"）不存在，它会初始化一个新的、空的检查点。
// 这个方法通常在扫描任务开始时调用，用于恢复之前的进度或从头开始。
func (m *Manager) load(querys []Query) {
	m.mu.Lock()
	defer m.mu.Unlock()

	checkpointFile := filepath.Join(m.dataPath, "checkpoint.json")
	if _, err := os.Stat(checkpointFile); os.IsNotExist(err) {
		// 如果检查点文件不存在，则创建一个新的检查点，表示从头开始。
		// 这种情况发生在首次运行或检查点文件被手动删除后。
		m.checkpoints = newCheckpoints(querys)
		return
	}

	data, err := os.ReadFile(checkpointFile)
	if err != nil {
		panic(fmt.Errorf("读取检查点文件失败: %w", err))
	}

	var checkpointsData checkpoints
	newCheckpoints := make(checkpoints)
	if err := json.Unmarshal(data, &checkpointsData); err != nil {
		// 如果JSON解析失败，可能意味着文件已损坏。
		panic(fmt.Errorf("解析检查点数据失败: %w", err))
	}

	// 检查是queries文件是否有新增的类别
	for _, query := range querys {
		if _, ok := checkpointsData[query]; !ok {
			newCheckpoints[query] = &Checkpoint{
				Query: query,
			}
		} else {
			newCheckpoints[query] = checkpointsData[query]
		}
	}

	m.checkpoints = newCheckpoints
}

// Save 方法将当前的检查点状态以易于阅读的JSON格式保存到文件中。
// 这个方法应该在每次检查点状态更新后，或在程序准备正常退出前被调用，
// 以确保最新的扫描进度被持久化，防止数据丢失。
func (m *Manager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	checkpointFile := filepath.Join(m.dataPath, "checkpoint.json")

	// 确保数据目录存在，如果不存在则以权限 0755 创建它。
	// 这是为了防止因目录不存在而导致文件写入失败。
	if err := os.MkdirAll(m.dataPath, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	data, err := json.MarshalIndent(m.checkpoints, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化检查点数据为JSON格式时失败: %w", err)
	}

	if err := os.WriteFile(checkpointFile, data, 0644); err != nil {
		return fmt.Errorf("写入检查点文件失败: %w", err)
	}

	return nil
}

// Update 方法用新的扫描进度信息原子性地更新内存中的检查点。
// 它接收当前的查询、页面、仓库、文件、提交、总结果数和已处理数作为参数，
// 并更新检查点的相应字段以及最后更新时间。
func (m *Manager) Update(query Query, page int, repo, file, commit string, totalResults, processed int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	checkpoint := m.checkpoints[query]
	checkpoint.Query = query
	checkpoint.LastPage = page
	checkpoint.LastRepo = repo
	checkpoint.LastFile = file
	checkpoint.LastCommit = commit
	checkpoint.LastUpdated = time.Now()
	checkpoint.TotalResults = totalResults
	checkpoint.Processed = processed
	return nil
}

// GetCheckpoint 方法返回当前检查点状态的一个安全副本。
// 通过返回一个值的副本而非指针，可以防止外部代码无意中修改了真实的检查点状态，
// 从而保证了管理器内部状态的完整性和线程安全。
func (m *Manager) GetCheckpoint(query Query) *Checkpoint {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// 返回一个检查点对象的副本，以避免外部修改。
	cp := *m.checkpoints[query]
	return &cp
}

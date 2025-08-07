package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/zzjbattlefield/hajimi-go/internal/checkpoint"
	"github.com/zzjbattlefield/hajimi-go/internal/config"
	"github.com/zzjbattlefield/hajimi-go/internal/data"
	"github.com/zzjbattlefield/hajimi-go/internal/extractor"
	githubclient "github.com/zzjbattlefield/hajimi-go/internal/github"
	"github.com/zzjbattlefield/hajimi-go/internal/logger"
	"github.com/zzjbattlefield/hajimi-go/internal/validator"
)

var saveFile *os.File

func main() {
	configFile := flag.String("config", ".env", "配置文件路径")
	flag.Parse()

	// 初始化日志记录器
	logger := logger.Log

	// 如果指定了配置文件则覆盖默认配置
	if *configFile != ".env" {
		// 在实际实现中，您可能希望加载不同的配置文件
		logger.Infof("使用配置文件: %s", *configFile)
	}

	// 验证必需的配置
	if len(config.Conf.GithubTokens) == 0 {
		logger.Error("GITHUB_TOKENS 是必需的")
		os.Exit(1)
	}

	// 初始化 GitHub 客户端
	client := githubclient.NewClient(config.Conf.GithubTokens)
	logger.Info("GitHub 客户端初始化完成")

	// 初始化密钥提取器
	extractor.NewExtractor()
	logger.Info("密钥提取器初始化完成")

	// 初始化验证器
	googleValidator := validator.NewGoogleValidator()
	secretValidator := validator.NewMultiValidator(
		googleValidator,
	)
	logger.Info("密钥验证器初始化完成")

	if err := os.MkdirAll("data", 0755); err != nil {
		panic(fmt.Errorf("创建数据目录失败: %w", err))
	}
	// 设置默认查询字符串
	var querys = []string{"AIzaSy in:file"}
	if config.Conf.QueriesFile != "" {
		// 从文件读取查询
		content, err := os.ReadFile(config.Conf.QueriesFile)
		if err != nil {
			logger.Errorf("读取查询文件失败: %v", err)
		}
		querys = strings.Split(string(content), "\n")
	}

	// 初始化检查点管理器
	checkpointManager := checkpoint.NewManager(config.Conf.DataPath, querys)
	// 确保在退出时保存检查点
	defer func() {
		data.CacheData.Close()
		if err := checkpointManager.Save(); err != nil {
			logger.Errorf("保存检查点失败: %v", err)
		}
	}()
	err := scanGitHub(client, secretValidator, checkpointManager, querys)
	if err != nil {
		logger.Errorf("扫描失败: %v", err)
		// 在错误时保存检查点
		if saveErr := checkpointManager.Save(); saveErr != nil {
			logger.Errorf("在错误时保存检查点失败: %v", saveErr)
		}
		os.Exit(1)
	}

	logger.Info("扫描成功完成")
}

// processCodeResults 处理一批代码搜索结果
func processCodeResults(ctx context.Context, codeResults []*github.CodeResult, secretValidator validator.Validator, checkpointManager *checkpoint.Manager, query string, page int, totalResults int, processed *int, cp *checkpoint.Checkpoint) error {
	for _, codeResult := range codeResults {
		// 提取仓库和文件信息
		repoName := ""
		if codeResult.Repository != nil && codeResult.Repository.FullName != nil {
			repoName = *codeResult.Repository.FullName
		}
		filename := ""
		if codeResult.Name != nil {
			filename = *codeResult.Name
		}

		commit := ""
		if codeResult.SHA != nil {
			commit = *codeResult.SHA
		}

		// 如果我们已经处理过这个文件则跳过（基于检查点）
		if cp.LastRepo != "" && cp.LastFile != "" && cp.LastCommit != "" {
			if repoName == cp.LastRepo && filename == cp.LastFile && commit == cp.LastCommit {
				// 跳过此文件并继续处理下一个
				logger.Log.Infof("跳过已处理的文件: %s/%s", repoName, filename)
				continue
			}
		}

		// 从文本匹配中获取文件内容
		content := ""
		if codeResult.TextMatches != nil {
			// 连接所有文本匹配
			var fragments []string
			for _, match := range codeResult.TextMatches {
				if match.Fragment != nil {
					fragments = append(fragments, *match.Fragment)
				}
			}
			content = strings.Join(fragments, "\n")
		}

		// 提取密钥
		secrets := extractor.SecretExtractor.Extract(content, filename, repoName, commit)

		// 验证并记录找到的密钥
		for _, secret := range secrets {
			// 将密钥类型转换为 validator.SecretType
			secretType := validator.SecretType(secret.Type)
			// 验证密钥
			validationResult, err := secretValidator.Validate(ctx, secretType, secret.Value)
			if err != nil {
				logger.Log.Errorf("验证密钥失败 - 类型: %s, 值: %s, ", secret.Type, secret.Value)
				// 即使验证失败也记录密钥
				logger.Log.Infof("找到密钥 - 类型: %s, 值: %s, 文件: %s, 仓库: %s, 行: %d",
					secret.Type, secret.Value, secret.File, secret.Repo, secret.Line)
				continue
			}

			// 记录验证结果
			if validationResult.Valid {
				logger.Log.Infof("找到有效密钥 - 类型: %s, 值: %s, 文件: %s, 仓库: %s, 行: %d, 详情: %s",
					secret.Type, secret.Value, secret.File, secret.Repo, secret.Line, validationResult.Details)
				data.CacheData.Add(secret.Value)
			} else {
				logger.Log.Infof("找到无效密钥 - 类型: %s, 值: %s, 文件: %s, 仓库: %s, 行: %d, 错误: %s, 详情: %s",
					secret.Type, secret.Value, secret.File, secret.Repo, secret.Line, validationResult.ErrorCode, validationResult.Details)
			}
		}

		// 定期更新检查点
		(*processed)++
		if (*processed)%10 == 0 || *processed == totalResults { // 每处理 10 个文件更新一次检查点
			if err := checkpointManager.Update(query, page, repoName, filename, commit, totalResults, *processed); err != nil {
				logger.Log.Errorf("更新检查点失败: %v", err)
			}

			// 保存检查点
			if err := checkpointManager.Save(); err != nil {
				logger.Log.Errorf("保存检查点失败: %v", err)
			}

			logger.Log.Infof("已处理 %d 个文件，保存检查点", *processed)
		}

		// 添加小延迟以避免速率限制
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func scanGitHub(client *githubclient.Client, secretValidator validator.Validator, checkpointManager *checkpoint.Manager, querys []string) error {
	ctx := context.Background()

	for _, query := range querys {
		// 获取检查点
		cp := checkpointManager.GetCheckpoint(query)

		// 如果我们有检查点，则从上次停止的地方恢复
		startPage := 1
		if cp.LastPage > 0 {
			startPage = cp.LastPage
			logger.Log.Infof("从第 %d 页恢复扫描", startPage)
		}

		// GitHub 搜索选项
		opts := &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    startPage,
				PerPage: 30,
			},
		}
		processed := cp.Processed
		logger.Log.Infof("开始使用查询: %s 进行扫描", query)
		// 搜索代码
		var result *github.CodeSearchResult
		var resp *github.Response
		var err error
		retry := 0
		for retry < 3 {
			if result, resp, err = client.SearchCode(ctx, query, opts); err != nil {
				logger.Log.Errorf("搜索代码失败,尝试更换令牌重试: %v", err)
				client.RotateToken()
				retry++
			} else {
				break
			}
		}
		// 如果我们是从头开始则更新总结果数
		if result.Total != nil && cp.TotalResults != *result.Total {
			cp.TotalResults = *result.Total
			logger.Log.Infof("找到 %d 个代码结果", *result.Total)
		}

		// 处理结果
		if err := processCodeResults(ctx, result.CodeResults, secretValidator, checkpointManager, query, opts.Page, cp.TotalResults, &processed, cp); err != nil {
			return err
		}

		// 处理分页
		for resp != nil && resp.NextPage > 0 {
			logger.Log.Infof("处理第 %d 页", resp.NextPage)

			// 更新下一页的选项
			opts.Page = resp.NextPage

			// 搜索下一页
			result, resp, err = client.SearchCode(ctx, query, opts)
			if err != nil {
				// 尝试轮换令牌并重试一次
				client.RotateToken()
				result, resp, err = client.SearchCode(ctx, query, opts)
				if err != nil {
					// 在返回错误前保存检查点
					if saveErr := checkpointManager.Save(); saveErr != nil {
						logger.Log.Errorf("在错误时保存检查点失败: %v", saveErr)
					}
					return fmt.Errorf("搜索代码失败: %w", err)
				}
			}

			// 处理下一页的结果
			if err := processCodeResults(ctx, result.CodeResults, secretValidator, checkpointManager, query, opts.Page, cp.TotalResults, &processed, cp); err != nil {
				return err
			}
		}

		// 更新最终检查点
		if err := checkpointManager.Update(query, opts.Page, "", "", "", cp.TotalResults, processed); err != nil {
			logger.Log.Errorf("更新最终检查点失败: %v", err)
		}
	}

	return nil
}

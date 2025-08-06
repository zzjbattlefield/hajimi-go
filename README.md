# Hajimi GO - GitHub 密钥扫描器

Hajimi GO 用于在 GitHub 上扫描gemini_key_api并且验证其有效性，提供断点续传功能，以确保在长时间扫描任务中不会丢失进度。

## 核心功能

- **GitHub 代码扫描**: 根据用户定义的查询或默认查询，在 GitHub 上搜索可能泄露的密钥。- **密钥提取**: 从搜索到的代码中提取多种类型的密钥（目前主要用于gemini_key_api）。
- **密钥验证**: 验证提取的密钥是否有效。- **断点续传**: 保存扫描进度，以便在任务中断后可以从上次停止的地方继续。- **令牌轮换**: 支持使用多个 GitHub 令牌，以避免因速率限制而导致扫描中断。
- **灵活配置**: 通过环境变量或 `.env` 文件轻松配置所有选项。

## 配置

您可以通过创建 `.env` 文件或设置环境变量来配置 Hajimi GO。

| 环境变量 | `.env` 文件中的键 | 描述 | 默认值 |
| --- | --- | --- | --- |
| `GITHUB_TOKENS` | `GITHUB_TOKENS` | 用于向 GitHub API 进行身份验证的 GitHub 个人访问令牌，多个令牌用逗号分隔。**这是必需的配置**。 | || `PROXY` | `PROXY` | 用于网络请求的代理服务器地址。 | |
| `DATA_PATH` | `DATA_PATH` | 用于存储检查点文件和有效密钥的数据目录路径。 | `/app/data` |
| `DATE_RANGE_DAYS` | `DATE_RANGE_DAYS` | 扫描最近多少天内更新的代码。 | `730` || `QUERIES_FILE` | `QUERIES_FILE` | 包含要用于搜索的查询的文件路径，每行一个查询。 | `queries.txt` |


### 示例 `.env` 文件您可以复制 [` .env.example`](.env.example:1) 文件来创建自己的 `.env` 文件：

```env
# .env
GITHUB_TOKENS=your_github_token_1,your_github_token_2
PROXY=http://your_proxy_server:port
DATA_PATH=./data
QUERIES_FILE=queries.txt
```

## 如何运行

### 先决条件

- [Go](https://golang.org/) (1.18 或更高版本)
- [Git](https://git-scm.com/)

### 从源码运行

1.  **克隆仓库**:
    ```bash
    git clone https://github.com/zzjbattlefield/hajimi-go.git
    cd hajimi-go
    ```

2.  **安装依赖**:
    ```bash
    make deps
    ```

3.  **配置**:
    创建 `.env` 文件并填入您的 GitHub 令牌。

4.  **运行**:
    ```bash
    make run
    ```

### 使用 `Makefile`

`Makefile` 提供了一些有用的命令来简化开发和构建过程：

- `make build`: 构建二进制文件。
- `make test`: 运行测试。
- `make run`: 运行应用程序。
- `make clean`: 清理构建产物。
- `make install`: 安装二进制文件。
- `make build-all`: 为 Linux、Windows 和 macOS 构建二进制文件。
- `make help`: 显示所有可用的 `make` 命令。


## 许可证

本项目根据 [MIT 许可证](LICENSE)授权。
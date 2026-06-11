# Slimming

Windows C盘清理CLI工具

## 功能

- 智能风险分级扫描
- 管道式架构清理（Discover → Classify → Report → Confirm → Execute）
- 交互式确认机制
- 可配置风险等级
- 自动生成清理报告

## 清理模块

| 模块 | 风险等级 | 说明 |
|------|----------|------|
| 临时文件 | low | Windows Temp、用户 %TEMP% |
| 回收站 | low | 清空回收站内容 |
| 缩略图缓存 | low | Thumbs.db、图标缓存 |
| 系统日志 | low | Windows 日志、CBS 日志 |
| 浏览器缓存 | medium | Chrome/Edge/Firefox 缓存目录 |
| Windows Update | medium | SoftwareDistribution 下载缓存 |
| 包管理器缓存 | medium | npm cache、pip cache、go mod cache |
| Windows.old | high | 旧系统安装残留 |
| 大文件扫描 | high | 扫描 >100MB 文件 |
| 重复文件 | high | 基于内容哈希检测 |

## 安装

```bash
go build -o slimming.exe .
```

## 使用

```bash
# 扫描可清理内容
slimming scan

# 清理C盘空间（交互式确认）
slimming clean

# 预览模式
slimming clean --dry-run

# 跳过确认，清理所有风险等级
slimming clean --all

# 查看配置
slimming config
```

## 配置

配置文件位置: `~/.slimming/config.toml`

```toml
# 默认风险等级覆盖
[risk]
# browser_cache = "low"
# windows_old = "medium"

# 大文件扫描阈值
[large_files]
min_size_mb = 100
exclude_paths = [
  "C:\\Users\\*\\Downloads",
  "D:\\",
]

# 重复文件扫描
[duplicates]
min_size_mb = 1
exclude_extensions = [".exe", ".dll", ".sys"]

# 清理报告
[report]
save_dir = "~/.slimming/reports"
keep_count = 30
```

## 报告

清理完成后会自动生成JSON格式的报告，保存在 `~/.slimming/reports/` 目录下。

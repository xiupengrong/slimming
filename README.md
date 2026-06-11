# Slimming

Windows C盘清理CLI工具

## 功能

- 智能风险分级扫描
- 管道式架构清理
- 安全确认机制
- 可配置风险等级

## 安装

```bash
go build -o slimming.exe .
```

## 使用

```bash
# 扫描可清理内容
slimming scan

# 清理C盘空间
slimming clean

# 预览模式
slimming clean --dry-run

# 只清理低风险项
slimming clean --risk low
```

## 配置

配置文件位置: `~/.slimming/config.toml`

```toml
[large_files]
min_size_mb = 100

[duplicates]
min_size_mb = 1
```

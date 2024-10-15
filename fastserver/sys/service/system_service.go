package service

import (
	"errors"
	"fastgin/config"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type SystemService struct {
}

func (s *SystemService) GetSystemInformation() map[string]any {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)
	getOSInfo := func() map[string]any {
		return map[string]any{
			"os":            runtime.GOOS,
			"num_cpu":       runtime.NumCPU(),
			"compiler":      runtime.Compiler,
			"go_version":    runtime.Version(),
			"num_goroutine": runtime.NumGoroutine(),
		}
	}
	getCpuInfo := func() map[string]any {
		cpuInfo := make(map[string]any)
		cpuInfo["states"], _ = cpu.Info()
		cpuInfo["count"], _ = cpu.Counts(false)
		cpuInfo["cpu_percent"], _ = cpu.Percent(0, false)
		return cpuInfo
	}
	getMemInfo := func() map[string]any {
		v, _ := mem.VirtualMemory()
		return map[string]any{
			"total":        v.Total / MB,
			"available":    v.Available / MB,
			"used":         v.Used / MB,
			"used_percent": v.UsedPercent,
		}
	}
	getDiskInfo := func() map[string]any {
		// 获取当前程序的路径
		executable, err := os.Executable()
		if err != nil {
			return map[string]any{}
		}

		// 获取当前程序所在的磁盘分区
		partition := executable[:strings.LastIndex(executable, string(os.PathSeparator))]

		// 获取该分区的使用信息
		usage, err := disk.Usage(partition)
		if err != nil {
			return map[string]any{}
		}

		// 返回该分区的使用信息
		return map[string]any{
			"partition":    partition,
			"total":        usage.Total / GB,
			"used":         usage.Used / GB,
			"free":         usage.Free / GB,
			"used_percent": usage.UsedPercent,
		}
	}

	info := make(map[string]any)
	info["os"] = getOSInfo()
	info["cpu"] = getCpuInfo()
	info["mem"] = getMemInfo()
	info["disk"] = getDiskInfo()
	info["app_version"] = config.AppVersion
	return info
}
func (s *SystemService) Reload() error {
	if runtime.GOOS == "windows" {
		return errors.New("系统不支持")
	}
	pid := os.Getpid()
	cmd := exec.Command("kill", "-1", strconv.Itoa(pid))
	return cmd.Run()
}
func (s *SystemService) Restart() error {
	// 获取当前程序的路径
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	// 启动新的进程
	cmd := exec.Command(executable)
	err = cmd.Start()
	if err != nil {
		config.Log.Error("重启服务失败!", err)
		return err
	}
	config.Log.Info("重启服务成功,将杀死老进程!")
	// 退出当前进程
	os.Exit(0)
	return nil
}

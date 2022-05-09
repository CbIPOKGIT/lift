package nanopi

import (
	"fmt"
	"testing"
)

var CPU CPUInfo

func TestGetLoadAverage(t *testing.T) {
	fmt.Print("GetLoadAverage: ")
	fmt.Println(GetLoadAverage())
}

func TestGetCpuTemp(t *testing.T) {
	fmt.Print("GetCpuTemp: ")
	fmt.Println(GetCpuTemp())
}

func TestGetCpuInfo(t *testing.T) {
	CPU, _ = GetCpuInfo()
	fmt.Printf("%#v\n", CPU)
}

func TestCPUInfo_NumProc(t *testing.T) {
	fmt.Println(CPU.NumProc())
}

func TestCPUInfo_GetCurrentFrequency(t *testing.T) {
	fmt.Println(CPU.GetCurrentFrequency())
}

func TestCPUInfo_GetIps(t *testing.T) {
	fmt.Print("GetIps: ")
	fmt.Println(GetIps())
}

func TestCPUInfo_GetUptime(t *testing.T) {
	fmt.Println("GetUptime", GetUptime())
}

func TestCPUInfo_GetMaininfo(t *testing.T) {
	fmt.Println("GetMainInfo", GetMainInfo())
}

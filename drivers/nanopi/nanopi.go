package nanopi

import (
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

type ProcInfo struct {
	Model    string
	Features string
	Arch     string
}

type CPUInfo struct {
	Proc     []ProcInfo
	Hardware string
	Serial   string
}

// const stdIps = []string{string(net.IPv4zero)}

func GetLoadAverage() (float32, float32, float32, error) {
	var m [3]float32
	f, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return m[0], m[1], m[2], err
	}

	tmp := strings.Split(string(f), " ")

	for i := 0; i < 3; i++ {
		tmp1, err := strconv.ParseFloat(tmp[i], 32)
		if err != nil {
			return m[0], m[1], m[2], err
		}
		m[i] = float32(tmp1)
	}

	return m[0], m[1], m[2], nil
}

func GetCpuTemp() (int, error) {
	var temp int
	f, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return temp, err
	}

	t, err := strconv.Atoi(strings.Trim(string(f), "\\ \n\r\t"))
	if err != nil {
		return 0, err
	}

	return t / 1000, nil
}

/*
func GetCpuFreq() ([4]uint, error) {

}

*/

func GetCpuInfo() (CPUInfo, error) {
	cpu := CPUInfo{}
	var nCpu int
	f, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return cpu, err
	}

	tmp := strings.Split(string(f), "\n")
	for _, j := range tmp {
		tmp1 := strings.Split(j, ":")
		if len(tmp1) != 2 {
			continue
		}
		switch {
		case strings.Contains(tmp1[0], "processor"):
			cpu.Proc = append(cpu.Proc, ProcInfo{})
			nCpu = len(cpu.Proc) - 1
		case strings.Contains(tmp1[0], "model name"):
			cpu.Proc[nCpu].Model = trim(tmp1[1])
		case strings.Contains(tmp1[0], "Features"):
			cpu.Proc[nCpu].Features = trim(tmp1[1])
		case strings.Contains(tmp1[0], "CPU architecture"):
			cpu.Proc[nCpu].Arch = trim(tmp1[1])
		case strings.Contains(tmp1[0], "Hardware"):
			cpu.Hardware = trim(tmp1[1])
		case strings.Contains(tmp1[0], "Serial"):
			cpu.Serial = trim(tmp1[1])
		}
	}

	return cpu, nil
}

func (c CPUInfo) NumProc() int {
	return len(c.Proc)
}

func (c CPUInfo) GetCurrentFrequency() ([]int, error) {
	cpus := make([]int, len(c.Proc))
	prefix := "/sys/bus/cpu/devices/cpu"
	suffix := "/cpufreq/cpuinfo_cur_freq"
	for i := 0; i < len(c.Proc); i++ {
		b, e := ioutil.ReadFile(prefix + strconv.Itoa(i) + suffix)
		if e != nil {
			return cpus, e
		}
		tmp, e := strconv.Atoi(trim(string(b)))
		if e != nil {
			return cpus, e
		}
		cpus[i] = tmp / 1000
	}

	return cpus, nil
}

func trim(s string) string {
	return strings.Trim(s, " \n\r\t")
}

func GetIps() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

func GetUptime() string {
	up, _ := host.Uptime()

	now := time.Now()
	ts := int64(up)
	timeFromTS := time.Unix(now.Unix()-ts, 0)
	diff := now.Sub(timeFromTS)
	return diff.String()
}

func GetMainInfo() string {
	firstIp := ""
	totalLoad := ""
	temp, _ := GetCpuTemp()
	ips, errIp := GetIps()
	if errIp == nil {
		firstIp = ips[0]
	}
	load1, load2, _, errLoad := GetLoadAverage()
	if errLoad == nil {
		totalLoad = " load: " + strconv.Itoa(int(load1)) + ", " + strconv.Itoa(int(load2))
	}

	return "t: " + strconv.Itoa(temp) + ", ip: " + firstIp + totalLoad
}

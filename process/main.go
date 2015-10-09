package main

import (
	"flag"
    "fmt"
    "syscall"
    "os/exec"
    "unsafe"
 //   "strconv"
"time"
)

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
    dwSize ulong
    cntUsage ulong
    th32ProcessID ulong
    th32DefaultHeapID ulong_ptr
    th32ModuleID ulong
    cntThreads ulong
    th32ParentProcessID ulong
    pcPriClassBase ulong
    dwFlags ulong
    szExeFile [260]byte
}

var compareProcess string = "iexplore.exe"
var cmd string = "C:\\Program Files\\Internet Explorer\\iexplore.exe"


var readme string = `每秒检测当前运行的windows系统是否正在运行某个进程，如果没有运行，执行一个指定的命令
两个参数，第一个为进程全名，第二个为无进程时执行的完整命令，默认参数的作用是始终有ie在运行`


func main() {
	flag.Parse()
	t := time.Tick(1 * time.Second)
	if flag.NArg() == 2 {
		compareProcess = flag.Arg(0)
		cmd = flag.Arg(1)
	}
        //v := <- t.C 
        //fmt.Println(v)
        fmt.Println(readme)
        //time.Sleep(3 * time.Second)
		run := false
        for now := range t { 
                // now := <- c
				run = isRunning()
                fmt.Println("run", now, run)
if(run == false){
	c := exec.Command(cmd)
	c.Start()
}
        }   
		
		
}
func isRunning()bool{
	
    kernel32 := syscall.NewLazyDLL("kernel32.dll");
    CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot");
    pHandle,_,_ := CreateToolhelp32Snapshot.Call(uintptr(0x2),uintptr(0x0));
    if int(pHandle)==-1 {
        return false;
    }
    Process32Next := kernel32.NewProc("Process32Next");
	run := false;
    for {
        var proc PROCESSENTRY32;
        proc.dwSize = ulong(unsafe.Sizeof(proc));
        if rt,_,_ := Process32Next.Call(uintptr(pHandle),uintptr(unsafe.Pointer(&proc)));int(rt)==1 {
			pName := string(proc.szExeFile[0:]);
			if(compareProcess != "" && pName[0:len(compareProcess)] == compareProcess) {
				run = true;
			}
            //fmt.Println("ProcessName : ",pName+"sssschrome.exe","j");
            //fmt.Println("ProcessID : "+strconv.Itoa(int(proc.th32ProcessID)));
        }else{
            break;
        }
    }
    CloseHandle := kernel32.NewProc("CloseHandle");
    _,_,_ = CloseHandle.Call(pHandle);
	return run;
}
func onTime(c <-chan time.Time) {
        for now := range c { 
                // now := <- c
                fmt.Println("onTime", now)
        }   
}


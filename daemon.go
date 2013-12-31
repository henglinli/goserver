package daemon

import (
        "syscall"
        "os"
        "log"
)

func Daemon(nochdir, noclose int) int {
        var ret, ret2 uintptr
        var err uintptr
	
        darwin := syscall.OS == "darwin"
	
        // already a daemon
        if syscall.Getppid() == 1 {
                return 0
        }
	
        // fork off the parent process
        ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
        if err != 0 {
                return -1
        }
	
        // failure
        if ret2 < 0 {
                os.Exit(-1)
        }
	
        // handle exception for darwin
        if darwin && ret2 == 1 {
                ret = 0
        }
	
        // if we got a good PID, then we call exit the parent process.
        if ret > 0 {
                os.Exit(0)
        }
	
        /* Change the file mode mask */
        _ = syscall.Umask(0)
	
        // create a new SID for the child process
        s_ret, s_errno := syscall.Setsid()
        if s_errno != 0 {
                log.Printf("Error: syscall.Setsid errno: %d", s_errno)
        }
        if s_ret < 0 {
                return -1
        }
	
        if nochdir == 0 {
                os.Chdir("/")
        }
	
        if noclose == 0 {
                f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
                if e == nil {
                        fd := f.Fd()
                        syscall.Dup2(fd, os.Stdin.Fd())
                        syscall.Dup2(fd, os.Stdout.Fd())
                        syscall.Dup2(fd, os.Stderr.Fd())
                }
        }
	
        return 0
}
// usage example: daemon(0, 0)

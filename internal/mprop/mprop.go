package mprop

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PicoOrg/AndroidBox/internal/ndk"
	"github.com/PicoOrg/AndroidBox/internal/util"
)

type MProp interface {
	Set(name, value string) (err error)
}

func NewMProp(logger util.Logger, initPid int) MProp {
	return &mprop{
		logger:  logger,
		initPid: initPid,
	}
}

type mprop struct {
	logger  util.Logger
	initPid int
}

func (i *mprop) Set(name, value string) (err error) {
	realValue, err := ndk.SystemPropertyGet(name)
	if err != nil {
		i.logger.Error("get property error", util.Fields{"name": name, "error": err})
		return
	}
	if realValue == value {
		i.logger.Debug("property already set", util.Fields{"name": name, "value": value})
		return
	}
	if strings.HasPrefix(name, "ro.") {
		return i.setEnforce(name, value)
	} else {
		return i.setNormal(name, value)
	}
}

func (i *mprop) setNormal(name, value string) (err error) {
	var realValue string
	for times := 0; times <= 10; times++ {
		err = ndk.SystemPropertySet(name, value)
		if err != nil {
			i.logger.Error("system_property_set error", util.Fields{"error": err})
			return
		}
		time.Sleep(time.Millisecond)
		realValue, err = ndk.SystemPropertyGet(name)
		if err != nil {
			i.logger.Error("system_property_get error", util.Fields{"error": err})
			return
		}
		if realValue == value {
			return
		}
	}
	err = fmt.Errorf("set property or, name: %s, value: %s, real value: %s", name, value, realValue)
	i.logger.Error("set property error", util.Fields{"name": name, "value": value, "real value": realValue, "error": err})
	return
}

func (i *mprop) setEnforce(name, value string) (err error) {

	err = syscall.PtraceAttach(i.initPid)
	if err != nil {
		i.logger.Error("ptrace attach error", util.Fields{"pid": i.initPid, "error": err})
		return
	}
	defer syscall.PtraceDetach(i.initPid)

	var wstatus syscall.WaitStatus
	for {
		_, err = syscall.Wait4(i.initPid, &wstatus, 0, nil)
		if wstatus.Stopped() {
			break
		}
	}

	maps, err := os.ReadFile(fmt.Sprintf("/proc/%d/maps", i.initPid))
	if err != nil {
		i.logger.Error("read maps error", util.Fields{"pid": i.initPid, "error": err})
		return
	}

	var st uintptr = 0
	// 00008000-000cb000 r-xp 00000000 00:01 6999       /init
	for _, line := range strings.Split(string(maps), "\n") {
		if strings.Contains(line, "/init") {
			var ms, me uint64
			lineSplit := strings.Split(strings.TrimSpace(line), " ")
			if len(lineSplit) < 3 {
				i.logger.Error("invalid line", util.Fields{"line": line, "error": err})
				return
			}
			mapHex := strings.Split(strings.TrimSpace(lineSplit[0]), "-")
			if len(mapHex) != 2 {
				i.logger.Error("invalid mapHex line", util.Fields{"line": line, "error": err})
				return
			}
			ms, err = strconv.ParseUint(mapHex[0], 16, 64)
			if err != nil {
				i.logger.Error("parse ms error", util.Fields{"line": line, "error": err})
				return
			}
			me, err = strconv.ParseUint(mapHex[1], 16, 64)
			if err != nil {
				i.logger.Error("parse me error", util.Fields{"line": line, "error": err})
				return
			}

			bufferSize := me - ms
			buffer := make([]byte, bufferSize)
			var rc int
			rc, err = syscall.PtracePeekData(i.initPid, uintptr(ms), buffer)
			if err != nil {
				i.logger.Error("peek data error", util.Fields{"line": line, "error": err, "buffer size": bufferSize, "rc": rc})
				return
			}
			for pos := 0; pos < len(buffer); pos++ {
				if (buffer[pos] == 0x72 || buffer[pos] == 0x73) && buffer[pos+1] == 0x6f && buffer[pos+2] == 0x2e && buffer[pos+3] == 0x00 {
					st = uintptr(ms) + uintptr(pos)
					break
				}
			}
			if st != 0 {
				break
			}
		}
	}
	if st == 0 {
		err = fmt.Errorf("st not found")
		i.logger.Error("magic string not found, please contract the author!", util.Fields{"pid": i.initPid, "error": err})
		return
	} else {
		i.logger.Debug("st found", util.Fields{"pid": i.initPid, "st": st})
	}

	if strings.HasPrefix(name, "ro.") {
		_, err = syscall.PtracePokeData(i.initPid, st, []byte{0x73, 0x6f, 0x2e, 0x00})
		if err != nil {
			i.logger.Error("ptrace poke data error", util.Fields{"pid": i.initPid, "error": err})
			return
		}
	}

	err = syscall.PtraceCont(i.initPid, 0)
	if err != nil {
		i.logger.Error("ptrace cont error", util.Fields{"pid": i.initPid, "error": err})
		return
	}

	err = i.setNormal(name, value)
	if err != nil {
		return
	}

	err = syscall.Kill(i.initPid, syscall.SIGSTOP)
	if err != nil {
		i.logger.Error("kill error", util.Fields{"pid": i.initPid, "error": err})
		return
	}

	for {
		_, err = syscall.Wait4(i.initPid, &wstatus, 0, nil)
		if err != nil {
			i.logger.Error("wait4 error", util.Fields{"pid": i.initPid, "error": err})
			return
		}
		if wstatus.Stopped() {
			break
		}
	}

	if st != 0 {
		_, err = syscall.PtracePokeData(i.initPid, st, []byte{0x72, 0x6f, 0x2e, 0x00})
		if err != nil {
			i.logger.Error("ptrace poke data error", util.Fields{"pid": i.initPid, "error": err})
			return
		}
	}

	buff := make([]byte, 4)
	_, err = syscall.PtracePeekData(i.initPid, st, buff)
	if err != nil {
		i.logger.Error("ptrace peek data error", util.Fields{"pid": i.initPid, "error": err})
		return
	}

	if buff[0] == 0x72 && buff[1] == 0x6f && buff[2] == 0x2e && buff[3] == 0x00 {
		i.logger.Debug("mprop magic string reset success", util.Fields{"st": st, "value": buff})
	}
	return nil
}

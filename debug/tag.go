package debug

import (
	"time"
	"core/util"
	"encoding/json"
	"bytes"
	"fmt"
	"github.com/verystar/golib/color"
	"strings"
)

type DebugTagData struct {
	Key     string
	Data    interface{}
	Stack   string
	Current string
}
type DebugTag struct {
	t         time.Time
	Data      []DebugTagData
	DebugFlag string
	SavePath  string
}

func NewDebugTag(options ...func(*DebugTag)) *DebugTag {
	debug := &DebugTag{
		DebugFlag: "on",
	}

	for _, option := range options {
		option(debug)
	}

	debug.Start()
	return debug
}

func (this *DebugTag) Start() {
	if this.DebugFlag == "off" {
		return
	}
	this.t = time.Now()
}

func (this *DebugTag) Tag(key string, data ...interface{}) {
	if this.DebugFlag == "off" {
		return
	}

	st := Stack(2)

	fmt.Println(color.Blue("[Debug Tag]") + " -------------------------> " + key + " <-------------------------")
	fmt.Println(color.Green(string(st)))
	fmt.Println(data...)

	this.Data = append(this.Data, DebugTagData{
		Key:     key,
		Data:    data,
		Stack:   string(st),
		Current: time.Now().Sub(this.t).String(),
	})
}

func (this *DebugTag) GetTagData() []DebugTagData {
	return this.Data
}

func (this *DebugTag) Save(dir string, format string, prefix ...string) error {
	pre := ""
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	if this.DebugFlag == "off" {
		return nil
	}

	now := time.Now()
	s := now.Format(format)
	filename := strings.TrimRight(this.SavePath, "/") + "/" + dir + "/" + pre + "_" + s + ".log"
	//buf , err := json.Marshal(this.Data)
	buf, err := json.MarshalIndent(this.Data, "", "    ")
	if err != nil {
		return err
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("\n[%v]\n", now.String()))
	buffer.Write(buf)
	buffer.WriteString("\n\n")
	return util.WriteToFile(filename, buffer.Bytes())
}

func (this *DebugTag) SaveToSecond(dir string, prefix ...string) error {
	return this.Save(dir, "2006-01-02-15-04-05", prefix...)
}

func (this *DebugTag) SaveToMinute(dir string, prefix ...string) error {
	return this.Save(dir, "2006-01-02-15-04", prefix...)
}

func (this *DebugTag) SaveToDay(dir string, prefix ...string) error {
	return this.Save(dir, "2006-01-02-15", prefix...)
}

package xsnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yyliziqiu/gdk/xfile"
	"github.com/yyliziqiu/gdk/xlog"
)

type Snap struct {
	path   string
	data   any
	dupAge time.Duration
}

func New(path string, data interface{}) *Snap {
	return New2(path, data, 0)
}

func New2(path string, data interface{}, dupAge time.Duration) *Snap {
	return &Snap{path: path, data: data, dupAge: dupAge}
}

func New3(path string, data interface{}, interval time.Duration, n int) *Snap {
	return New2(path, data, interval*time.Duration(n)+time.Second)
}

func (s *Snap) Path() string {
	return s.path
}

func (s *Snap) Data() any {
	return s.data
}

func (s *Snap) DupAge() time.Duration {
	return s.dupAge
}

// Save 保存
func (s *Snap) Save() error {
	return Save(s.path, s.data)
}

// Load 加载
func (s *Snap) Load() error {
	return Load(s.path, s.data)
}

// Dup 生成副本
func (s *Snap) Dup() error {
	return Dup(s.path, s.data, s.dupAge)
}

// Save 保存快照
func Save(path string, data interface{}) error {
	err := xfile.MakeDir(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("mkdir snap dir %s failed [%v]", filepath.Dir(path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snap data %s failed [%v]", path, err)
	}

	temp := path + ".temp"
	err = os.WriteFile(temp, bytes, 0644)
	if err != nil {
		return fmt.Errorf("save snap data %s failed [%v]", path, err)
	}
	err = os.Rename(temp, path)
	if err != nil {
		return fmt.Errorf("rename snap file %s failed [%v]", path, err)
	}

	return nil
}

// Load 加载快照
func Load(path string, data interface{}) error {
	ok, err := xfile.Exist(path)
	if err != nil {
		return fmt.Errorf("check snap file %s failed [%v]", path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load snap file %s failed [%v]", path, err)
	}
	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, data)
}

// Dup 生成快照副本
func Dup(path string, data any, age time.Duration) error {
	if age <= 0 {
		return errors.New("duplication age must be greater than zero")
	}

	// 生成快照存储目录
	dir := filepath.Join(filepath.Dir(path), filepath.Base(path)+"-dup")

	// 清理过期快找副本
	_ = filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			if !os.IsNotExist(err) {
				xlog.Errorf("clear snap duplications failed, error: %v.", err)
			}
			return nil
		}
		if info.IsDir() || info.ModTime().After(time.Now().Add(-age)) {
			return nil
		}
		return os.Remove(file)
	})

	// 保存最新快照
	return Save(filepath.Join(dir, time.Now().Format("20060102150405")), data)
}

package xsnap

import (
	"context"
	"reflect"
	"time"

	"github.com/yyliziqiu/gdk/xlog"
	"github.com/yyliziqiu/gdk/xtime"
)

type Persistent interface {
	Load() error
	Save(exit bool) error
}

type PersistConfig struct {
	Name     string
	Interval time.Duration
}

type GetPersistConfig interface {
	PersistConfig() PersistConfig
}

func Persist(ctx context.Context, pts []Persistent) error {
	// load
	err := persistLoad(pts)
	if err != nil {
		return err
	}

	// save
	for _, pt := range pts {
		go runPersistSave(ctx, pt)
	}

	return nil
}

func persistLoad(pts []Persistent) error {
	t1 := xtime.NewTimer()
	for _, pt := range pts {
		t2 := xtime.NewTimer()
		if err := pt.Load(); err != nil {
			xlog.Errorf("Load snap failed, name: %s, error: %v.", persistentName(pt), err)
			return err
		}
		xlog.Infof("Load snap succeed, name: %s, cost: %s.", persistentName(pt), t2.Stops())
	}
	xlog.Infof("Load all snaps compeleted, cost: %s.", t1.Stops())

	return nil
}

func persistentName(pt Persistent) string {
	pc := persistentConfig(pt)
	if pc.Name != "" {
		return pc.Name
	}
	typ := reflect.TypeOf(pt)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ.Name()
}

func persistentConfig(pt any) PersistConfig {
	if c, ok := pt.(GetPersistConfig); ok {
		return c.PersistConfig()
	}
	return PersistConfig{}
}

func runPersistSave(ctx context.Context, pt Persistent) {
	pc := persistentConfig(pt)
	if pc.Interval <= 0 {
		<-ctx.Done()
		_ = persistSave(pt, true)
		return
	}

	ticker := time.NewTicker(pc.Interval)
	for {
		select {
		case <-ticker.C:
			_ = persistSave(pt, false)
		case <-ctx.Done():
			_ = persistSave(pt, true)
			return
		}
	}
}

func persistSave(pt Persistent, exit bool) error {
	timer := xtime.NewTimer()

	err := pt.Save(exit)
	if err != nil {
		xlog.Errorf("Save snap failed, name: %s, error: %v.", persistentName(pt), err)
	} else {
		xlog.Infof("Save snap succeed, name: %s, cost: %s.", persistentName(pt), timer.Stops())
	}

	return err
}

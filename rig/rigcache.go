package rig

import (
	"fmt"
	"time"

	"github.com/cgrates/ltcache"

	"github.com/dh1tw/goHamlib"
)

type RigCache struct {
	Rig   *goHamlib.Rig
	cache *ltcache.Cache
}

func NewRigCache(r *goHamlib.Rig, d time.Duration) *RigCache {
	return &RigCache{
		Rig:   r,
		cache: ltcache.NewCache(25000, d, true, nil),
	}
}

type getmode struct {
	mode     goHamlib.Mode
	pb_width int
}

func (rig *RigCache) GetMode(vfo goHamlib.VFOType) (mode goHamlib.Mode, pb_width int, err error) {
	id := fmt.Sprintf("getmode-%v", vfo)
	v, exists := rig.cache.Get(id)
	if exists {
		gm := v.(*getmode)
		return gm.mode, gm.pb_width, nil
	}
	m, p, err := rig.Rig.GetMode(vfo)
	if err != nil {
		return m, p, err
	}
	rig.cache.Set(id, &getmode{m, p}, nil)
	return m, p, err
}

type getfreq struct {
	freq float64
}

func (rig *RigCache) GetFreq(vfo goHamlib.VFOType) (freq float64, err error) {
	id := fmt.Sprintf("getfreq-%v", vfo)
	v, exists := rig.cache.Get(id)
	if exists {
		gm := v.(*getfreq)
		return gm.freq, nil
	}
	f, err := rig.Rig.GetFreq(vfo)
	if err != nil {
		return f, err
	}
	rig.cache.Set(id, &getfreq{f}, nil)
	return f, err
}

func (rig *RigCache) SetFreq(vfo goHamlib.VFOType, freq float64) error {
	if err := rig.Rig.SetFreq(vfo, freq); err != nil {
		return err
	}
	id := fmt.Sprintf("getfreq-%v", vfo)
	rig.cache.Set(id, &getfreq{freq}, nil)
	return nil
}

func (rig *RigCache) SetMode(vfo goHamlib.VFOType, mode goHamlib.Mode, pbWidth int) error {
	if err := rig.Rig.SetMode(vfo, mode, pbWidth); err != nil {
		return err
	}

	id := fmt.Sprintf("getmode-%v", vfo)
	rig.cache.Set(id, &getmode{mode, pbWidth}, nil)
	return nil
}

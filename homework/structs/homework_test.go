package main

import (
	"fmt"
	"log"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	// head
	goldSz       = 31
	manaSz       = 10
	healthSz     = 10
	respectSz    = 4
	strengthSz   = 4
	experienceSz = 4

	goldOffset       = 0
	manaOffset       = goldOffset + goldSz
	healthOffset     = manaOffset + manaSz
	respectOffset    = healthOffset + healthSz
	strengthOffset   = respectOffset + respectSz
	experienceOffset = strengthOffset + strengthSz

	// tail
	levelSz  = 4
	houseSz  = 1
	gunSz    = 1
	familySz = 1
	typeSz   = 2

	levelOffset  = 0
	houseOffset  = levelOffset + levelSz
	gunOffset    = houseOffset + houseSz
	familyOffset = gunOffset + gunSz
	typeOffset   = familyOffset + familySz
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		var buf [42]byte
		copy(buf[:], name)
		person.name = buf
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.ox = int32(x)
		person.oy = int32(y)
		person.oz = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		if gold < 0 || gold > math.MaxInt32 {
			log.Printf("wrong amount of gold")
			return
		}

		var mask uint64 = (1<<goldSz - 1) << goldOffset // 000..0001111.111 (31 times "one")

		person.packHead &^= mask // reset for gold
		person.packHead |= uint64(gold) >> goldOffset
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 || mana > 1_000 {
			log.Printf("wrong amount of mana")
			return
		}

		var mask uint64 = (1<<manaSz - 1) << manaOffset

		person.packHead &^= mask
		person.packHead |= uint64(mana) << manaOffset
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 || health > 1_000 {
			log.Printf("wrong amount of mana")
			return
		}

		var mask uint64 = (1<<healthSz - 1) << healthOffset

		person.packHead &^= mask
		person.packHead |= uint64(health) << healthOffset
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 || respect > 10 {
			log.Printf("wrong amount of mana")
			return
		}

		var mask uint64 = (1<<respectSz - 1) << respectOffset

		person.packHead &^= mask
		person.packHead |= uint64(respect) << respectOffset
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 || strength > 10 {
			log.Printf("wrong amount of mana")
			return
		}

		var mask uint64 = (1<<strengthSz - 1) << strengthOffset

		person.packHead &^= mask
		person.packHead |= uint64(strength) << strengthOffset
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 || experience > 15 {
			log.Printf("wrong amount of experience")
			return
		}

		var mask uint64 = (1<<experienceSz - 1) << experienceOffset

		person.packHead &^= mask
		person.packHead |= uint64(experience) << experienceOffset
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level < 0 || level > 10 {
			log.Printf("wrong level value")
			return
		}

		var mask uint16 = (1<<levelSz - 1) << levelOffset

		person.packTail &^= mask
		person.packTail |= uint16(level) << levelOffset

		fmt.Printf("%d\n", level)
		fmt.Printf("%b\n", mask)
		fmt.Printf("%b\n", uint64(level)<<levelOffset)
		fmt.Printf("%b\n", person.packTail)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		var mask uint16 = (1<<houseSz - 1) << houseOffset
		person.packTail &^= mask
		person.packTail |= 1 << houseOffset
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		var mask uint16 = (1<<gunSz - 1) << gunOffset
		person.packTail &^= mask
		person.packTail |= 1 << gunOffset
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		var mask uint16 = (1<<familySz - 1) << familyOffset
		person.packTail &^= mask
		person.packTail |= 1 << familyOffset
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType < 0 || personType > 3 {
			log.Printf("wrong type value")
			return
		}

		var mask uint16 = (1<<typeSz - 1) << typeOffset

		person.packTail &^= mask
		person.packTail |= uint16(personType) << typeOffset
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	packHead uint64   // 0-63    |(63)experience(60)|(59)strength(56)|(55)respect(52)|(51)health(42)|(41)mana(32)|(31)gold(0)|
	ox       int32    // 64-95
	oy       int32    // 96-127
	oz       int32    // 128-159
	packTail uint16   // 160-175 |(8)type(7)|family(6)|gun(5)|house(4)|(3)level(0)|
	name     [42]byte // 176-511
	// 512 / 8 = 64
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{}
	for _, opt := range options {
		opt(&gp)
	}
	return gp
}

func (p *GamePerson) Name() string {
	for i, b := range p.name {
		if b == 0 {
			return string(p.name[:i])
		}
	}
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.ox)
}

func (p *GamePerson) Y() int {
	return int(p.oy)
}

func (p *GamePerson) Z() int {
	return int(p.oz)
}

func (p *GamePerson) Gold() int {
	const mask = (1 << goldSz) - 1
	return int((p.packHead >> goldOffset) & uint64(mask))
}

func (p *GamePerson) Mana() int {
	const mask = (1 << manaSz) - 1
	return int((p.packHead >> manaOffset) & uint64(mask))
}

func (p *GamePerson) Health() int {
	const mask = (1 << healthSz) - 1
	return int((p.packHead >> healthOffset) & uint64(mask))
}

func (p *GamePerson) Respect() int {
	const mask = (1 << respectSz) - 1
	return int((p.packHead >> respectOffset) & uint64(mask))
}

func (p *GamePerson) Strength() int {
	const mask = (1 << strengthSz) - 1
	return int((p.packHead >> strengthOffset) & uint64(mask))
}

func (p *GamePerson) Experience() int {
	const mask = (1 << experienceSz) - 1
	return int((p.packHead >> experienceOffset) & uint64(mask))
}

func (p *GamePerson) Level() int {
	const mask = (1 << levelSz) - 1
	return int((p.packTail >> levelOffset) & uint16(mask))
}

func (p *GamePerson) Type() int {
	const mask = (1 << typeSz) - 1
	return int((p.packTail >> typeOffset) & uint16(mask))
}
func (p *GamePerson) HasHouse() bool {
	mask := uint16(1) << houseOffset
	return (p.packTail & mask) != 0
}

func (p *GamePerson) HasGun() bool {
	mask := uint16(1) << gunOffset
	return (p.packTail & mask) != 0
}

func (p *GamePerson) HasFamily() bool {
	mask := uint16(1) << familyOffset
	return (p.packTail & mask) != 0
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	// manaType
	manaMask        = 0x03FF
	personTypeMask  = 0xFC00
	personTypeShift = 10

	// healthHouseGunFamily
	healthMask  = 0x07FF
	healthShift = 3
	houseBit    = 0x0001
	gunBit      = 0x0002
	familyBit   = 0x0004

	// levelExpirienceStrengthRespect
	respectMask     = 0x000F
	respectShift    = 12
	strengthMask    = 0x000F
	strengthShift   = 8
	experienceMask  = 0x000F
	experienceShift = 4
	levelMask       = 0x000F
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(person *GamePerson) {
		copy(person.name[:], name)
		for i := len(name); i < len(person.name); i++ {
			person.name[i] = 0
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaType |= int16(mana)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaType |= int16(personType) << personTypeShift
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunFamily |= int16(health) << healthShift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunFamily |= houseBit
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunFamily |= gunBit
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.healthHouseGunFamily |= familyBit
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExpirienceStrengthRespect |= uint16(respect) << respectShift
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExpirienceStrengthRespect |= uint16(strength) << strengthShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExpirienceStrengthRespect |= uint16(experience) << experienceShift
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.levelExpirienceStrengthRespect |= uint16(level)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z                        int32
	gold                           int32
	name                           [42]byte
	levelExpirienceStrengthRespect uint16
	manaType                       int16
	healthHouseGunFamily           int16
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	length := 0
	for i, b := range p.name {
		if b == 0 {
			break
		}
		length = i + 1
	}
	return unsafe.String((*byte)(unsafe.Pointer(&p.name[0])), length)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaType & manaMask)
}

func (p *GamePerson) Type() int {
	return int(p.manaType >> personTypeShift)
}

func (p *GamePerson) Health() int {
	return int(p.healthHouseGunFamily>>healthShift) & healthMask
}

func (p *GamePerson) HasHouse() bool {
	return p.healthHouseGunFamily&houseBit != 0
}

func (p *GamePerson) HasGun() bool {
	return p.healthHouseGunFamily&gunBit != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.healthHouseGunFamily&familyBit != 0
}

func (p *GamePerson) Respect() int {
	return int((p.levelExpirienceStrengthRespect >> respectShift) & respectMask)
}

func (p *GamePerson) Strength() int {
	return int((p.levelExpirienceStrengthRespect >> strengthShift) & strengthMask)
}

func (p *GamePerson) Experience() int {
	return int((p.levelExpirienceStrengthRespect >> experienceShift) & experienceMask)
}

func (p *GamePerson) Level() int {
	return int((p.levelExpirienceStrengthRespect) & levelMask)
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
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(person *GamePerson) {
		if len(name) > 0 {
			person.attributes4p |= int32(name[0]) << 24
		}
		if len(name) > 1 {
			person.attributes3p |= int32(name[1]) << 20
		}
		if len(name) > 2 {
			copy(person.name[:], name[2:])
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
		person.attributes3p |= int32(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes3p |= int32(health) << 10
	}
}

func (p *GamePerson) HasHouse() bool {
	return p.attributes3p&HasHouse != 0
}

func (p *GamePerson) HasGun() bool {
	return p.attributes3p&HasGun != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.attributes3p&HasFamily != 0
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes4p |= int32(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes4p |= int32(strength) << 4
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes4p |= int32(experience) << 8
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes4p |= int32(level) << 12
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes3p |= HasHouse
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes3p |= HasGun
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes3p |= HasFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes4p |= int32(personType) << 24
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	HasHouse  int32 = 1 << 28
	HasGun    int32 = 1 << 29
	HasFamily int32 = 1 << 30
)

type GamePerson struct {
	x, y, z      int32
	gold         int32
	attributes4p int32
	attributes3p int32
	name         [40]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	name := make([]byte, 0, 42)

	firstByte := (p.attributes4p >> 24)
	secondByte := byte(p.attributes3p >> 20)

	if firstByte != 0 {
		name = append(name, byte(firstByte))
	}
	if secondByte != 0 {
		name = append(name, secondByte)
	}

	for _, b := range p.name {
		if b == 0 {
			break
		}
		name = append(name, b)
	}

	return string(name)
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
	return int((p.attributes3p)) & 0x03FF
}

func (p *GamePerson) Health() int {
	return int((p.attributes3p)>>10) & 0x03FF
}

func (p *GamePerson) Respect() int {
	return int((p.attributes4p) & 0x0F)
}

func (p *GamePerson) Strength() int {
	return int((p.attributes4p >> 4) & 0x0F)
}

func (p *GamePerson) Experience() int {
	return int((p.attributes4p >> 8) & 0x0F)
}

func (p *GamePerson) Level() int {
	return int((p.attributes4p >> 12) & 0x0F)
}

func (p *GamePerson) Type() int {
	return int((p.attributes4p >> 16) & 0x0F)
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

package tests

import (
	"fmt"
	"testing"
)

type PS5 interface {
	StartGPUEngine()
	GetPrice() int64
}

type PS5WithCD struct {
}

func (cd *PS5WithCD) StartGPUEngine() {
	fmt.Println("StartGPUEngine")
}

func (ps *PS5WithCD) GetPrice() int64 {
	return 5000
}

type PS5WithDigital struct {
}

func (dt *PS5WithDigital) StartGPUEngine() {
	fmt.Println("PS5WithDigital StartGPUEngine")
}

func (dt *PS5WithDigital) GetPrice() int64 {
	return 3600
}

type PS5MachinePlus struct {
	ps5 PS5
}

func (plus *PS5MachinePlus) SetMachine(ps5 PS5) {
	plus.ps5 = ps5
}

func (plus *PS5MachinePlus) StartGPUEngine() {
	plus.ps5.StartGPUEngine()
	fmt.Println("性能升级 plus")
}

func (plus *PS5MachinePlus) GetPrice() int64 {
	return plus.ps5.GetPrice() + 500
}

type Ps5WithTopicColor struct {
	ps5 PS5
}

func (color *Ps5WithTopicColor) SetPs5WithTopicColor(ps5 PS5) {
	color.ps5 = ps5
}

func (color *Ps5WithTopicColor) StartGPUEngine() {
	color.ps5.StartGPUEngine()
	fmt.Println("定制主题颜色")
}

func (color *Ps5WithTopicColor) GetPrice() int64 {
	return color.ps5.GetPrice() + 200
}

func TestDecorator(t *testing.T) {

	plus:=&PS5MachinePlus{}
	cd:=&PS5WithCD{}
	plus.SetMachine(cd)
	plus.StartGPUEngine()
	plus.GetPrice()

	colorPlus:=Ps5WithTopicColor{}
	colorPlus.SetPs5WithTopicColor(plus)
	colorPlus.StartGPUEngine()
	colorPlus.GetPrice()

}
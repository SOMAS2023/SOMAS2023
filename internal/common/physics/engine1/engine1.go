package engine1

import (
	//"fmt"
	utils "SOMAS2023/internal/common/utils"
	"fmt"
	"math"
	//baseAgent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	//"github.com/google/uuid"
)

func Test(){
	fmt.Println("7")
}

// type ICounterAgent interface {
// 	baseagent.IAgent[ICounterAgent]
// 	DoCount()
// 	GetCount()
// }
type XY_map struct{
	*utils.Coordinates
}

type SingleAgent struct{
	*utils.Forces
	Mass float64
}
type MegaBike struct{
	Agents []SingleAgent
	Acceleration float64
	Velocity float64
	Orientation float64
}
type Motion interface{
	Add_Agent()
	Calculate_Orientation()
	Calculate_ResultantAcceleration()
	Calculate_Velocity()
	Update_loc()
}
func (mb *MegaBike) Add_Agent(agent SingleAgent){
	mb.Agents = append(mb.Agents, agent)
}

func (mb *MegaBike)Calculate_Orientation(){
	if len(mb.Agents)==0{
		return
	}
	Total_turning:=0.0
	for _,agent := range mb.Agents{
		Total_turning+=float64(agent.Turning)
	}
	Average_turning:=Total_turning/float64(len(mb.Agents))
	mb.Orientation+=(Average_turning)
}
func (mb *MegaBike) Calculate_ResultantAcceleration(){
	force_map:=4.0
	if len(mb.Agents)==0{
		return
	}
	Total_pedal:=0.0
	Total_brake:=0.0
	Total_mass:=0.0
	for _,agent := range mb.Agents{
		Total_mass+=float64(agent.Mass)
		if agent.Pedal !=0{
			Total_pedal+=float64(agent.Pedal)
		}else{
			Total_brake+=float64(agent.Brake)
		}
	}
	F:=force_map*(float64(Total_pedal)-float64(Total_brake))
	mb.Acceleration= (F/float64(Total_mass))
}
func (mb *MegaBike) Calculate_Velocity(dt float64){
	if mb.Velocity+mb.Acceleration*dt<0{
		mb.Velocity=0
	}else{
		mb.Velocity+=mb.Acceleration*dt
	}
}
func (Map *XY_map) Update_loc(mb *MegaBike){
	Map.X+=mb.Velocity*float64(math.Cos(float64(math.Pi*mb.Orientation)))
	Map.Y+=mb.Velocity*float64(math.Sin(float64(math.Pi*mb.Orientation)))
}



package engine2
import (
	//"fmt"
	"math"
	utils "SOMAS2023/internal/common/utils"
	//baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
)



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
	Mass float32
}
type MegaBike struct{
	Agents []SingleAgent
	Acceleration float64
	Velocity float64
	Orientation float64
}

type Motion interface{
	Add_Agent()
	Vector_addition()
	Calculate_Velocity()
	Update_loc()
}

func (mb *MegaBike) Add_Agent(agent SingleAgent){
	mb.Agents = append(mb.Agents, agent)
}
func (mb *MegaBike)Force_decomposation(agent SingleAgent)(f_x float64, f_y float64){
	f_x=agent.Pedal*float64(math.Cos(float64(math.Pi*(agent.Turning+mb.Orientation))))
	f_y=agent.Pedal*float64(math.Sin(float64(math.Pi*(agent.Turning+mb.Orientation))))
	return
}
func (mb *MegaBike)Vector_addition(){
	if len(mb.Agents)==0{
		return
	}
	F_X := 0.0
	F_Y := 0.0
	Mass := 0.0
	Brake := 0.0
	for _,agent := range mb.Agents{
		f_x,f_y:=mb.Force_decomposation(agent)
		F_X+=f_x
		F_Y+=f_y
		Brake+=agent.Brake
		Mass+=float64(agent.Mass)
	}
	F:=math.Sqrt(F_X*F_X+F_Y*F_Y)
	mb.Orientation = math.Atan2(F_Y,F_X)
	mb.Acceleration = (F-Brake)/Mass
}

func (mb *MegaBike) Calculate_Velocity(dt float64){
	if mb.Velocity+mb.Acceleration*dt<0{
		mb.Velocity=0
	}else{
		mb.Velocity+=mb.Acceleration*dt
	}
}
func (Map *XY_map) Update_loc(mb *MegaBike){
	Map.X+=float64(mb.Velocity*float64(math.Cos(math.Pi*mb.Orientation)))
	Map.Y+=float64(mb.Velocity*float64(math.Sin(math.Pi*mb.Orientation)))
}

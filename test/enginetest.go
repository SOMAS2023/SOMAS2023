package test
import (
	//"SOMAS2023/internal/server"
	physics "SOMAS2023/internal/common/physics"
	utils "SOMAS2023/internal/common/utils"
	"fmt"
)

func TestEngine(){
	fmt.Println("Test for physical engine")
	var preState utils.PhysicalState
	preState.Position.X = 1
	preState.Position.Y = 2
	preState.Velocity = 0
	preState.Acceleration = 0
	preState.Mass = 1
	var targetCoordinate utils.Coordinates
	targetCoordinate.X = 5
	targetCoordinate.Y = 15

	var orientation float64
	var l2distance float64
	var thresdistance float64 = 0.5
	var pedalforce float64 = 2
	var newState utils.PhysicalState
	var count int =0

	orientation = physics.ComputeOrientation(preState.Position,targetCoordinate)
	l2distance = physics.ComputeDistance(preState.Position,targetCoordinate)
	fmt.Println("Orientation is ", 180*orientation)
	fmt.Println("L2 Distance is ",l2distance)
	format := "Position: (%.2f, %.2f), Acceleration: %.2f, Velocity: %.2f, Mass: %.2f\n"

	for l2distance>thresdistance{
		if count>2{
			pedalforce = pedalforce-0.1
			count = 0
		}
		orientation = physics.ComputeOrientation(preState.Position,targetCoordinate)
		newState = physics.GenerateNewState(preState,pedalforce,orientation)
		fmt.Printf(format,
			newState.Position.X, newState.Position.Y,
			newState.Acceleration, newState.Velocity, newState.Mass)
		l2distance = physics.ComputeDistance(newState.Position,targetCoordinate)
		preState=newState
		if newState.Velocity==0{
			count = count+1
		}	
	}
	print("Successfully Reach the target")
}
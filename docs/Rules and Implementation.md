# Rules and Implementation Details

## Agent Goal / Winning Condition
The winning agent must be alive and have the highest points (obtained via lootboxes) at the end of the game.


## Bikers and MultiBike Forces

1. **Agent Parameters:**
   - Each agent has three parameters: Pedaling, Braking, and Turning forces.

2. **MegaBike Physics Parameters:**
   - MegaBike parameters for physics engine include Velocity and Orientation.

3. **Orientation Value:**
   - The current value of orientation for MegaBike is referred to as the offset.

4. **Turning Angle Calculation:**
   - The turning angle depends on the TurningDecision.SteeringForce and TurningDecision.SteerBike. If agents do not want to steer, they must set their TurningDecision.SteerBike to false and their steering will not have an impact on the direction of the bike. If an agent does want to steer, they must submit is a force from -1 to 1 which maps to -180° to 180°. This will then be summed with all other agents on the bike who have set TurningDecision.SteerBike to true and averaged to output a new orientation for the bike.

5. **Orientation Update:**
   - The updated Orientation is calculated by adding the Turning Angle to the Offset: `Offset = Offset + Turning Angle`.

6. **Post-Turning Forces Application:**
   - After turning, all the pedaling force and braking force will be applied in the direction of the updated orientation.

7. **Velocity Constraint:**
   - The Velocity will not drop below zero; hence, the bike will not move backwards.
   - The Velocity has a maximum value that it cannot exceed

8. **Drag Force**
   - There is a drag force that is propotional to Velocity squared.

<img src="../docs/Images/MultibikeForceOrientation.png" alt="MultiBike Force and Orientation Diagram" width="500"/> 

## Lootbox Collision
When a Megabike collides with a lootbox:
   1. All agents on the bike receive the same eneregy, irrespective of the lootbox colour.
   2. Agents of the same colour as the lootbox will receive a set number of points each.
   3. If more than one bike colides with a lootbox during one epoch, the energy will be split between the bikes equally.

## Awdi Collision
An Awdi targets the slowest bike. When an Awdi collides with a lootbox:
   1. All agents on the bike die.

## Physics Boundaries
There is no physical boundary, however lootboxes only spawn in a set area of the map. 
Therefore there is no incentive to go further off the map, but if you do you want to you will not be penalized.

## Resource Allocation Voting
- Each agent votes by passing in an array which contains the distribution of your vote for each agent (including themselves),
 normalized to one. This function takes in this array from each agent, sums up the votes for each agent and normalises the array to one. 
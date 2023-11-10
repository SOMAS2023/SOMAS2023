# AUDI
This is a brief intro of package 'Audi'.

Considering Audi is an agent, it might have complex function in the future. So instead of a go file, a package was reserved.
## Component
#### Audi
* Audi act as an agent.
* The bike of Audi only has Audi on it.
* When Audi drive, the bike of Audi moving in constant speed. To implement this, Audi will only pedal with 1 unit force to get bike moving, or brake with 1 unit force when Audi decide to stop.
* Audi will seek target with a constant logic. Currently, Audi only seek for stopped MegaBike.
## To be continued
* How Audi find its own position and How Audi get states of all existing MegaBike is waiting to be implemented, due to the unimplemented func GetGameState()
* How Audi get into his AudiBike with no one else on it is undecided yet.
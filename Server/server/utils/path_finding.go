package utils

import (
	"slices"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Node in Minecraft represents 1 block with extra information.
type Node struct {
	// hCost represents an estimated cost of reaching the goal from a given node. In this case, we use the distance to calculate the cost.
	hCost float64

	// Parent represents the parent Node to which it can be traced back to.
	Parent *Node

	// XYZ represents the position of the Node
	XYZ mgl64.Vec3
}

// gCost represents an estimated cost of reaching a given node from the start. In this case, we use the distance to calculate the cost.
func (n *Node) gCost() float64 {
	return n.gCostRecursive(n)
}

// gCostRecursive is a recursive helper function for gCost
func (n *Node) gCostRecursive(parent *Node) float64 {
	if parent == nil {
		return 1
	}
	return Distance(n.XYZ, parent.XYZ) + n.gCostRecursive(parent.Parent)
}

// fCost is the sum of gCost & hCost. It is used to check the next best node to explore. It has been multiplied by 1000 to avoid nodes with the same numbers.
func (n *Node) fCost() float64 {
	return (n.gCost() + n.hCost) * 10000000
}

// AStar represents the algorithm responsible for finding the best path possible.
type AStar struct {
	// TX represents the world in which the path finding will take place.
	TX *world.Tx

	openNodes   []*Node
	closedNodes []*Node
}

// lowestFCostNode finds the node in openNodes with the lowest F cost.
func (s AStar) lowestFCostNode() *Node {
	var minNode *Node
	for _, n := range s.openNodes {
		if minNode == nil || n.fCost() < minNode.fCost() {
			minNode = n
		}
	}
	return minNode
}

// Start To commence the exploration for the best path. It returns the Node to which you can trace the path using that Node's previous Parents.
// e.g.:
//
//	go func() {
//		pl.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
//			current, timedOut := entity.AStar{TX: tx}.Start(u.WTData.Volume.Pos1.Vec3(), u.WTData.Volume.Pos2.Vec3(), entity.DefaultAlgorithmSettings())
//			if timedOut {
//				pl.Message("Timed out!")
//			}
//			for current != nil {
//	         	//You can use 'current.XYZ' here to do whatever you want.
//				tx.SetBlock(cube.PosFromVec3(current.XYZ), block.Glass{}, nil)
//				current = current.Parent
//			}
//		})
//	}()
//
// NOTE #1: Make sure you use this function inside a Go routine as it takes time to process.
// NODE #2: The second return boolean value is when the process times out; this means it reached maximum checking distance OR maximum time.
//

func (s AStar) Start(start mgl64.Vec3, end mgl64.Vec3, settings AlgorithmSettings) (*Node, bool) {
	s.openNodes = append(s.openNodes, &Node{hCost: Distance(start, end), XYZ: start})

	startTime := time.Now()
	for s.lowestFCostNode() != nil {
		current := s.lowestFCostNode()
		s.openNodes = Filter[*Node](s.openNodes, func(n1 *Node) bool {
			return n1.XYZ != current.XYZ
		})
		s.closedNodes = append(s.closedNodes, current)

		if current.XYZ == end {
			return current, false
		}

		i := 0
		c := current
		for c != nil {
			c = c.Parent
			i++
		}
		if i > settings.MaxDistanceToCheck {
			continue
		}

		neighborXYZs := []mgl64.Vec3{
			current.XYZ.Add(mgl64.Vec3{1, 0, 0}),
			current.XYZ.Add(mgl64.Vec3{-1, 0, 0}),
			current.XYZ.Add(mgl64.Vec3{0, 1, 0}),
			current.XYZ.Add(mgl64.Vec3{0, -1, 0}),
			current.XYZ.Add(mgl64.Vec3{0, 0, 1}),
			current.XYZ.Add(mgl64.Vec3{0, 0, -1}),
		}

		for _, xyz := range neighborXYZs {
			var oldNode *Node
			for _, n := range s.openNodes {
				if n.XYZ == xyz {
					oldNode = n
					break
				}
			}

			newNode := &Node{hCost: Distance(end, xyz), XYZ: xyz}
			if oldNode != nil {
				newNode.Parent = oldNode.Parent
			}

			if settings.Validator(xyz, s.TX) && !slices.ContainsFunc[[]*Node, *Node](s.closedNodes, func(node *Node) bool {
				return newNode.XYZ == node.XYZ
			}) {
				if oldNode == nil || newNode.gCost() < oldNode.gCost() {
					newNode.Parent = current
					if oldNode == nil {
						s.openNodes = append(s.openNodes, newNode)
					}
				}
			}
		}

		if time.Now().Sub(startTime) > settings.MaxTimeToCheck {
			return current, true
		}
	}
	return nil, true
}

// AlgorithmSettings provides several options to customize your path finding.
type AlgorithmSettings struct {
	// MaxDistanceToCheck aborts the process if all paths have been explored to the maximum distance.
	MaxDistanceToCheck int
	// MaxTimeToCheck aborts process once it reaches a certain duration.
	MaxTimeToCheck time.Duration
	// Validator To specify the properties of the Nodes that can be traversed (compatible as a path)
	Validator func(v mgl64.Vec3, w *world.Tx) bool
}

// DefaultAlgorithmSettings is an example of how it should be done.
func DefaultAlgorithmSettings() AlgorithmSettings {
	return AlgorithmSettings{
		MaxDistanceToCheck: 50,
		MaxTimeToCheck:     5 * time.Second,
		Validator: func(v mgl64.Vec3, tx *world.Tx) bool {
			_, ok1 := tx.Block(cube.PosFromVec3(v)).Model().(model.Empty)
			_, ok2 := tx.Block(cube.PosFromVec3(v.Add(mgl64.Vec3{0, -1, 0}))).Model().(model.Empty)
			return ok1 && !ok2
		},
	}
}

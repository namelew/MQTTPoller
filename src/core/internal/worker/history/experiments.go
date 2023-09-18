package history

import (
	"os"

	"github.com/namelew/mqtt-poller/src/core/packages/utils"
)

type OngoingExperiments struct {
	root *OngoingExperiment
}

type OngoingExperiment struct {
	Id       int64
	Finished bool
	Proc     *os.Process
	height   int
	right    *OngoingExperiment
	left     *OngoingExperiment
}

func CreateRegister(id int64, proc *os.Process) OngoingExperiment {
	return OngoingExperiment{id, false, proc, 1, nil, nil}
}

func (tree *OngoingExperiments) Add(node *OngoingExperiment) {
	if tree.root == nil {
		tree.root = node
	} else {
		tree.root.add(node)
	}
	tree.root = tree.root.rebalanceTree()
}

func (node *OngoingExperiment) add(newNode *OngoingExperiment) {
	if newNode.Id <= node.Id {
		if node.left == nil {
			node.left = newNode
		} else {
			node.left.add(newNode)
		}
	} else {
		if node.right == nil {
			node.right = newNode
		} else {
			node.right.add(newNode)
		}
	}
}

func (tree *OngoingExperiments) Remove(id int64) {
	tree.root = tree.root.remove(id)
}

func (tree *OngoingExperiments) Search(id int64) *OngoingExperiment {
	return tree.root.search(id)
}

func (node *OngoingExperiment) search(id int64) *OngoingExperiment {
	if node == nil {
		return nil
	}
	if id < node.Id {
		return node.left.search(id)
	} else if id > node.Id {
		return node.right.search(id)
	} else {
		return node
	}
}

func (node *OngoingExperiment) remove(id int64) *OngoingExperiment {
	if node == nil {
		return nil
	}
	if id < node.Id {
		node.left = node.left.remove(id)
	} else if id > node.Id {
		node.right = node.right.remove(id)
	} else {
		if node.left != nil && node.right != nil {
			rightMinNode := node.right.findSmallest()
			node.Id = rightMinNode.Id
			node.Finished = rightMinNode.Finished
			node.Proc = rightMinNode.Proc
			node.right = node.right.remove(rightMinNode.Id)
		} else if node.left != nil {
			node = node.left
		} else if node.right != nil {
			node = node.right
		} else {
			node = nil
			return node
		}

	}
	return node.rebalanceTree()
}

func (node *OngoingExperiment) rebalanceTree() *OngoingExperiment {
	if node == nil {
		return node
	}
	node.recalculateHeight()

	balanceFactor := node.left.getHeight() - node.right.getHeight()
	if balanceFactor == -2 {
		if node.right.left.getHeight() > node.right.right.getHeight() {
			node.right = node.right.rotateRight()
		}
		return node.rotateLeft()
	} else if balanceFactor == 2 {
		if node.left.right.getHeight() > node.left.left.getHeight() {
			node.left = node.left.rotateLeft()
		}
		return node.rotateRight()
	}
	return node
}

func (node *OngoingExperiment) rotateLeft() *OngoingExperiment {
	newRoot := node.right
	node.right = newRoot.left
	newRoot.left = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *OngoingExperiment) rotateRight() *OngoingExperiment {
	newRoot := node.left
	node.left = newRoot.right
	newRoot.right = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *OngoingExperiment) getHeight() int {
	if node == nil {
		return 0
	}
	return node.height
}

func (node *OngoingExperiment) recalculateHeight() {
	node.height = 1 + utils.Max(node.left.getHeight(), node.right.getHeight())
}

func (node *OngoingExperiment) findSmallest() *OngoingExperiment {
	if node.left != nil {
		return node.left.findSmallest()
	} else {
		return node
	}
}

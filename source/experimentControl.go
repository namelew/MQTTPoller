package main

import "os"

type ongoingExperiments struct{
	root *ongoingExperiment
}

type ongoingExperiment struct{
	id int64
	finished bool
	proc *os.Process
	height int
	right *ongoingExperiment 
	left *ongoingExperiment
}

func (tree *ongoingExperiments) add(node *ongoingExperiment){
	if tree.root == nil{
		tree.root = node
	} else{
		tree.root.add(node)
	}
	tree.root = tree.root.rebalanceTree()
}

func (node *ongoingExperiment) add(newNode *ongoingExperiment){
	if newNode.id <= node.id{
		if node.left == nil{
			node.left = newNode
		}else{
			node.left.add(newNode)
		}
	} else{
		if node.right == nil{
			node.right = newNode
		} else{
			node.right.add(newNode)
		}
	}
}

func (tree *ongoingExperiments) remove(id int64){
	tree.root  = tree.root.remove(id)
}

func (tree *ongoingExperiments) search(id int64) *ongoingExperiment{
	return tree.root.search(id)
}

func (node *ongoingExperiment) search(id int64) *ongoingExperiment{
	if node == nil {
		return nil
	}
	if id < node.id {
		return node.left.search(id)
	} else if id > node.id {
		return node.right.search(id)
	} else {
		return node
	}
}

func (node *ongoingExperiment) remove(id int64) *ongoingExperiment {
	if node == nil {
		return nil
	}
	if id < node.id {
		node.left = node.left.remove(id)
	} else if id > node.id {
		node.right = node.right.remove(id)
	} else {
		if node.left != nil && node.right != nil {
			rightMinNode := node.right.findSmallest()
			node.id = rightMinNode.id
			node.finished = rightMinNode.finished
			node.proc = rightMinNode.proc
			node.right = node.right.remove(rightMinNode.id)
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

func (node *ongoingExperiment) rebalanceTree() *ongoingExperiment {
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

func (node *ongoingExperiment) rotateLeft() *ongoingExperiment {
	newRoot := node.right
	node.right = newRoot.left
	newRoot.left = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *ongoingExperiment) rotateRight() *ongoingExperiment {
	newRoot := node.left
	node.left = newRoot.right
	newRoot.right = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *ongoingExperiment) getHeight() int {
	if node == nil {
		return 0
	}
	return node.height
}

func (node *ongoingExperiment) recalculateHeight() {
	node.height = 1 + max(node.left.getHeight(), node.right.getHeight())
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (node *ongoingExperiment) findSmallest() *ongoingExperiment {
	if node.left != nil {
		return node.left.findSmallest()
	} else {
		return node
	}
}
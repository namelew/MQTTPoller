package main

import (
	"fmt"
)

type status struct{
	Type string `json:"type"`
	Status string `json:"status"`
	Attr command `json:"attr"`
}

type session struct{
	Id int
	Finish bool
	Status status
	LogLevel int
	ToleranceLevel int
}

type experimentLog struct{
	id int64
	attempts int
	cmd command
	height int
	err bool
	finished bool
	left *experimentLog
	right *experimentLog
}

type experimentHistory struct{
	root *experimentLog
}

func (tree *experimentHistory) Add(id int64, cmd command, attemps int){
	node := experimentLog{id, attemps, cmd, 1,false, false,nil, nil}
	if tree.root == nil{
		tree.root = &node
	} else{
		tree.root.add(&node)
	}
	tree.root = tree.root.rebalanceTree()
}

func (node *experimentLog) add(newNode *experimentLog){
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

func (tree *experimentHistory) remove(id int64){
	tree.root  = tree.root.remove(id)
}

func (tree *experimentHistory) search(id int64) *experimentLog{
	return tree.root.search(id)
}

func (node *experimentLog) search(id int64) *experimentLog{
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

func (node *experimentLog) remove(id int64) *experimentLog {
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
			node.cmd = rightMinNode.cmd
			node.finished = rightMinNode.finished
			node.err = rightMinNode.err
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

func (node *experimentLog) rebalanceTree() *experimentLog {
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

func (node *experimentLog) rotateLeft() *experimentLog {
	newRoot := node.right
	node.right = newRoot.left
	newRoot.left = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *experimentLog) rotateRight() *experimentLog {
	newRoot := node.left
	node.left = newRoot.right
	newRoot.right = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *experimentLog) getHeight() int {
	if node == nil {
		return 0
	}
	return node.height
}

func (node *experimentLog) recalculateHeight() {
	node.height = 1 + max(node.left.getHeight(), node.right.getHeight())
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (node *experimentLog) findSmallest() *experimentLog {
	if node.left != nil {
		return node.left.findSmallest()
	} else {
		return node
	}
}

func (node *experimentLog) findLarger() *experimentLog{
	if node.right != nil {
		return node.right.findLarger()
	} else {
		return node
	}
}

func (tree *experimentHistory)sweep(test func(a *experimentLog) bool) *experimentLog{
	return tree.root.sweep(test)
}

func (node *experimentLog) sweep(test func(a *experimentLog) bool) *experimentLog{
	var result *experimentLog = nil
	if node.left != nil && result == nil{
		if test(node.left){
			result = node.left
		} else{
			result = node.left.sweep(test)
		}
	}

	if node.right != nil && result == nil{
		if test(node.right){
			result = node.left
		} else{
			result = node.right.sweep(test)
		}
	}

	if test(node){
		result = node
	}

	return result
}

func (tree *experimentHistory) Print() {
	tree.root.print()
}

func (node *experimentLog) print() {
	if node.left != nil {
		node.left.print()
	}
	fmt.Println("\n-----------------------")
	fmt.Printf("ID: %d\n", node.id)
	fmt.Printf("Finish: %t\n", node.finished)
	fmt.Println("\n-----------------------")
	if node.right != nil {
		node.right.print()
	}
}

func (tree *experimentHistory) GetUnfinish() *experimentLog{
	return tree.sweep(func(a *experimentLog) bool{
		return !a.finished
	})
}

func (tree *experimentHistory) Truncate(){
	tree.root = nil
}
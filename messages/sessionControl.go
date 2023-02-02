package messages

type Status struct {
	Type   string  `json:"type"`
	Status string  `json:"status"`
	Attr   Command `json:"attr"`
}

type Session struct {
	Id             int
	Finish         bool
	Status         Status
	LogLevel       int
	ToleranceLevel int
}

type ExperimentLog struct {
	Id       int64
	Attempts int
	Cmd      Command
	height   int
	Err      bool
	Finished bool
	left     *ExperimentLog
	right    *ExperimentLog
}

type ExperimentHistory struct {
	root *ExperimentLog
}

func (tree *ExperimentHistory) Add(id int64, cmd Command, attemps int) {
	node := ExperimentLog{id, attemps, cmd, 1, false, false, nil, nil}
	if tree.root == nil {
		tree.root = &node
	} else {
		tree.root.add(&node)
	}
	tree.root = tree.root.rebalanceTree()
}

func (node *ExperimentLog) add(newNode *ExperimentLog) {
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

func (tree *ExperimentHistory) Remove(id int64) {
	tree.root = tree.root.remove(id)
}

func (tree *ExperimentHistory) Search(id int64) *ExperimentLog {
	return tree.root.search(id)
}

func (node *ExperimentLog) search(id int64) *ExperimentLog {
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

func (node *ExperimentLog) remove(id int64) *ExperimentLog {
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
			node.Cmd = rightMinNode.Cmd
			node.Finished = rightMinNode.Finished
			node.Err = rightMinNode.Err
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

func (node *ExperimentLog) rebalanceTree() *ExperimentLog {
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

func (node *ExperimentLog) rotateLeft() *ExperimentLog {
	newRoot := node.right
	node.right = newRoot.left
	newRoot.left = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *ExperimentLog) rotateRight() *ExperimentLog {
	newRoot := node.left
	node.left = newRoot.right
	newRoot.right = node

	node.recalculateHeight()
	newRoot.recalculateHeight()
	return newRoot
}

func (node *ExperimentLog) getHeight() int {
	if node == nil {
		return 0
	}
	return node.height
}

func (node *ExperimentLog) recalculateHeight() {
	node.height = 1 + max(node.left.getHeight(), node.right.getHeight())
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (node *ExperimentLog) findSmallest() *ExperimentLog {
	if node.left != nil {
		return node.left.findSmallest()
	} else {
		return node
	}
}

func (tree *ExperimentHistory) FindLarger() *ExperimentLog {
	return tree.root.findLarger()
}

func (node *ExperimentLog) findLarger() *ExperimentLog {
	if node.right != nil {
		return node.right.findLarger()
	} else {
		return node
	}
}

func (tree *ExperimentHistory) Sweep(test func(a *ExperimentLog) bool) *ExperimentLog {
	return tree.root.sweep(test)
}

func (node *ExperimentLog) sweep(test func(a *ExperimentLog) bool) *ExperimentLog {
	var result *ExperimentLog = nil
	if node.left != nil && result == nil {
		if test(node.left) {
			result = node.left
		} else {
			result = node.left.sweep(test)
		}
	}

	if node.right != nil && result == nil {
		if test(node.right) {
			result = node.left
		} else {
			result = node.right.sweep(test)
		}
	}

	if test(node) {
		result = node
	}

	return result
}

func (tree *ExperimentHistory) Print(array []interface{}) {
	tree.root.print(array)
}

func (node *ExperimentLog) print(array []interface{}) {
	if node == nil {
		return
	}
	data := make(map[string]interface{})
	data["Id"] = node.Id
	data["Command"] = node.Cmd
	data["Finished"] = node.Finished

	if len(array) == 1 {
		array[0] = data
	} else {
		array = append(array, data)
	}

	if node.left != nil {
		node.left.print(array)
	}

	if node.right != nil {
		node.right.print(array)
	}
}

func (tree *ExperimentHistory) GetUnfinish() *ExperimentLog {
	return tree.Sweep(func(a *ExperimentLog) bool {
		return !a.Finished
	})
}

func (tree *ExperimentHistory) Truncate() {
	tree.root = nil
}

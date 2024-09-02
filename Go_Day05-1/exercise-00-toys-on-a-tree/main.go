package main

import (
	"fmt"
)

// TreeNode представляет собой узел бинарного дерева
type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

// areToysBalanced проверяет, уравновешено ли дерево по количеству игрушек на левой и правой ветвях
func areToysBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return countToys(root.Left) == countToys(root.Right)
}

// countToys считает количество игрушек в поддереве
func countToys(node *TreeNode) int {
	if node == nil {
		return 0
	}
	count := 0
	if node.HasToy {
		count = 1
	}
	return count + countToys(node.Left) + countToys(node.Right)
}

// drawTree отрисовывает бинарное дерево в виде строк
func drawTree(node *TreeNode, prefix string, isLast bool) {
	if node != nil {
		// Определяем символ для текущего узла
		var nodeSymbol string
		if node.HasToy {
			nodeSymbol = "1"
		} else {
			nodeSymbol = "0"
		}

		// Отрисовка текущего узла
		if isLast {
			fmt.Println(prefix + "└── " + nodeSymbol)
			prefix += "    "
		} else {
			fmt.Println(prefix + "├── " + nodeSymbol)
			prefix += "│   "
		}

		// Рекурсивно отрисовываем поддеревья
		if node.Left != nil || node.Right != nil {
			drawTree(node.Left, prefix, node.Right == nil)
			drawTree(node.Right, prefix, true)
		}
	}
}

// promptForBool запрашивает у пользователя значение bool
func promptForBool(prompt string) bool {
	var input string
	for {
		fmt.Printf("%s (true/false): ", prompt)
		fmt.Scanln(&input)
		if input == "true" {
			return true
		} else if input == "false" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'true' or 'false'.")
		}
	}
}

// promptForSubtree запрашивает, какие поддеревья создать для узла
func promptForSubtree(prompt string) string {
	var input string
	for {
		fmt.Printf("%s (left/right/both/none): ", prompt)
		fmt.Scanln(&input)
		if input == "left" || input == "right" || input == "both" || input == "none" {
			return input
		} else {
			fmt.Println("Invalid input. Please enter 'left', 'right', 'both', or 'none'.")
		}
	}
}

func main() {
	// Запрашиваем у пользователя значения для каждого узла
	rootToy := promptForBool("Does the root node have a toy?")
	leftToy := promptForBool("Does the left child node have a toy?")
	rightToy := promptForBool("Does the right child node have a toy?")

	// Проверяем, какие поддеревья создать для левого дочернего узла
	var leftSubtree *TreeNode
	leftSubtreeChoice := promptForSubtree("Which subtree does the left child node have?")
	switch leftSubtreeChoice {
	case "left":
		leftLeftToy := promptForBool("  Does the left subtree, left child node have a toy?")
		leftSubtree = &TreeNode{
			HasToy: leftToy,
			Left:   &TreeNode{HasToy: leftLeftToy},
		}
	case "right":
		leftRightToy := promptForBool("  Does the left subtree, right child node have a toy?")
		leftSubtree = &TreeNode{
			HasToy: leftToy,
			Right:  &TreeNode{HasToy: leftRightToy},
		}
	case "both":
		leftLeftToy := promptForBool("  Does the left subtree, left child node have a toy?")
		leftRightToy := promptForBool("  Does the left subtree, right child node have a toy?")
		leftSubtree = &TreeNode{
			HasToy: leftToy,
			Left:   &TreeNode{HasToy: leftLeftToy},
			Right:  &TreeNode{HasToy: leftRightToy},
		}
	default:
		leftSubtree = &TreeNode{HasToy: leftToy}
	}

	// Проверяем, какие поддеревья создать для правого дочернего узла
	var rightSubtree *TreeNode
	rightSubtreeChoice := promptForSubtree("Which subtree does the right child node have?")
	switch rightSubtreeChoice {
	case "left":
		rightLeftToy := promptForBool("  Does the right subtree, left child node have a toy?")
		rightSubtree = &TreeNode{
			HasToy: rightToy,
			Left:   &TreeNode{HasToy: rightLeftToy},
		}
	case "right":
		rightRightToy := promptForBool("  Does the right subtree, right child node have a toy?")
		rightSubtree = &TreeNode{
			HasToy: rightToy,
			Right:  &TreeNode{HasToy: rightRightToy},
		}
	case "both":
		rightLeftToy := promptForBool("  Does the right subtree, left child node have a toy?")
		rightRightToy := promptForBool("  Does the right subtree, right child node have a toy?")
		rightSubtree = &TreeNode{
			HasToy: rightToy,
			Left:   &TreeNode{HasToy: rightLeftToy},
			Right:  &TreeNode{HasToy: rightRightToy},
		}
	default:
		rightSubtree = &TreeNode{HasToy: rightToy}
	}

	// Создаем дерево на основе введенных данных
	root := &TreeNode{
		HasToy: rootToy,
		Left:   leftSubtree,
		Right:  rightSubtree,
	}

	// Отрисовываем дерево в терминале
	fmt.Println("Binary Tree:")
	drawTree(root, "", true)

	// Проверяем, уравновешено ли дерево
	if areToysBalanced(root) {
		fmt.Println("The tree is balanced - true.")
	} else {
		fmt.Println("The tree is not balanced - false.")
	}
}

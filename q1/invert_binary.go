package q1

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

// invertTree inverts a binary tree without using recursion.
// using bfs to traverse the tree and swap the left and right children of each node.
func invertTree(root *TreeNode) *TreeNode {
    if root == nil {
        return nil
    }

    queue := []*TreeNode{root}
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        // Swap the left and right children
        current.Left, current.Right = current.Right, current.Left

        // If the left child exists, add it to the queue
        if current.Left != nil {
            queue = append(queue, current.Left)
        }

        // If the right child exists, add it to the queue
        if current.Right != nil {
            queue = append(queue, current.Right)
        }
    }

    return root
}
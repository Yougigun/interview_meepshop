package q1

import (
	"reflect"
	"testing"
)

// createBinaryTree creates a binary tree from a slice of interface{} values.
func TestInvertTree(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{} // Use interface{} to handle nil values in the tree
		expected []interface{}
	}{
		{
			name:     "Example 1",
			input:    []interface{}{5, 3, 8, 1, 7, 2, 6},
			expected: []interface{}{5, 8, 3, 6, 2, 7, 1},
		},
		{
			name:     "Example 2",
			input:    []interface{}{6, 8, 9},
			expected: []interface{}{6, 9, 8},
		},
		{
			name:     "Example 3",
			input:    []interface{}{5, 3, 8, 1, 7, 2, 6, 100, 3, -1},
			expected: []interface{}{5, 8, 3, 6, 2, 7, 1, nil, nil, nil, nil, nil, -1, 3, 100},
		},
		{
			name:     "Example 4",
			input:    []interface{}{},
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputTree := createBinaryTree(tt.input)
			expectedTree := createBinaryTree(tt.expected)
			invertedTree := invertTree(inputTree)
			if !reflect.DeepEqual(treeToSlice(invertedTree), treeToSlice(expectedTree)) {
				t.Errorf("%s: got %v, want %v", tt.name, treeToSlice(invertedTree), tt.expected)
			}
		})
	}
}

// Helper function to create a binary tree from a slice of integers.
// This function assumes the input slice represents a binary tree
// in level order and uses 'nil' to indicate missing nodes.
func createBinaryTree(values []interface{}) *TreeNode {
	if len(values) == 0 {
		return nil
	}

	root := &TreeNode{Val: values[0].(int)}
	queue := []*TreeNode{root}

	for i := 1; i < len(values); i += 2 {
		current := queue[0]
		queue = queue[1:]

		if values[i] != nil {
			current.Left = &TreeNode{Val: values[i].(int)}
			queue = append(queue, current.Left)
		}
		if i+1 < len(values) && values[i+1] != nil {
			current.Right = &TreeNode{Val: values[i+1].(int)}
			queue = append(queue, current.Right)
		}
	}

	return root
}

// treeToSlice converts a binary tree to a slice in level order.
// This helps in comparing the actual and expected trees.
func treeToSlice(root *TreeNode) []interface{} {
	var result []interface{}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == nil {
			result = append(result, nil)
			continue
		}
		result = append(result, current.Val)
		queue = append(queue, current.Left, current.Right)
	}
	return result
}

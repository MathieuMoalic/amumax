package script

func Contains(tree, search Expr) bool {
	if tree == search {
		return true
	}
	children := tree.Child()
	for _, e := range children {
		if Contains(e, search) {
			return true
		}
	}
	return false
}

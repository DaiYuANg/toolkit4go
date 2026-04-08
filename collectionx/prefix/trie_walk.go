package prefix

func (t *Trie[V]) collectPairs(node *trieNode[V], path []rune, out *[]keyValue[V]) {
	if node == nil {
		return
	}
	if node.hasValue {
		*out = append(*out, keyValue[V]{
			key:   string(path),
			value: node.value,
		})
	}
	if node.children.Len() == 0 {
		return
	}

	for _, ch := range node.childKeys {
		child, _ := node.children.Get(ch)
		t.collectPairs(child, append(path, ch), out)
	}
}

func (t *Trie[V]) collectKeys(node *trieNode[V], path []rune, out *[]string) {
	if node == nil {
		return
	}
	if node.hasValue {
		*out = append(*out, string(path))
	}
	if node.children.Len() == 0 {
		return
	}

	for _, ch := range node.childKeys {
		child, _ := node.children.Get(ch)
		t.collectKeys(child, append(path, ch), out)
	}
}

func (t *Trie[V]) collectValues(node *trieNode[V], out *[]V) {
	if node == nil {
		return
	}
	if node.hasValue {
		*out = append(*out, node.value)
	}
	if node.children.Len() == 0 {
		return
	}

	for _, ch := range node.childKeys {
		child, _ := node.children.Get(ch)
		t.collectValues(child, out)
	}
}

func (t *Trie[V]) rangePrefix(node *trieNode[V], path []rune, fn func(key string, value V) bool) bool {
	if node == nil {
		return true
	}
	if node.hasValue && !fn(string(path), node.value) {
		return false
	}
	if node.children.Len() == 0 {
		return true
	}

	for _, ch := range node.childKeys {
		child, _ := node.children.Get(ch)
		if !t.rangePrefix(child, append(path, ch), fn) {
			return false
		}
	}
	return true
}

func (t *Trie[V]) pairsWithPrefix(prefix string) []keyValue[V] {
	if t == nil || t.root == nil {
		return nil
	}
	startNode, ok := t.findNode(prefix)
	if !ok {
		return nil
	}

	out := make([]keyValue[V], 0, startNode.valueCount)
	t.collectPairs(startNode, []rune(prefix), &out)
	return out
}

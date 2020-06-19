package bounded_tree

import (
	"fmt"
	"testing"

	"github.com/chrislusf/seaweedfs/weed/util"
)


var (

	visitFn = func(path util.FullPath) (childDirectories []string, err error) {
		fmt.Printf("  visit %v ...\n", path)
		switch path {
		case "/":
			return []string{"a", "g", "h"}, nil
		case "/a":
			return []string{"b", "f"}, nil
		case "/a/b":
			return []string{"c", "e"}, nil
		case "/a/b/c":
			return []string{"d"}, nil
		case "/a/b/c/d":
			return []string{"i", "j"}, nil
		case "/a/b/c/d/i":
			return []string{}, nil
		case "/a/b/c/d/j":
			return []string{}, nil
		case "/a/b/e":
			return []string{}, nil
		case "/a/f":
			return []string{}, nil
		}
		return nil, nil
	}


	printMap = func(m map[string]*Node) {
		for k := range m {
			println("  >", k)
		}
	}


)

func TestBoundedTree(t *testing.T) {

	// a/b/c/d/i
	// a/b/c/d/j
	// a/b/c/d
	// a/b/e
	// a/f
	// g
	// h

	tree := NewBoundedTree()

	tree.EnsureVisited(util.FullPath("/a/b/c"), visitFn)

	printMap(tree.root.Children)

	a := tree.root.getChild("a")

	b := a.getChild("b")
	if !b.isVisited() {
		t.Errorf("expect visited /a/b")
	}
	c := b.getChild("c")
	if !c.isVisited() {
		t.Errorf("expect visited /a/b/c")
	}

	d := c.getChild("d")
	if d.isVisited() {
		t.Errorf("expect unvisited /a/b/c/d")
	}

	tree.EnsureVisited(util.FullPath("/a/b/c/d"), visitFn)
	tree.EnsureVisited(util.FullPath("/a/b/c/d/i"), visitFn)
	tree.EnsureVisited(util.FullPath("/a/b/c/d/j"), visitFn)
	tree.EnsureVisited(util.FullPath("/a/b/e"), visitFn)
	tree.EnsureVisited(util.FullPath("/a/f"), visitFn)

	printMap(tree.root.Children)

}

func TestEmptyBoundedTree(t *testing.T) {

	// g
	// h

	tree := NewBoundedTree()

	visitFn := func(path util.FullPath) (childDirectories []string, err error) {
		fmt.Printf("  visit %v ...\n", path)
		switch path {
		case "/":
			return []string{"g", "h"}, nil
		}
		t.Fatalf("expected visit %s", path)
		return nil, nil
	}

	tree.EnsureVisited(util.FullPath("/a/b"), visitFn)

	tree.EnsureVisited(util.FullPath("/a/b"), visitFn)

	printMap(tree.root.Children)

	println(tree.HasVisited(util.FullPath("/a/b")))
	println(tree.HasVisited(util.FullPath("/a")))
	println(tree.HasVisited(util.FullPath("/g")))
	println(tree.HasVisited(util.FullPath("/g/x")))

}

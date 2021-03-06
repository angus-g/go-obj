// Package obj allows simple parsing of Wavefront OBJ files
// into slices of vertices, normals and elements.
package obj

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Parse a vertex (normal) string, a list of whitespace-separated
// floating point numbers.
func parseVertex(t []string) []float32 {
	x, _ := strconv.ParseFloat(t[0], 32)
	y, _ := strconv.ParseFloat(t[1], 32)
	z, _ := strconv.ParseFloat(t[2], 32)

	return []float32{float32(x), float32(y), float32(z)}
}

// Parse an element string, a list of whitespace-separated elements.
// Elements are of the form "<vi>/<ti>/<ni>" where indices are the
// vertex, texture coordinate and normal, respectively.
func parseElement(t []string) [][3]int32 {
	e := make([][3]int32, len(t))

	for i := 0; i < len(t); i++ {
		f := strings.Split(t[i], "/")

		for j := 0; j < len(f); j++ {
			// for now, just grab the vertex index
			if x, err := strconv.ParseInt(f[j], 10, 32); err == nil {
				e[i][j] = int32(x) - 1 // convert to 0-indexing
			} else {
				e[i][j] = -1
			}
		}
	}

	// convert quads to triangles
	if len(t) > 3 {
		e = append(e, e[0], e[2])
	}

	return e
}

func Parse(filename string) ([]float32, []float32) {
	fp, _ := os.Open(filename)
	scanner := bufio.NewScanner(fp)

	vertices := [][]float32{}
	normals := [][]float32{}
	elements := [][3]int32{}

	vertOut := []float32{}
	normOut := []float32{}

	for scanner.Scan() {
		toks := strings.Fields(strings.TrimSpace(scanner.Text()))

		switch toks[0] {
		case "v":
			vertices = append(vertices, parseVertex(toks[1:]))
		case "vn":
			normals = append(normals, parseVertex(toks[1:]))
		case "f":
			elements = append(elements, parseElement(toks[1:])...)
		}
	}

	for _, e := range elements {
		if e[0] >= 0 {
			vertOut = append(vertOut, vertices[e[0]]...)
		}
		if e[2] >= 0 {
			normOut = append(normOut, normals[e[2]]...)
		}
	}

	return vertOut, normOut
}

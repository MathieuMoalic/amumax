package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/util"
)

func init() {
	DeclFunc("Chunk", Mx3chunks, "")
}

type Chunk struct {
	len int
	nb  int
}

type Chunks struct {
	x Chunk
	y Chunk
	z Chunk
	c Chunk
}

type RequestedChunking struct {
	x int
	y int
	z int
	c int
}

func Mx3chunks(x, y, z, c int) RequestedChunking {
	return RequestedChunking{x, y, z, c}
}

func NewChunks(q Quantity, c RequestedChunking) Chunks {
	size := SizeOf(q)
	return Chunks{
		NewChunk(size[0], c.x, 0),
		NewChunk(size[1], c.y, 1),
		NewChunk(size[2], c.z, 2),
		NewChunk(q.NComp(), c.c, 3),
	}
}

func NewChunk(length, nb_of_chunks, N_index int) Chunk {
	name := []string{"Nx", "Ny", "Nz", "comp"}[N_index]
	if nb_of_chunks < 1 || (nb_of_chunks > length) {
		util.Fatal("Error: The number of chunks must be between 1 and ", name)
	}
	new_nb_of_chunks := closestDivisor(length, nb_of_chunks)
	if new_nb_of_chunks != nb_of_chunks {
		LogOut("Warning: The number of chunks for", name, "has been automatically resized from", nb_of_chunks, "to", new_nb_of_chunks)
	}
	nb_of_chunks = new_nb_of_chunks
	return Chunk{length / nb_of_chunks, nb_of_chunks}
}

func closestDivisor(N int, D int) int {
	closest := 0
	minDist := math.MaxInt32
	for i := 1; i <= N; i++ {
		if N%i == 0 {
			dist := i - D
			if dist < 0 {
				dist = -dist
			}
			if dist < minDist {
				minDist = dist
				closest = i
			}
		}
	}
	return closest
}

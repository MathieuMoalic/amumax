package engine_old

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/log_old"
)

type chunk struct {
	len int
	nb  int
}

type chunks struct {
	x chunk
	y chunk
	z chunk
	c chunk
}

type requestedChunking struct {
	x int
	y int
	z int
	c int
}

func mx3chunks(x, y, z, c int) requestedChunking {
	return requestedChunking{x, y, z, c}
}

func newChunks(q Quantity, c requestedChunking) chunks {
	size := sizeOf(q)
	return chunks{
		newChunk(size[0], c.x, 0),
		newChunk(size[1], c.y, 1),
		newChunk(size[2], c.z, 2),
		newChunk(q.NComp(), c.c, 3),
	}
}

func newChunk(length, nb_of_chunks, N_index int) chunk {
	name := []string{"Nx", "Ny", "Nz", "comp"}[N_index]
	if nb_of_chunks < 1 || (nb_of_chunks > length) {
		log_old.Log.ErrAndExit("Error: The number of chunks must be between 1 and %v", name)
	}
	new_nb_of_chunks := closestDivisor(length, nb_of_chunks)
	if new_nb_of_chunks != nb_of_chunks {
		log_old.Log.Info("Warning: The number of chunks for %v has been automatically resized from %v to %v", name, nb_of_chunks, new_nb_of_chunks)
	}
	nb_of_chunks = new_nb_of_chunks
	return chunk{length / nb_of_chunks, nb_of_chunks}
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

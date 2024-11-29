package chunk

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/quantity"
)

type chunk struct {
	Len   int
	Count int
}

type Chunks struct {
	X chunk
	Y chunk
	Z chunk
	C chunk
}

type RequestedChunking struct {
	X int
	Y int
	Z int
	C int
}

func CreateRequestedChunk(x, y, z, c int) RequestedChunking {
	return RequestedChunking{x, y, z, c}
}

func NewChunks(log *log.Logs, q quantity.Quantity, c RequestedChunking) Chunks {
	return Chunks{
		newChunk(log, q.Size()[0], c.X, 0),
		newChunk(log, q.Size()[1], c.Y, 1),
		newChunk(log, q.Size()[2], c.Z, 2),
		newChunk(log, q.NComp(), c.C, 3),
	}
}

func newChunk(log *log.Logs, length, nb_of_chunks, N_index int) chunk {
	name := []string{"Nx", "Ny", "Nz", "comp"}[N_index]
	if nb_of_chunks < 1 || (nb_of_chunks > length) {
		log.ErrAndExit("Error: The number of chunks must be between 1 and %v", name)
	}
	new_nb_of_chunks := closestDivisor(length, nb_of_chunks)
	if new_nb_of_chunks != nb_of_chunks {
		log.Info("Warning: The number of chunks for %v has been automatically resized from %v to %v", name, nb_of_chunks, new_nb_of_chunks)
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

package engine

import (
	"math"

	"github.com/MathieuMoalic/amumax/src/log"
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

func newChunk(length, nbOfChunks, NIndex int) chunk {
	name := []string{"Nx", "Ny", "Nz", "comp"}[NIndex]
	if nbOfChunks < 1 || (nbOfChunks > length) {
		log.Log.ErrAndExit("Error: The number of chunks must be between 1 and %v", name)
	}
	newNbOfChunks := closestDivisor(length, nbOfChunks)
	if newNbOfChunks != nbOfChunks {
		log.Log.Info("Warning: The number of chunks for %v has been automatically resized from %v to %v", name, nbOfChunks, newNbOfChunks)
	}
	nbOfChunks = newNbOfChunks
	return chunk{length / nbOfChunks, nbOfChunks}
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

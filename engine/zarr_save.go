package engine

// Bookkeeping for auto-saving quantities at given intervals.

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/DataDog/zstd"

	"github.com/MathieuMoalic/amumax/cuda"
	"github.com/MathieuMoalic/amumax/data"
	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

func init() {
	DeclFunc("AutoSaveAs", zAutoSaveAs, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("AutoSave", zAutoSave, "Auto save space-dependent quantity every period (s) as the zarr standard.")
	DeclFunc("SaveAs", zSaveAs, "Save space-dependent quantity as the zarr standard.")
	DeclFunc("Save", zSave, "Save space-dependent quantity as the zarr standard.")
	DeclFunc("chunkxyzc", chunkxyzc, "chunkxyzc")
}

var zGroups []string
var zArrays = make(map[string]*zArray)
var chunks Chunks

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

func chunkxyzc(x, y, z, c int) {
	meshsize := globalmesh_.Size()
	// fmt.Println(meshsize, x, y, z)
	if meshsize[0] == 0 {
		util.Fatal("Error: You have to define the mesh before defining the chunks")
	} else if (z < 1) || (y < 1) || (x < 1) || (c < 1) {
		util.Fatal("Error: Chunks must be strictly positive")
	} else if (z > meshsize[Z]) || (y > meshsize[Y]) || (x > meshsize[X]) {
		util.Fatal("Error: Chunks must be smaller or equal to the number of cells")
	} else if (meshsize[Z]%z != 0) || (meshsize[Y]%y != 0) || (meshsize[X]%x != 0) {
		util.Fatal("Error: Chunks must fit an integer number of times in the mesh")
	} else if (c != 1) && (c != 3) {
		util.Fatal("Error: Chunks for the magnetization components can only be 1 or 3")
	} else if x*y*z*c*4 < 1000 {
		util.Fatal("Error: Chunks are too small, chunks around 1 MB give the best performance, current chunks: ", float32(x*y*z*c*4)/1e6, " MB")
	} else {
		chunks = Chunks{
			Chunk{x, meshsize[X] / x},
			Chunk{y, meshsize[Y] / y},
			Chunk{z, meshsize[Z] / z},
			Chunk{c, 3 / c},
		}
		fmt.Println("Chunk size: ", float32(x*y*z*c*4)/1e6, " MB")
	}
}

type zArray struct {
	name   string
	q      Quantity
	period float64                // How often to save
	start  float64                // Starting point
	count  int                    // Number of times it has been autosaved
	save   func(Quantity, string) // called to do the actual save
	times  []float64
}

// returns true when the time is right to save.
func (a *zArray) needSave() bool {
	return a.period != 0 && (Time-a.start)-float64(a.count)*a.period >= a.period
}
func (a *zArray) SaveAttrs() {
	// it's stupid and wasteful but it works
	// keeping the whole array of times wastes a few MB of RAM
	u, err := json.Marshal(zarr.Zattrs{Buffer: a.times})
	util.FatalErr(err)
	httpfs.Remove(OD() + a.name + "/.zattrs")
	httpfs.Put(OD()+a.name+"/.zattrs", []byte(string(u)))

}

func zAutoSave(q Quantity, period float64) {
	zAutoSaveAs(q, NameOf(q), period)
}

func zInitArray(name string, q Quantity, period float64) {
	httpfs.Mkdir(OD() + name)
	var b []float64
	zArrays[name] = &zArray{name, q, period, Time, -1, zSaveAs, b} // init count to -1 allows save at t=0

}

// period == 0 stops autosaving.
func zAutoSaveAs(q Quantity, name string, period float64) {
	if period == 0 {
		delete(zArrays, name)
	} else {
		zInitArray(name, q, period)
	}
}

func zSaveAs(q Quantity, name string) {
	if _, ok := zArrays[name]; !ok {
		zInitArray(name, q, 0)
	}
	zArrays[name].times = append(zArrays[name].times, Time)
	zArrays[name].SaveAttrs()
	buffer := ValueOf(q)
	defer cuda.Recycle(buffer)
	data := buffer.HostCopy() // must be copy (async io)
	t := zArrays[name].count  // no desync this way
	queOutput(func() { zSyncSave(data, name, t) })
	zArrays[name].count++

}

// Save once, with auto file name
func zSave(q Quantity) {
	zSaveAs(q, NameOf(q))
}

func IntToByteArray(num int32) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

// synchronous chunky save
func zSyncSave(array *data.Slice, qname string, time int) {
	data := array.Tensors()
	size := array.Size()
	ncomp := array.NComp()

	// fmt.Println(chunks)
	// for every chunk
	count := 0
	for icx := 0; icx < chunks.x.nb; icx++ {
		for icy := 0; icy < chunks.y.nb; icy++ {
			for icz := 0; icz < chunks.z.nb; icz++ {
				for icc := 0; icc < chunks.c.nb; icc++ {
					// fmt.Println("-----------------------")
					// fmt.Println(icx, icy, icz, icc)
					// fmt.Println("-----------------------")
					// for every cell in that chunk
					bdata := []byte{}
					var bytes []byte
					f, err := httpfs.Create(fmt.Sprintf(OD()+"%s/%d.%d.%d.%d.%d", qname, time+1, icz, icy, icx, icc))
					util.FatalErr(err)
					defer f.Close()
					for iz := 0; iz < chunks.z.len; iz++ {
						z := icz*chunks.z.len + iz
						for iy := 0; iy < chunks.y.len; iy++ {
							y := icy*chunks.y.len + iy
							for ix := 0; ix < chunks.x.len; ix++ {
								x := icx*chunks.x.len + ix
								// fmt.Println(count, x, y, z)
								for ic := 0; ic < chunks.c.len; ic++ {
									c := icc*chunks.c.len + ic
									bytes = (*[4]byte)(unsafe.Pointer(&data[c][z][y][x]))[:]
									// bytes = IntToByteArray(int32(count))
									// fmt.Println(count, int32(count), bytes, byte(count))
									for k := 0; k < 4; k++ {
										bdata = append(bdata, bytes[k])
									}
								}
								count++
							}
						}
					}
					compressedData, err := zstd.Compress(nil, bdata)
					util.FatalErr(err)
					f.Write(compressedData)
				}
			}
		}
	}
	//.zarray file
	zarr.SaveFileZarray(fmt.Sprintf(OD()+"%s/.zarray", qname), size, ncomp, time+1, chunks.z.len, chunks.y.len, chunks.x.len, chunks.c.len)
}

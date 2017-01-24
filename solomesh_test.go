package recast

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aurelien-rainone/go-detour/recast"
)

func compareFiles(fn1, fn2 string) (bool, error) {
	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	var (
		f1, f2 []byte
		err    error
	)
	f1, err = ioutil.ReadFile(fn1)
	if err != nil {
		return false, err
	}

	f2, err = ioutil.ReadFile(fn2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(f1, f2), nil
}

func testCreateSoloNavMesh(t *testing.T, meshname string) {
	var (
		meshName, orgNavMeshName, tmpBin string
		ctx                              *recast.BuildContext
		soloMesh                         *SoloMesh
	)

	meshName = "testdata/" + meshname + ".obj"
	orgNavMeshName = "testdata/" + meshname + ".org.bin"
	tmpBin = "out.bin"

	ctx = recast.NewBuildContext(false)
	soloMesh = NewSoloMesh(ctx)
	if !soloMesh.Load(meshName) {
		t.Fatalf("couldn't load mesh %v", meshName)
	}
	navMesh, ok := soloMesh.Build()
	if !ok {
		t.Fatalf("couldn't build navmesh for %v", meshname)
	}

	navMesh.SaveToFile(tmpBin)
	t.Logf("%v navmesh successfully built", meshname)
	ok, err := compareFiles(tmpBin, orgNavMeshName)
	if err != nil {
		t.Fatalf("couldn't compare %v and %v, %v", tmpBin, orgNavMeshName, err)
	}
	if !ok {
		t.Fatalf("%v and %v are different", tmpBin, orgNavMeshName)
	}
	ctx.DumpLog("Navmesh Build log")
}

func TestCreateDungeonNavMesh(t *testing.T) {
	testCreateSoloNavMesh(t, "dungeon")
}

func TestCreateCubeNavMesh(t *testing.T) {
	testCreateSoloNavMesh(t, "cube")
}

func benchmarkCreateNavMesh(b *testing.B, meshname string) {
	meshName := "testdata/" + meshname + ".obj"

	soloMesh := NewSoloMesh(recast.NewBuildContext(false))
	if !soloMesh.Load(meshName) {
		b.Fatalf("couldn't load mesh %v", meshName)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, ok := soloMesh.Build()
		if !ok {
			b.Fatalf("couldn't build navmesh for %v", meshname)
		}
	}
}

func BenchmarkCreateCubeNavMesh(b *testing.B) {
	benchmarkCreateNavMesh(b, "cube")
}
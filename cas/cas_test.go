package cas_test

import (
	"testing"
)

func TestStorage(t *testing.T) {
	/*
		assert(t, os.RemoveAll("temp"))
		var (
			fsys = os.DirFS("../testdata")
			gm   = memgrf.NewMapper()
			gs   = memgrf.NewStore()
		)
		bs, err := localdisk.NewStorage("temp/file", localdisk.Buckets(1, 2))
		assert(t, err)

		g := grf.New(gm, gs)
		g.Register(
			grf.TypeInfo{
				Type:     "file",
				DataType: cas.NodeData{},
			})

		s := cas.New(bs, g, "file", file.MediaAttrs(), file.Digests(crypto.MD5))
		assert(t, fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			n, err := s.Add(blob.FromFS(fsys, path))
			if err != nil {
				return err
			}
			t.Logf("%s %v", path, n)
			return nil
		}))
		keyMapData, err := json.MarshalIndent(gm, "", "\t")
		assert(t, err)
		assert(t, localfile.Write("temp/keymap.json", blob.FromBytes(keyMapData), 0600))

		graphData, err := json.MarshalIndent(gs, "", "\t")
		assert(t, err)
		assert(t, localfile.Write("temp/graph.json", blob.FromBytes(graphData), 0600))
	*/
}

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

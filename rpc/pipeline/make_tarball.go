package pipeline

import (
	"archive/tar"
	"bytes"
	"io"

	"dflimg/lib/cher"
)

// MakeTarball makes a TAR of the collection of files you provide
func MakeTarball(p *Pipeline) (bool, error) {
	if !p.context.multifile {
		return true, nil
	}

	// we are assuming by this stage we have a collection of files and nothing
	// else wants them, so no highlighting has occured yet

	p.w.Header().Set("Content-Type", "application/x-tar")
	p.w.Header().Set("Content-Disposition", "attachment; filename=collection.tar")

	tarWriter := tar.NewWriter(p.w)
	defer tarWriter.Close()

	dupeNames := make(map[string]struct{})

	for _, i := range p.rwqs {
		rwq := i

		name := rwq.qi.Original

		if rwq.r.Name != nil {
			name = *rwq.r.Name
		}

		if _, ok := dupeNames[name]; ok {
			return false, cher.New("duplicate_name", cher.M{"dupes": dupeNames})
		}

		dupeNames[name] = struct{}{}

		err := tarWriter.WriteHeader(&tar.Header{
			Name:    name,
			Size:    int64(len(p.contents[rwq.r.ID].bytes)),
			Mode:    int64(0777),
			ModTime: *p.contents[rwq.r.ID].modtime,
		})
		if err != nil {
			return false, cher.New("cant_write_tar_header", cher.M{"resource_id": rwq.r.ID})
		}

		_, err = io.Copy(tarWriter, bytes.NewReader(p.contents[rwq.r.ID].bytes))
		if err != nil {
			return false, cher.New("cant_write_tar_body", cher.M{"resource_id": rwq.r.ID})
		}
	}

	return false, nil
}

package bins

func (d *binsDir) NewUp(name string, b io.Reader) error {
	// Create and open a temp file (acts as a lock and a temporary place for incomplete bytes).
	// Should reside in the bins directory to guarantee atomic rename.
	tempPath, err := d.tempPathCreate(name)
	if err != nil {
		return multierr.Combine(errors.New("create temp file failed"), err)
	}

	defer func() {
		err := removeForce(tempPath)
		if err != nil {
			logs.Warnf("remove temp file failed; %v", err)
		}
	}()

	tempFile, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return multierr.Combine(errors.New("open temp file failed"), err)
	}

	defer func() {
		err := tempFile.Close()
		if err != nil {
			logs.Warnf("close temp file failed; %v", err)
		}
	}()

	// Write into the temp file.
	_, err = tempFile.ReadFrom(b)
	if err != nil {
		return multierr.Combine(errors.New("read/write failed"), err)
	}

	// Move the temp file to the actual file.
	path := d.binPath(name)
	err = os.Rename(tempPath, path)
	if err != nil {
		return multierr.Combine(errors.New("rename temp file failed"), err)
	}

	return nil
}

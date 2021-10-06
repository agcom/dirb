package bins

func mergeJsnObjRec(j1, j2 interface{}) interface{} {
	if j1m, j1ok := j1.(map[string]interface{}); j1ok {
		if j2m, j2ok := j2.(map[string]interface{}); j2ok {
			return mergeMapsRec(j1m, j2m)
		}
	}

	return j2
}

func mergeMapsRec(m1, m2 map[string]interface{}) map[string]interface{} {
	r := make(map[string]interface{}, len(m1)+len(m2))
	for k, v1 := range m1 {
		r[k] = v1
	}

	for k, v2 := range m2 {
		if vr, ok := r[k]; ok {
			if vrm, okr := vr.(map[string]interface{}); okr {
				if v2m, ok2 := v2.(map[string]interface{}); ok2 {
					v2 = mergeMapsRec(vrm, v2m)
				}
			}
		}

		r[k] = v2
	}

	return r
}

func (b *Bins) Pch(name string, j interface{}) error {
	jo, err := b.Get(name)
	if err != nil {
		return err
	}

	jn := mergeJsnObjRec(jo, j)

	return b.NewUp(name, jn)
}

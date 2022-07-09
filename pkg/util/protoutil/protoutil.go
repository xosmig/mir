package protoutil

type PbWrapper[T any] interface {
	Pb() T
}

//func SliceFromPb[W PbWrapper[W, PB], PB any](pbs []PB) []W {
//	ws := make([]W, len(pbs))
//	for i := range pbs {
//		ws[i] = ws[i].FromPb(pbs[i])
//	}
//	return ws
//}

func SliceToPb[PB any, W PbWrapper[PB]](ws []W) []PB {
	pbs := make([]PB, len(ws))
	for i := range ws {
		pbs[i] = ws[i].Pb()
	}
	return pbs
}

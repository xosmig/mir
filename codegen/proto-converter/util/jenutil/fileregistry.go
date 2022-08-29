package jenutil

//import "github.com/dave/jennifer/jen"
//
//type FileRegistry map[string]*jen.File
//
//func NewFileRegistry() *FileRegistry {
//	return &FileRegistry{
//		fileByPbPackagePath: make(map[string]*jen.File),
//	}
//}
//
//func (fr *FileRegistry) GetFile(outputPackage string) *jen.File {
//	if g, ok := fr.fileByPbPackagePath[outputPackage]; ok {
//		return g
//	}
//
//	g := jen.NewFilePath(outputPackage)
//	fr.fileByPbPackagePath[outputPackage] = g
//	return g
//}

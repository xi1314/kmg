package kmgGoParser

import (
	"fmt"
	"path"
)

func MustWriteGoTypes(thisPackagePath string, typi Type) (s string, addPkgPathList []string) {
	switch typ := typi.(type) {
	case FuncType:
		panic("TODO")
	case *NamedType:
		if thisPackagePath == typ.PackagePath {
			return typ.Name, nil
		}
		fmt.Println("[" + typ.Name + "][" + typ.PackagePath + "]")
		return path.Base(typ.PackagePath) + "." + typ.Name, []string{typ.PackagePath}
	case PointerType:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem)
		return "*" + s, addPkgPathList
	case SliceType:
		s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem)
		return "[]" + s, addPkgPathList
	case MapType:
		ks, kaddPkgPathList := MustWriteGoTypes(thisPackagePath, typ.Key)
		vs, vaddPkgPathList := MustWriteGoTypes(thisPackagePath, typ.Value)
		return "map[" + ks + "]" + vs, append(kaddPkgPathList, vaddPkgPathList...)
	case BuiltinType:
		return string(typ), nil
	/*
		case *types.Basic:
			return typ.String(), nil
		case *types.Named:
			if typ.Obj().Pkg() == nil {
				return typ.Obj().Name(), nil
			}
			typPkgPath := typ.Obj().Pkg().Path()
			if thisPackagePath == typPkgPath {
				return typ.Obj().Name(), nil
			}
			return path.Base(typPkgPath) + "." + typ.Obj().Name(), []string{typPkgPath}
		case *types.Pointer:
			s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
			return "*" + s, addPkgPathList
		case *types.Slice:
			s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
			return "[]" + s, addPkgPathList
		case *types.Interface:
			return typ.String(), nil
		//s, addPkgPathList = MustWriteGoTypes(thisPackagePath, typ.Elem())
		//return "[]" + s, addPkgPathList
	*/
	default:
		panic(fmt.Errorf("[MustWriteGoTypes] Not implement go/types [%T]",
			typi))
	}
	return "", nil
}

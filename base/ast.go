package base

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

//获取常量对应的注释
func Ast(fileName string) map[int] string{
	names := map[int] string{}
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs(fileName)
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return names
	}
	nLen := len(f.Decls)
	for i := 0; i < nLen; i++{
		constVal := 0
		decl, ok := f.Decls[i].(*ast.GenDecl)
		if ok && decl.Tok == token.CONST{
			for _, v1 := range decl.Specs{
				val, ok := v1.(*ast.ValueSpec)
				if ok {
					constVal++
					if len(val.Values) >= 1{
						val, ok := val.Values[0].(*ast.BasicLit)
						if ok {
							constVal = Int(val.Value)
						}
					}
					names[constVal] = strings.Trim(val.Comment.Text(), "\n")
				}
			}
		}
	}

	return names
}
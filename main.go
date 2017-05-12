// Package channeler generates synced type using channels.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/channeler/utils"
)

var name = "channeler"
var version = "0.0.0"

func main() {

	var help bool
	var h bool
	var ver bool
	var v bool
	var outPkg string
	flag.BoolVar(&help, "help", false, "Show help.")
	flag.BoolVar(&h, "h", false, "Show help.")
	flag.BoolVar(&ver, "version", false, "Show version.")
	flag.BoolVar(&v, "v", false, "Show version.")
	flag.StringVar(&outPkg, "p", "", "Package name of the new code.")

	flag.Parse()

	if ver || v {
		showVer()
		return
	}
	if help || h {
		showHelp()
		return
	}

	if flag.NArg() < 1 {
		panic("wrong usage")
	}
	args := flag.Args()

	out := ""
	if args[0] == "-" {
		args = args[1:]
		out = "-"
	}

	todos, err := utils.NewTransformsArgs(utils.GetPkgToLoad()).Parse(args)
	if err != nil {
		panic(err)
	}

	filesOut := utils.NewFilesOut("github.com/mh-cbon/" + name)

	for _, todo := range todos.Args {
		srcName := todo.FromTypeName
		destName := todo.ToTypeName
		toImport := todo.FromPkgPath
		if todo.FromPkgPath == "" {
			log.Println("Skipped ", srcName)
			continue
		}
		prog := astutil.GetProgramFast(toImport).Package(toImport)
		foundTypes := astutil.FindTypes(prog)
		foundMethods := astutil.FindMethods(prog)
		foundCtors := astutil.FindCtors(prog, foundTypes)

		fileOut := filesOut.Get(todo.ToPath)
		fileOut.PkgName = outPkg

		if fileOut.PkgName == "" {
			fileOut.PkgName = findOutPkg(todo)
		}

		if todo.FromPkgPath != todo.ToPkgPath {
			fileOut.AddImport(todo.FromPkgPath, "")
		}
		if todo.FromPkgPath != todo.ToPkgPath {
			fileOut.AddImport(todo.FromPkgPath, "")
		}

		res := processType(destName, srcName, foundCtors, foundMethods)
		io.Copy(&fileOut.Body, &res)
	}

	filesOut.Write(out)
}

func findOutPkg(todo utils.TransformArg) string {
	if todo.ToPkgPath != "" {
		prog := astutil.GetProgramFast(todo.ToPkgPath)
		if prog != nil {
			pkg := prog.Package(todo.ToPkgPath)
			return pkg.Pkg.Name()
		}
	}
	if todo.ToPkgPath == "" {
		prog := astutil.GetProgramFast(utils.GetPkgToLoad())
		if len(prog.Imported) < 1 {
			panic("impossible, add [-p name] option")
		}
		for _, p := range prog.Imported {
			return p.Pkg.Name()
		}
	}
	if strings.Index(todo.ToPkgPath, "/") > -1 {
		return filepath.Base(todo.ToPkgPath)
	}
	return todo.ToPkgPath
}

func showVer() {
	fmt.Printf("%v %v\n", name, version)
}

func showHelp() {
	showVer()
	fmt.Println()
	fmt.Println("Usage")
	fmt.Println()
	fmt.Printf("  %v [-p name] [...types]\n\n", name)
	fmt.Printf("  types:  A list of types such as src:dst.\n")
	fmt.Printf("          A type is defined by its package path and its type name,\n")
	fmt.Printf("          [pkgpath/]name\n")
	fmt.Printf("          If the Package path is empty, it is set to the package name being generated.\n")
	// fmt.Printf("          If the Package path is a directory relative to the cwd, and the Package name is not provided\n")
	// fmt.Printf("          the package path is set to this relative directory,\n")
	// fmt.Printf("          the package name is set to the name of this directory.\n")
	fmt.Printf("          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName \n")
	fmt.Printf("  -p:     The name of the package output.\n")
	fmt.Println()
}

func processType(destName, srcName string, foundCtors map[string]*ast.FuncDecl, foundMethods map[string][]*ast.FuncDecl) bytes.Buffer {

	var b bytes.Buffer
	dest := &b

	fmt.Fprintf(dest, `
// %v is channeled.
type %v struct{
	embed %v
	ops chan func()
	stop chan bool
	tick chan bool
}
		`, destName, destName, srcName)

	ctorParams := ""
	ctorParamsInvokation := ""
	ctorName := ""
	ctorIsPointer := false
	srcConcrete := astutil.GetUnpointedType(srcName)
	if x, ok := foundCtors[srcConcrete]; ok {
		withEllipse := astutil.MethodHasEllipse(x)
		ctorParamsInvokation = astutil.MethodParamNamesInvokation(x, withEllipse)
		ctorParams = astutil.MethodParams(x)
		ctorIsPointer = astutil.MethodReturnPointer(x)
		ctorName = "New" + srcConcrete
	}

	if !(astutil.IsAPointedType(srcName) == ctorIsPointer) {
		ctorParams = ""
	}

	fmt.Fprintf(dest, `// New%v constructs a channeled version of %v
func New%v(%v) *%v {
	ret := &%v{
		ops: make(chan func()),
		tick: make(chan bool),
		stop: make(chan bool),
	}
`,
		destName, srcName, destName, ctorParams, destName, destName)

	if ctorName != "" && astutil.IsAPointedType(srcName) == ctorIsPointer {
		fmt.Fprintf(dest, "	ret.embed = %v(%v)\n", ctorName, ctorParamsInvokation)
	}
	fmt.Fprintf(dest, "	go ret.Start()\n")
	fmt.Fprintf(dest, "	return ret\n")
	fmt.Fprintf(dest, "}\n")

	fmt.Fprintf(dest, `// Start the main loop
func (t *%v) Start(){
	for {
		select{
		case op:=<-t.ops:
			op()
			t.tick<-true
		case <-t.stop:
			return
		}
	}
}`,
		destName)

	fmt.Fprintln(dest)
	fmt.Fprintf(dest, `// Stop the main loop
func (t *%v) Stop(){
	t.stop <- true
}`,
		destName)
	fmt.Fprintln(dest)

	for _, m := range foundMethods[srcConcrete] {
		withEllipse := astutil.MethodHasEllipse(m)
		paramNames := astutil.MethodParamNamesInvokation(m, withEllipse)
		// receiverName := astutil.ReceiverName(m)
		methodName := astutil.MethodName(m)
		varExpr := ""
		assignExpr := ""
		callExpr := fmt.Sprintf("%v.embed.%v(%v)", "t", methodName, paramNames)
		returnExpr := ""
		methodReturnTypes := astutil.MethodReturnTypes(m)
		if len(methodReturnTypes) > 0 {
			retVars := astutil.MethodReturnVars(m)
			varExpr = fmt.Sprintf("var ")
			for i, r := range retVars {
				varExpr += r + " " + methodReturnTypes[i] + ", "
			}
			varExpr = varExpr[:len(varExpr)-2]
			assignExpr = fmt.Sprintf("%v = ", strings.Join(retVars, ", "))
			returnExpr = fmt.Sprintf(`
				return %v
				`, strings.Join(retVars, ", "))
		}
		sExpr := fmt.Sprintf(`
	%v
	t.ops<-func() {%v%v}
	<-t.tick
	%v

`, varExpr, assignExpr, callExpr, returnExpr)

		sExpr = fmt.Sprintf(`func(){%v}`, sExpr)
		expr, err := parser.ParseExpr(sExpr)
		if err != nil {
			panic(err)
		}
		astutil.SetReceiverName(m, "t")
		astutil.SetReceiverTypeName(m, destName)
		astutil.SetReceiverPointer(m, true)
		m.Body = expr.(*ast.FuncLit).Body
		fmt.Fprintf(dest, "// %v is channeled\n", methodName)
		m.Doc = nil // clear the doc.
		fmt.Fprintf(dest, "%v\n", astutil.Print(m))
	}

	return b
}

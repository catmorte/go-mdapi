package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/catmorte/go-mdapi/internal/converters"
	"github.com/catmorte/go-mdapi/internal/file"
	"github.com/catmorte/go-mdapi/internal/parser"
	"github.com/catmorte/go-mdapi/internal/types"
	varsPkg "github.com/catmorte/go-mdapi/internal/vars"
	"github.com/spf13/cobra"
)

var (
	mdPath  string
	vars    = map[string]string{}
	cfgPath string

	resultFolder = ".result"
)

func assert(err error, s string, args ...any) {
	if err != nil {
		fmt.Println(fmt.Sprintf(s, args...), err)
		os.Exit(1)
	}
}

func assertOK(ok bool, s string, args ...any) {
	if !ok {
		fmt.Println(fmt.Sprintf(s, args...))
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-mdapi",
	Short: "go-mdapi is a sample CLI application to call api declared in structured md file",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("go-mdapi is a sample CLI application to call api declared in structured md file. use --help for detail")
		fmt.Println()
		fmt.Println("each template supports the following fields")
		fmt.Println(" - CURDIR - current directory")
		fmt.Println(" - CURFILE - current file w/o extension")
		fmt.Println(" - RESULTDIR - result directory")
		fmt.Println()
		fmt.Println("each var supports the following filters")
		for _, v := range converters.SupportedConvs() {
			fmt.Println(" - " + v)
		}
		fmt.Println()
		fmt.Println("internal types")
		for _, v := range types.InternalTypes() {
			fmt.Println(" - " + v.GetName())
		}
		fmt.Println()

	},
}

var varTypesCmd = &cobra.Command{
	Use:   "var_types",
	Short: "returns all available var types",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lenArgs := len(args)
		switch lenArgs {
		case 0:
			for _, v := range file.GetSupportedTypes() {
				fmt.Println(v)
			}
		default:
			c, err := file.GetTypeDescription(args[0])
			assert(err, "failed to get type description")
			fmt.Println(c)
		}
	},
}

var typeVarsCmd = &cobra.Command{
	Use:   "type_vars",
	Short: "returns all possible type's vars",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dts, err := types.GetDefinedTypes(cfgPath)
		assert(err, "failed to get defined types")
		dt, err := dts.FindByName(args[0])
		for _, v := range dt.GetVars() {
			fmt.Println(v)
		}
	},
}

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "returns all available types declared in $HOME/.config/go-mdapi folder",
	Run: func(cmd *cobra.Command, args []string) {
		definedTypes, err := types.GetDefinedTypes(cfgPath)
		assert(err, "can't get defined types")
		for _, v := range definedTypes {
			fmt.Println(v.GetName())
		}
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate api of type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dts, err := types.GetDefinedTypes(cfgPath)
		assert(err, "failed to get defined types")
		dt, err := dts.FindByName(args[0])
		assert(err, "failed to get defined type")
		fmt.Println(dt.NewAPI())
	},
}

var varsCmd = &cobra.Command{
	Use:   "vars [var_name] [index]",
	Short: "shows all the vars in format name:type:count",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		fileData, err := parser.ParseMarkdownFile(mdPath)
		assert(err, "failed to open file")
		lenArgs := len(args)
		switch lenArgs {
		case 0:
			for _, v := range fileData.Vars {
				fmt.Printf("%s:%s:%d", v.Nam, v.Typ, len(v.Vals))
				fmt.Println()
			}
		default:
			c, ok := fileData.GetVarByName(args[0])
			assertOK(ok, "unknown var")
			switch lenArgs {
			case 1:
				fmt.Printf("has values: %d", len(c.Vals))
			case 2:
				index, err := strconv.Atoi(args[1])
				assert(err, "failed to parse index")
				assertOK((index >= 0) && (index < len(c.Vals)), "index out of bounds")
				fmt.Println(c.Vals[index].Typ)
				fmt.Println(c.Vals[index].Val)
			}
		}
	},
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compile the api",
	Args:  cobra.MaximumNArgs(1), // Allow at most 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		fileData, err := parser.ParseMarkdownFile(mdPath)
		assert(err, "failed to parse file")
		curdir := filepath.Dir(mdPath)
		curfile := strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))
		resdir := filepath.Join(curdir, resultFolder, curfile)
		if vars == nil {
			vars = map[string]string{}
		}
		allFields := varsPkg.Vars(vars)
		allFields.SetCurrentDir(curdir)
		allFields.SetCurrentFile(curfile)
		allFields.SetResultDir(resdir)
		allFields, err = fileData.Compute(vars)
		assert(err, "failed to compute")
		dts, err := types.GetDefinedTypes(cfgPath)
		assert(err, "failed to get defined types")
		dt, err := dts.FindByName(fileData.Typ.Typ)
		assert(err, "failed to get defined type")
		err = fileData.Typ.Fields.Compute(allFields, true)
		assert(err, "failed to parse type fields")

		err = dt.Compile(allFields)
		assert(err, "failed to run")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the api",
	Args:  cobra.MaximumNArgs(1), // Allow at most 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		fileData, err := parser.ParseMarkdownFile(mdPath)
		assert(err, "failed to parse file")
		curdir := filepath.Dir(mdPath)
		curfile := strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))
		resdir := filepath.Join(curdir, resultFolder, curfile)
		if vars == nil {
			vars = map[string]string{}
		}
		allFields := varsPkg.Vars(vars)
		allFields.SetCurrentDir(curdir)
		allFields.SetCurrentFile(curfile)
		allFields.SetResultDir(resdir)
		allFields, err = fileData.Compute(vars)
		assert(err, "failed to compute")
		dts, err := types.GetDefinedTypes(cfgPath)
		assert(err, "failed to get defined types")
		dt, err := dts.FindByName(fileData.Typ.Typ)
		assert(err, "failed to get defined type")
		err = fileData.Typ.Fields.Compute(allFields, true)
		assert(err, "failed to parse type fields")
		_, err = os.Stat(resdir)
		if !os.IsNotExist(err) {
			counter := 1
			var newPath string
			for {
				newPath = filepath.Join(curdir, resultFolder, fmt.Sprintf("%s_%d", curfile, counter))
				_, err = os.Stat(newPath)
				if os.IsNotExist(err) {
					err = os.Rename(resdir, newPath)
					assert(err, "failed to rename")
					break
				}
				counter++
			}
		}

		err = os.MkdirAll(resdir, 0o755)
		assert(err, "failed to create result dir")
		err = dt.Run(allFields)
		assert(err, "failed to run")
		varsFile := filepath.Join(resdir, ".vars")
		jsonVarsRaw, err := json.MarshalIndent(allFields, "", " ")
		assert(err, "failed to convert fields to json")
		err = os.WriteFile(varsFile, jsonVarsRaw, 0x775)
		assert(err, "failed to write vars")

		err = fileData.After.Compute(allFields, true)
		assert(err, "failed to compute after")
		for _, v := range fileData.After {
			afterField := filepath.Join(resdir, v.Nam)
			err = os.WriteFile(afterField, []byte(allFields[v.Nam]), 0x775)
			assert(err, "failed to write %s", v.Nam)
		}
		fmt.Println(allFields.GetResultDir())
	},
}

func defineFileFlag(c *cobra.Command) {
	c.PersistentFlags().StringVarP(&mdPath, "file", "f", "", "path to the file to read (required)")
	c.MarkPersistentFlagRequired("file")
}

func main() {
	defineFileFlag(varsCmd)
	defineFileFlag(runCmd)
	defineFileFlag(compileCmd)
	runCmd.Flags().StringToStringVar(&vars, "vars", nil, "key-value parameters (e.g. --vars key1=value1 --vars key2=value2)")
	compileCmd.Flags().StringToStringVar(&vars, "vars", nil, "key-value parameters (e.g. --vars key1=value1 --vars key2=value2)")
	rootCmd.AddCommand(varsCmd)
	rootCmd.AddCommand(typesCmd)
	rootCmd.AddCommand(varTypesCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(typeVarsCmd)

	dirname, err := os.UserHomeDir()
	assert(err, "can't get user's home dir")

	cfgPath = filepath.Join(dirname, ".config", "go-mdapi")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

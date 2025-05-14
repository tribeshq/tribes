package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/abigen"
)

const rollupsContractsUrl = "https://registry.npmjs.org/@cartesi/rollups/-/rollups-2.0.0-rc.18.tgz"
const baseContractsPath = "package/out/"
const bindingPkg = "rollups_contracts"

type contractBinding struct {
	jsonPath        string
	custom_typeName string
	outFile         string
}

var bindings = []contractBinding{
	{
		jsonPath:        baseContractsPath + "IInputBox.sol/IInputBox.json",
		custom_typeName: "IInputBox",
		outFile:         "iinputbox.go",
	},
	{
		jsonPath:        baseContractsPath + "IApplication.sol/IApplication.json",
		custom_typeName: "IApplication",
		outFile:         "iapplication.go",
	},
	{
		jsonPath:        baseContractsPath + "IERC20Portal.sol/IERC20Portal.json",
		custom_typeName: "IERC20Portal",
		outFile:         "ierc20portal.go",
	},
}

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", cwd)

	contractsZip := downloadContracts(rollupsContractsUrl)
	defer contractsZip.Close()
	contractsTar := unzip(contractsZip)
	defer contractsTar.Close()

	files := make(map[string]bool)
	for _, b := range bindings {
		files[b.jsonPath] = true
	}
	contents := readFilesFromTar(contractsTar, files)

	for _, b := range bindings {
		content := contents[b.jsonPath]
		if content == nil {
			log.Fatal("missing contents for ", b.jsonPath)
		}
		generateBinding(b, content)
	}
}

// Exit if there is any error.
func checkErr(context string, err any) {
	if err != nil {
		log.Fatal(context, ": ", err)
	}
}

// Download the contracts from rollupsContractsUrl.
// Return the buffer with the contracts.
func downloadContracts(url string) io.ReadCloser {
	log.Print("downloading contracts from ", url)
	response, err := http.Get(url)
	checkErr("download tgz", err)
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		log.Fatal("invalid status: ", response.Status)
	}
	return response.Body
}

// Decompress the buffer with the contracts.
func unzip(r io.Reader) io.ReadCloser {
	log.Print("unziping contracts")
	gzipReader, err := gzip.NewReader(r)
	checkErr("unziping", err)
	return gzipReader
}

// Read the required files from the tar.
// Return a map with the file contents.
func readFilesFromTar(r io.Reader, files map[string]bool) map[string][]byte {
	contents := make(map[string][]byte)
	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		checkErr("read tar", err)
		if files[header.Name] {
			contents[header.Name], err = io.ReadAll(tarReader)
			checkErr("read tar", err)
		}
	}
	return contents
}

// Get the .abi key from the json
func getAbi(rawJson []byte) []byte {
	var contents struct {
		Abi json.RawMessage `json:"abi"`
	}
	err := json.Unmarshal(rawJson, &contents)
	checkErr("decode json", err)
	return contents.Abi
}

// Generate the Go bindings for the contracts.
func generateBinding(b contractBinding, content []byte) {
	var (
		sigs         []map[string]string
		abis         = []string{string(getAbi(content))}
		bins         = []string{""}
		custom_types = []string{b.custom_typeName}
		libs         = make(map[string]string)
		aliases      = make(map[string]string)
	)
	code, err := abigen.Bind(custom_types, abis, bins, sigs, bindingPkg, libs, aliases)
	checkErr("generate binding", err)

	// Get the absolute path for the output file
	absPath, err := filepath.Abs(b.outFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path for %s: %v", b.outFile, err)
	}
	log.Printf("Generating binding to: %s", absPath)

	// Ensure the output directory exists
	dirPath := filepath.Dir(absPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	checkErr("create directory '"+dirPath+"' for output file '"+absPath+"'", err)

	const fileMode = 0600
	err = os.WriteFile(absPath, []byte(code), fileMode)
	checkErr("write binding file '"+absPath+"'", err)
	log.Print("generated binding ", absPath)
}

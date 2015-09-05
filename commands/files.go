package commands

import (
	"github.com/eris-ltd/eris-cli/files"

	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/spf13/cobra"
)

// Primary Files Sub-Command
// Flags to add: ipfsHost
var Files = &cobra.Command{
	Use:   "files",
	Short: "Manage Files containers for your Application.",
	Long: `The files subcommand is used to import, and export
files to and from IPFS for use on the host machine.

These commands are provided in addition to the various
functionality which is included throughout the tool, such as
services import or services export which operate more
precisely. The eris files command is used as a general wrapper
around an IPFS gateway which would be running as eris services ipfs.

At times, due to the manner in which IPFS boots files commands
will fail. If you get errors when running eris files commands
then please run [eris services start ipfs] give that a second
or two to boot and then retry the eris files command which failed.`,
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

// build the files subcommand
func buildFilesCommand() {
	Files.AddCommand(filesImport)
	Files.AddCommand(filesExport)
	Files.AddCommand(filesCache)
	Files.AddCommand(filesCat)
	Files.AddCommand(filesList)
	Files.AddCommand(filesCached)
	addFilesFlags()
}

var filesImport = &cobra.Command{
	Use:   "get [hash] [fileName]",
	Short: "Pull files from IPFS via a hash and save them locally.",
	Long: `Pull files from IPFS via a hash and save them locally.

Optionally pass in a csv with: get --csv=[fileName]`,
	Run: FilesGet,
}

var filesExport = &cobra.Command{
	Use:   "put [fileName]",
	Short: "Post files to IPFS.",
	Long: `Post files to IPFS.

Optionally post all contents of a directory with: put [dirName] --dir`,
	Run: FilesPut,
}

var filesCache = &cobra.Command{
	Use:   "cache [fileHash]",
	Short: "Cache files to IPFS.",
	Long: `Cache files to IPFS' local daemon.

Caches a files locally via IPFS pin, by hash.
Optionally pass in a csv with: cache --csv=[fileName]
Note: "put" will "cache" recursively by default`,
	Run: FilesPin,
}

var filesCat = &cobra.Command{
	Use:   "cat [fileHash]",
	Short: "Cat the contents of a file from IPFS.",
	Long:  "Cat the contents of a file from IPFS.",
	Run:   FilesCat,
}

var filesList = &cobra.Command{
	Use:   "ls [objectHash]",
	Short: "List links from an IPFS object.",
	//TODO test listing up and down through DAG / Zach just learn the DAG.
	Long: "Lists object named by [objectHash/Path] and displays the link it contains.",
	Run:  FilesList,
}

var filesCached = &cobra.Command{
	Use:   "cached",
	Short: "Lists files cached locally.",
	Long:  `Displays list of files cached locally.`,
	Run:   FilesManageCached,
}

//--------------------------------------------------------------
// cli flags
func addFilesFlags() {

	filesImport.Flags().StringVarP(&do.CSV, "csv", "", "", "specify a .csv with entries of format: hash,fileName")
	filesImport.Flags().StringVarP(&do.NewName, "dirname", "", "", "name of new directory to dump IPFS files from --csv")
	filesExport.Flags().StringVarP(&do.Gateway, "gateway", "", "", "specify a hosted gateway. default is IPFS' gateway; type \"eris\" for our gateway, or use your own with \"http://yourhost\"")
	//TODO `put files --dir -r` once pr to ipfs is merged
	filesExport.Flags().BoolVarP(&do.AddDir, "dir", "", false, "add all files from a directory (note: this will not create an ipfs object). returns a log file (ipfs_hashes.csv) to pass into `eris files get`")

	//command will ignore fileName but that's ok
	filesCache.Flags().StringVarP(&do.CSV, "csv", "", "", "specify a .csv with entries of format: hash,fileName")

	filesCached.Flags().BoolVarP(&do.Rm, "rma", "", false, "remove all cached files")
	filesCached.Flags().StringVarP(&do.Hash, "rm", "", "", "remove a cached file by hash")
}

func FilesGet(cmd *cobra.Command, args []string) {
	if do.CSV == "" {
		IfExit(ArgCheck(2, "eq", cmd, args))
		do.Name = args[0]
		do.Path = args[1]
	} else {
		do.Name = ""
		do.Path = ""
	}
	IfExit(files.GetFiles(do))
}

func FilesPut(cmd *cobra.Command, args []string) {
	IfExit(ArgCheck(1, "eq", cmd, args))

	do.Name = args[0]
	err := files.PutFiles(do)
	IfExit(err)
	logger.Println(do.Result)
}

func FilesPin(cmd *cobra.Command, args []string) {
	if do.CSV == "" {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		do.Name = args[0]
	} else {
		do.Name = ""
	}
	err := files.PinFiles(do)
	IfExit(err)
	logger.Println(do.Result)
}

func FilesCat(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}
	do.Name = args[0]
	err := files.CatFiles(do)
	IfExit(err)
	logger.Println(do.Result)

}

func FilesList(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}
	do.Name = args[0]
	err := files.ListFiles(do)
	IfExit(err)
	logger.Println(do.Result)
}

func FilesManageCached(cmd *cobra.Command, args []string) {
	err := files.ManagePinned(do)
	IfExit(err)
	logger.Println(do.Result)
}

package export

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

const (
	OutputTypeSingleFile = "single-file"
	OutputTypeDirectory  = "dir"
)

func writeExportedIntents(files []cloudapi.ClientIntentsFileRepresentation, outputLocation string, outputType string) error {
	if outputLocation == "" {
		printIntentFilesToStdout(files)
		return nil
	}

	switch outputType {
	case OutputTypeSingleFile:
		err := writeIntentsFiles(outputLocation, files)
		if err != nil {
			return err
		}
		output.PrintStderr("Successfully wrote intents into %s", outputLocation)
	case OutputTypeDirectory:
		files = filterDuplicateFilenames(files)

		if len(files) == 0 {
			output.PrintStderr("No intent files to write.")
			return nil
		}

		err := os.MkdirAll(outputLocation, 0700)
		if err != nil {
			return fmt.Errorf("could not create dir %s: %w", outputLocation, err)
		}

		for _, file := range files {
			filePath := filepath.Join(outputLocation, file.FileName)
			err := writeIntentsFiles(filePath, []cloudapi.ClientIntentsFileRepresentation{file})
			if err != nil {
				return err
			}
		}
		output.PrintStderr("Successfully wrote intents into %s", outputLocation)
	default:
		return fmt.Errorf("unexpected output type %s, use one of (%s, %s)", outputType, OutputTypeSingleFile, OutputTypeDirectory)
	}

	return nil
}

func filterDuplicateFilenames(files []cloudapi.ClientIntentsFileRepresentation) []cloudapi.ClientIntentsFileRepresentation {
	filesByFilename := lo.GroupBy(files, func(file cloudapi.ClientIntentsFileRepresentation) string {
		return file.FileName
	})

	hasUniqueFilename := func(file cloudapi.ClientIntentsFileRepresentation, _ int) bool {
		return len(filesByFilename[file.FileName]) == 1
	}
	filesWithUniqueFilename := lo.Filter(files, hasUniqueFilename)
	filesWithDuplicateFilenames := lo.Reject(files, hasUniqueFilename)

	if len(filesWithDuplicateFilenames) > 0 {
		duplicateFilenames := lo.Uniq(lo.Map(filesWithDuplicateFilenames, func(file cloudapi.ClientIntentsFileRepresentation, _ int) string {
			return file.FileName
		}))
		output.PrintStderr("Duplicate filenames detected, omitting the following files:")
		for _, filename := range duplicateFilenames {
			output.PrintStderr("- %s", filename)
		}

		output.PrintStderr("Consider filtering by clusters, namespaces, or services to avoid duplicate filenames.")
	}

	return filesWithUniqueFilename
}

func printIntentFilesToStdout(files []cloudapi.ClientIntentsFileRepresentation) {
	formatted := output.FormatClientIntentsFiles(files)
	output.PrintStdout(formatted)
}

func writeIntentsFiles(filePath string, files []cloudapi.ClientIntentsFileRepresentation) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	formatted := output.FormatClientIntentsFiles(files)
	_, err = f.WriteString(formatted)
	if err != nil {
		return err
	}
	return nil
}

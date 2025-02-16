package export

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"os"
	"path/filepath"
)

const (
	OutputTypeSingleFile = "single-file"
	OutputTypeDirectory  = "dir"
)

type IntentsWriter struct {
	OutputLocation   string
	OutputType       string
	WithDiffComments bool
}

func NewIntentsWriter(outputLocation string, outputType string, withDiffComments bool) *IntentsWriter {
	return &IntentsWriter{
		OutputLocation:   outputLocation,
		OutputType:       outputType,
		WithDiffComments: withDiffComments,
	}
}

func (w *IntentsWriter) WriteExportedIntents(files []cloudapi.ClientIntentsFileRepresentation) error {
	files = w.filterDuplicateFilenames(files)
	if len(files) == 0 {
		prints.PrintCliStderr("No intent files to write.")
		return nil
	}

	if w.OutputLocation == "" {
		w.printIntentFilesToStdout(files)
		return nil
	}

	switch w.OutputType {
	case OutputTypeSingleFile:
		err := w.writeIntentsToFile(w.OutputLocation, files)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("Successfully wrote intents into %s", w.OutputLocation)
	case OutputTypeDirectory:
		err := w.writeIntentsToDir(w.OutputLocation, files)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Successfully wrote intents into %s", w.OutputLocation)
	default:
		return fmt.Errorf("unexpected output type %s, use one of (%s, %s)", w.OutputType, OutputTypeSingleFile, OutputTypeDirectory)
	}

	return nil
}

func (w *IntentsWriter) filterDuplicateFilenames(files []cloudapi.ClientIntentsFileRepresentation) []cloudapi.ClientIntentsFileRepresentation {
	filesByFilename := lo.GroupBy(files, func(file cloudapi.ClientIntentsFileRepresentation) string {
		return file.NamespacedFileName
	})

	hasUniqueFilename := func(file cloudapi.ClientIntentsFileRepresentation, _ int) bool {
		return len(filesByFilename[file.NamespacedFileName]) == 1
	}
	filesWithUniqueFilename := lo.Filter(files, hasUniqueFilename)
	filesWithDuplicateFilenames := lo.Reject(files, hasUniqueFilename)

	if len(filesWithDuplicateFilenames) > 0 {
		duplicateFilenames := lo.Uniq(lo.Map(filesWithDuplicateFilenames, func(file cloudapi.ClientIntentsFileRepresentation, _ int) string {
			return file.NamespacedFileName
		}))
		prints.PrintCliStderr("Duplicate filenames detected, omitting the following files:")
		for _, filename := range duplicateFilenames {
			prints.PrintCliStderr("- %s", filename)
		}

		prints.PrintCliStderr("Consider filtering by clusters, namespaces, or services to avoid duplicate filenames.")
	}

	return filesWithUniqueFilename
}

func (w *IntentsWriter) printIntentFilesToStdout(files []cloudapi.ClientIntentsFileRepresentation) {
	formatted := output.FormatClientIntentsFiles(files, w.WithDiffComments)
	prints.PrintCliOutput(formatted)
}

func (w *IntentsWriter) writeIntentsToDir(dirPath string, files []cloudapi.ClientIntentsFileRepresentation) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		return fmt.Errorf("could not create dir %s: %w", dirPath, err)
	}

	for _, file := range files {
		filePath := filepath.Join(dirPath, file.NamespacedFileName)
		err := w.writeIntentsToFile(filePath, []cloudapi.ClientIntentsFileRepresentation{file})
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *IntentsWriter) writeIntentsToFile(filePath string, files []cloudapi.ClientIntentsFileRepresentation) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	formatted := output.FormatClientIntentsFiles(files, w.WithDiffComments)
	_, err = f.WriteString(formatted)
	if err != nil {
		return err
	}
	return nil
}

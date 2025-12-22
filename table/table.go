package table

import (
	"io"

	"github.com/loft-sh/log"
	"github.com/loft-sh/log/scanner"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/sirupsen/logrus"
)

func PrintTable(s log.Logger, header []string, values [][]string) {
	PrintTableWithOptions(s, header, values, nil)
}

// PrintTableWithOptions prints a table with header columns and string values
func PrintTableWithOptions(s log.Logger, header []string, values [][]string, modify func(table *tablewriter.Table)) {
	reader, writer := io.Pipe()
	defer writer.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		sa := scanner.NewScanner(reader)
		for sa.Scan() {
			s.WriteString(logrus.InfoLevel, "  "+sa.Text()+"\n")
		}
	}()

	alignment := make(tw.Alignment, len(header))
	for i := range alignment {
		alignment[i] = tw.AlignLeft
	}

	opts := []tablewriter.Option{
		tablewriter.WithAlignment(alignment),
	}

	table := tablewriter.NewTable(writer, opts...)

	if modify != nil {
		modify(table)
	}

	// Render
	_, _ = writer.Write([]byte("\n"))
	table.Header(header)
	table.Bulk(values)
	_ = table.Render()
	_, _ = writer.Write([]byte("\n"))
	_ = writer.Close()
	<-done
}

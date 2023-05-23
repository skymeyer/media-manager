package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"go.skymeyer.dev/media-manager/pkg/downloader"
)

func main() {
	if err := New().Execute(); err != nil {
		os.Exit(1)
	}
}

func New() *cobra.Command {

	app := &cobra.Command{
		Use:   "mm",
		Short: "Media Manager",
	}

	app.AddCommand(
		NewDownloadCmd(),
	)

	return app
}

func NewDownloadCmd() *cobra.Command {

	var (
		outputDir = "/home/runner/download"
		urls      []string
		fromFile  string
	)

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download media",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Initialize downloader
			dnl, err := downloader.New(
				downloader.WithWriters(cmd.OutOrStdout(), cmd.ErrOrStderr()),
				downloader.WithOutputDir(outputDir),
				downloader.WithYTArguments([]string{
					"--write-description", // dump additional info for tag correction
				}),
			)
			if err != nil {
				return err
			}

			// Read download URLs from file if given
			if fromFile != "" {
				file, err := os.Open(fromFile)
				if err != nil {
					return err
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					urls = append(urls, scanner.Text())
				}

				if err := scanner.Err(); err != nil {
					return err
				}
			}

			// Iterate over URLs
			for _, u := range urls {
				target, err := url.Parse(u)
				if err != nil {
					return fmt.Errorf("invalid url %q", u)
				}
				fmt.Printf("Downloading %s ...\n", u)
				if err := dnl.Download(target); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&urls, "urls", urls, "List of URLs to download")
	cmd.Flags().StringVar(&fromFile, "from-file", fromFile, "File containing download URLs")
	cmd.Flags().StringVar(&outputDir, "output-dir", outputDir, "Directory to store downloads")

	return cmd
}

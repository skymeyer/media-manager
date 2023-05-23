package downloader

import (
	"fmt"
	"io"
	"net/url"
	"os/exec"
)

const (
	BIN_YT_DLP = "yt-dlp"
	BIN_FFMPEG = "ffmpeg"
	BIN_FFPLAY = "ffplay"
	BIN_SCDL   = "scdl"
)

var (
	requiredBins = []string{
		BIN_YT_DLP,
		BIN_FFMPEG,
		BIN_FFPLAY,
		BIN_SCDL,
	}
)

type Option func(*Downloader)

func New(opts ...Option) (*Downloader, error) {

	// Verify third party binaries are available
	for _, bin := range requiredBins {
		if _, err := exec.LookPath(bin); err != nil {
			return nil, fmt.Errorf("%s not found", bin)
		}
	}

	// Create default downloader
	dl := &Downloader{
		ytFormat: "bestaudio[ext=m4a]",
	}

	// Apply options
	for _, opt := range opts {
		opt(dl)
	}

	return dl, nil
}

type Downloader struct {
	outputDir string
	ytFormat  string
	ytArgs    []string
	scArgs    []string
	outWriter io.Writer
	errWriter io.Writer
}

func (d *Downloader) Download(target *url.URL) error {
	switch target.Host {
	case "youtube.com", "www.youtube.com", "music.youtube.com":
		return d.downloadYouTube(target)
	case "soundcloud.com":
		return d.downloadSoundCloud(target)
	default:
		return fmt.Errorf("cannot handle url for host %q", target.Host)
	}
}

// downloads files with metadata, cropped album art with filename as Artist - Title.m4a
//
// yt-dlp.exe playlist_url_or_video_url_here -f "bestaudio[ext=m4a]" --embed-thumbnail --convert-thumbnail jpg
//
//	--exec-before-download "ffmpeg -i %(thumbnails.-1.filepath)q -vf crop=\"'if(gt(ih,iw),iw,ih)':'if(gt(iw,ih),ih,iw)'\" _%(thumbnails.-1.filepath)q"
//	--exec-before-download "rm %(thumbnails.-1.filepath)q"
//	--exec-before-download "mv _%(thumbnails.-1.filepath)q %(thumbnails.-1.filepath)q"
//	--output "%(artist)s - %(title)s.%(ext)s"
func (d *Downloader) downloadYouTube(target *url.URL) error {
	var cmdArgs = []string{
		target.String(),
		"--no-update",
		"--no-colors",
		"--format", d.ytFormat,
		"--embed-thumbnail",
		"--convert-thumbnail", "jpg",
		"--output", "%(artist)s - %(title)s.%(ext)s",
	}
	return d.runCmd(exec.Command(BIN_YT_DLP, append(cmdArgs, d.ytArgs...)...))
}

func (d *Downloader) downloadSoundCloud(target *url.URL) error {
	var cmdArgs = []string{
		"-l", target.String(),
		"--addtofile",      // Add artist to filename if missing
		"--overwrite",      // Overwrite file if it already exists
		"--extract-artist", // Set artist tag from title instead of username
		"--original-art",   // Download original cover art
		"--original-name",  // Do not change name of original file downloads
		"--debug",          // Set log level to DEBUG
	}
	return d.runCmd(exec.Command(BIN_SCDL, append(cmdArgs, d.scArgs...)...))
}

func (d *Downloader) runCmd(cmd *exec.Cmd) error {
	cmd.Dir = d.outputDir
	cmd.Stdout = d.outWriter
	cmd.Stderr = d.errWriter

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

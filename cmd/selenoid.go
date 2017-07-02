package cmd

import (
	"fmt"
	"github.com/aerokube/cm/selenoid"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	registryUrl            = "https://registry.hub.docker.com"
	defaultBrowsersJsonURL = "https://raw.githubusercontent.com/aerokube/cm/master/browsers.json"
)

var (
	lastVersions    int
	tmpfs           int
	operatingSystem string
	arch            string
	version         string
	browsers        string
	browsersJSONUrl string
	configDir       string
	skipDownload    bool
	vnc             bool
	force           bool
	args            string
)

func init() {
	initFlags()

	selenoidCmd.AddCommand(selenoidDownloadCmd)
	selenoidCmd.AddCommand(selenoidConfigureCmd)
	selenoidCmd.AddCommand(selenoidStartCmd)
	selenoidCmd.AddCommand(selenoidStopCmd)
	selenoidCmd.AddCommand(selenoidUpdateCmd)
	selenoidCmd.AddCommand(selenoidCleanupCmd)
	selenoidCmd.AddCommand(selenoidStatusCmd)

	selenoidUICmd.AddCommand(selenoidDownloadUICmd)
	selenoidUICmd.AddCommand(selenoidStartUICmd)
	selenoidUICmd.AddCommand(selenoidStopUICmd)
	selenoidUICmd.AddCommand(selenoidUpdateUICmd)
	selenoidUICmd.AddCommand(selenoidCleanupUICmd)
	selenoidUICmd.AddCommand(selenoidUIStatusCmd)
}

func initFlags() {
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidStopCmd,
		selenoidUpdateCmd,
		selenoidCleanupCmd,
		selenoidStatusCmd,
	} {
		c.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidStopCmd,
		selenoidUpdateCmd,
		selenoidCleanupCmd,
		selenoidStatusCmd,
	} {
		c.Flags().StringVarP(&configDir, "config-dir", "c", getSelenoidConfigDir(), "directory to save files")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadUICmd,
		selenoidStartUICmd,
		selenoidStopUICmd,
		selenoidUpdateUICmd,
		selenoidCleanupUICmd,
		selenoidUIStatusCmd,
	} {
		c.Flags().StringVarP(&configDir, "config-dir", "c", getSelenoidUIConfigDir(), "directory to save files")
	}

	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
		selenoidDownloadUICmd,
		selenoidStartUICmd,
		selenoidUpdateUICmd,
	} {
		c.Flags().StringVarP(&operatingSystem, "operating-system", "o", runtime.GOOS, "target operating system (drivers only)")
		c.Flags().StringVarP(&arch, "architecture", "a", runtime.GOARCH, "target architecture (drivers only)")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
	} {
		c.Flags().StringVarP(&version, "version", "v", selenoid.Latest, "desired version; default is latest release")
		c.Flags().StringVarP(&browsers, "browsers", "b", "", "comma separated list of browser names to process")
		c.Flags().StringVarP(&browsersJSONUrl, "browsers-json", "j", defaultBrowsersJsonURL, "browsers JSON data URL (in most cases never need to be set manually)")
		c.Flags().BoolVarP(&skipDownload, "no-download", "n", false, "only output config file without downloading images or drivers")
		c.Flags().StringVarP(&registry, "registry", "r", registryUrl, "Docker registry to use")
		c.Flags().IntVarP(&lastVersions, "last-versions", "l", 2, "process only last N versions (Docker only)")
		c.Flags().IntVarP(&tmpfs, "tmpfs", "t", 0, "add tmpfs volume sized in megabytes (Docker only)")
		c.Flags().BoolVarP(&vnc, "vnc", "s", false, "download containers with VNC support (Docker only)")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidDownloadUICmd,
		selenoidStartUICmd,
	} {
		c.Flags().BoolVarP(&force, "force", "f", false, "force action")
	}
	for _, c := range []*cobra.Command{
		selenoidStartCmd,
		selenoidUpdateCmd,
		selenoidStartUICmd,
		selenoidUpdateUICmd,
	} {
		c.Flags().StringVarP(&args, "args", "g", "", "additional service arguments (e.g. \"-limit 5\")")
	}
}

func createLifecycle() (*selenoid.Lifecycle, error) {
	config := selenoid.LifecycleConfig{
		Quiet:     quiet,
		Force:     force,
		ConfigDir: configDir,
		Browsers:  browsers,
		Download:  !skipDownload,
		Args:      args,

		LastVersions: lastVersions,
		RegistryUrl:  registry,
		Tmpfs:        tmpfs,
		VNC:          vnc,

		BrowsersJsonUrl: browsersJSONUrl,
		OS:              operatingSystem,
		Arch:            arch,
		Version:         version,
	}
	return selenoid.NewLifecycle(&config)
}

var selenoidCmd = &cobra.Command{
	Use:   "selenoid",
	Short: "Download, configure and run Selenoid",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func getConfigDir(elem ...string) string {
	usr, err := user.Current()
	var p string
	if err != nil {
		p = filepath.Join(elem...)
	}
	p = filepath.Join(append([]string{usr.HomeDir}, elem...)...)
	ap, _ := filepath.Abs(p)
	return ap
}

func getSelenoidConfigDir() string {
	return getConfigDir(".aerokube", "selenoid")
}

func getSelenoidUIConfigDir() string {
	return getConfigDir(".aerokube", "selenoid-ui")
}

func stderr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a)
}
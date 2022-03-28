package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/umeier/mktorrent/torrent"
	"net/url"
	"os"
	p "path"
	"path/filepath"
)

var Tracker []string
var BaseUrls []string

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&Tracker, "tracker", "t", []string{}, "Tracker to use")
	_ = rootCmd.MarkFlagRequired("tracker")
	rootCmd.PersistentFlags().StringArrayVarP(&BaseUrls, "base-urls", "b", []string{}, "BaseUrl(s) for webseed")
}

var rootCmd = &cobra.Command{
	Use:   "mktorrent [options] [files]",
	Short: "Create torrent-files for given files",
	Args:  cobra.MinimumNArgs(1),
	Run:   createTorrents,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createTorrents(cmd *cobra.Command, args []string) {
	var ann = []string{}
	var urls = []string{}
	for _, path := range args {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
			os.Exit(1)
		}
		base := filepath.Base(path)
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, tracker := range Tracker {
			ann = append(ann, tracker)
		}

		for _, bu := range BaseUrls {
			u, err := url.Parse(bu)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			u.Path = p.Join(u.Path, base)
			urls = append(urls, u.String())
		}

		tr, err := torrent.MakeTorrent(f, base, ann, urls)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tf, err := os.Create(path + ".torrent")
		defer tf.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tr.Save(tf)

	}
}

package cmd

import (
	"fmt"
	"iht/pkg/cfg"

	"github.com/spf13/cobra"
)

var p115Cmd = &cobra.Command{
	Use:   "p115",
	Short: "pan 115 commands",
	Long:  ``,
}

// 导出115目录树到文件
var p115exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export 115 directory tree to file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}
		fmt.Println(cookie)

		cid, _ := cmd.Flags().GetString("cid")
		fmt.Println(cid)
	},
}

// 将目录树文件生成目录结构
var p115treeifyCmd = &cobra.Command{
	Use:   "treeify",
	Short: "transform 115 directory tree file to directory structure",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}
		fmt.Println(cookie)

		cid, _ := cmd.Flags().GetString("cid")
		fmt.Println(cid)
	},
}

// 自动维护strm媒体库
var p115strmCmd = &cobra.Command{
	Use:   "strm",
	Short: "auto export and maintain strm media library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}
		fmt.Println(cookie)

		cid, _ := cmd.Flags().GetString("cid")
		fmt.Println(cid)
	},
}

func init() {
	p115Cmd.PersistentFlags().StringP("cookie", "c", "", "pan cookie")
	p115Cmd.PersistentFlags().StringP("cid", "i", "", "pan target id")
	p115Cmd.PersistentFlags().StringP("location", "l", "", "pan target path")
	p115Cmd.PersistentFlags().StringP("file", "f", "", "directory tree file path")
	p115Cmd.PersistentFlags().StringP("base", "b", "", "media library directory path")
	p115Cmd.PersistentFlags().StringP("url", "u", "", "media url")
	p115Cmd.PersistentFlags().StringP("prefix", "p", "", "media url prefix path")

	p115Cmd.AddCommand(p115exportCmd)
	p115Cmd.AddCommand(p115treeifyCmd)
	p115Cmd.AddCommand(p115strmCmd)

	rootCmd.AddCommand(p115Cmd)
}

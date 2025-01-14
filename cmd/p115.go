package cmd

import (
	"fmt"
	"iht/pkg/cfg"
	"iht/pkg/pan"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xshrim/gol"
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

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P115.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P115.Cpath
		}

		fpath, _ := cmd.Flags().GetString("fpath")
		if fpath == "" {
			fpath = cfg.Conf.P115.Fpath
		}

		gol.Info("Start to export directory tree")

		client := pan.P115{
			Cookie: cookie,
		}

		if cid == "" {
			if fitem, err := client.FetchItem(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v", err)
				return
			} else {
				if fitem.Fid == "" {
					cid = fitem.Cid
				} else {
					cid = fitem.Fid
				}
			}
		}

		if fpath, err := client.ExportTreeToFile(cid, fpath); err != nil {
			gol.Errorf("Export directory tree failed: %v\n", err)
			return
		} else {
			gol.Infof("Export directory tree to file: %s\n", fpath)
		}

	},
}

// 将目录树文件生成目录结构
var p115treeifyCmd = &cobra.Command{
	Use:   "treeify",
	Short: "transform 115 directory tree file to directory structure",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fpath, _ := cmd.Flags().GetString("fpath")
		if fpath == "" {
			fpath = cfg.Conf.P115.Fpath
		}

		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			url = cfg.Conf.P115.Url
		}

		prefix, _ := cmd.Flags().GetString("prefix")
		if prefix == "" {
			prefix = cfg.Conf.P115.Prefix
		}

		library, _ := cmd.Flags().GetString("library")
		if library == "" {
			library = cfg.Conf.P115.Library
		}

		gol.Info("Start to convert directory tree to library")

		if err := pan.Treeify(fpath, url, prefix, library); err != nil {
			gol.Errorf("Convert directory tree to library failed: %v\n", err)
		} else {
			gol.Info("Convert directory tree to library success")
		}

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

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P115.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P115.Cpath
		}

		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			url = cfg.Conf.P115.Url
		}

		prefix, _ := cmd.Flags().GetString("prefix")
		if prefix == "" {
			prefix = cfg.Conf.P115.Prefix
		}

		library, _ := cmd.Flags().GetString("library")
		if library == "" {
			library = cfg.Conf.P115.Library
		}

		gol.Info("Start to export directory tree")

		client := pan.P115{
			Cookie: cookie,
		}

		if cid == "" {
			if fitem, err := client.FetchItem(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				if fitem.Fid == "" {
					cid = fitem.Cid
				} else {
					cid = fitem.Fid
				}
			}
		}

		var paths []string
		if attr, err := client.FetchAttr(cid); err == nil {
			for _, path := range attr.Paths {
				if fmt.Sprintf("%v", path.Fid) == "0" {
					continue
				}
				paths = append(paths, path.Fname)
			}
		}

		prefix = filepath.Join(prefix, strings.Join(paths, "/"))

		data, err := client.ExportTree(cid)
		if err != nil {
			gol.Errorf("Export directory tree failed: %v\n", err)
			return
		}
		gol.Info("Export directory tree succeed")

		gol.Info("Start to convert directory tree to library")
		if err := pan.Tree2Lib(url, prefix, library, data); err != nil {
			gol.Errorf("Convert directory tree to library failed: %v\n", err)
			return
		}
		gol.Info("Convert directory tree to library succeed")

	},
}

// 获取指定ID文件的路径
var p115pathCmd = &cobra.Command{
	Use:   "path",
	Short: "get 115 file path",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P115.Cid
		}

		gol.Info("Start to get file attribute")

		client := pan.P115{
			Cookie: cookie,
		}

		if fileattr, err := client.FetchAttr(cid); err != nil {
			gol.Errorf("Get file path failed: %v\n", err)
		} else {
			var paths []string
			for _, p := range fileattr.Paths {
				paths = append(paths, p.Fname)
			}
			paths = append(paths, fileattr.Name)

			fmt.Printf("Get file path succeed: %s\n", strings.Join(paths, "/"))
		}

	},
}

// 获取指定路径文件的ID
var p115cidCmd = &cobra.Command{
	Use:   "cid",
	Short: "get 115 file cid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P115.Cpath
		}

		gol.Info("Start to get file cid")

		client := pan.P115{
			Cookie: cookie,
		}

		if fitem, err := client.FetchItem(cpath); err != nil {
			gol.Errorf("Get cid failed: %v\n", err)
		} else {
			var cid, fid string
			ftype := "文件"
			if fitem.Fid == "" {
				cid = fitem.Cid
				fid = fitem.Pid
				ftype = "目录"
			} else {
				cid = fitem.Fid
				fid = fitem.Cid
			}

			fmt.Printf("文件编号: %s\n", cid)
			fmt.Printf("文件名称: %s\n", fitem.Name)
			fmt.Printf("文件类型: %s\n", ftype)
			fmt.Printf("父级编号: %v\n", fid)
			fmt.Printf("文件哈希: %s\n", fitem.Sha)
			fmt.Printf("更新时间: %v\n", fitem.Te)
			fmt.Printf("文件提取码: %s\n", fitem.Pc)
		}

	},
}

// 获取指定文件属性
var p115attrCmd = &cobra.Command{
	Use:   "attr",
	Short: "get 115 file attribute",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P115.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P115.Cpath
		}

		gol.Info("Start to get file attribute")

		client := pan.P115{
			Cookie: cookie,
		}

		if cid == "" {
			if fitem, err := client.FetchItem(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				if fitem.Fid == "" {
					cid = fitem.Cid
				} else {
					cid = fitem.Fid
				}
			}
		}

		if fileattr, err := client.FetchAttr(cid); err != nil {
			gol.Errorf("Get file attribute failed: %v\n", err)
		} else {
			ftype := "文件"
			if fileattr.Category == "0" {
				ftype = "目录"
			}
			var paths []string
			for _, p := range fileattr.Paths {
				paths = append(paths, fmt.Sprintf("%s[%v]", p.Fname, p.Fid))
			}

			fmt.Printf("文件编号: %s\n", cid)
			fmt.Printf("文件名称: %s\n", fileattr.Name)
			fmt.Printf("文件类型: %s\n", ftype)
			fmt.Printf("文件数量: %v\n", fileattr.Count)
			fmt.Printf("目录数量: %v\n", fileattr.Fcount)
			fmt.Printf("文件大小: %s\n", fileattr.Size)
			fmt.Printf("播放时长: %d\n", fileattr.Plong)
			fmt.Printf("更新时间: %s\n", fileattr.Utime)
			fmt.Printf("文件描述: %s\n", fileattr.Desc)
			fmt.Printf("文件提取码: %s\n", fileattr.Pc)
			fmt.Printf("文件路径: %s\n", strings.Join(paths, "/"))
		}

	},
}

// 获取目录文件列表
var p115listCmd = &cobra.Command{
	Use:   "list",
	Short: "get 115 file list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cookie, _ := cmd.Flags().GetString("cookie")
		if cookie == "" {
			cookie = cfg.Conf.P115.Cookie
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P115.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P115.Cpath
		}

		gol.Info("Start to get file list")

		client := pan.P115{
			Cookie: cookie,
		}

		if cid == "" {
			if fitem, err := client.FetchItem(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				if fitem.Fid == "" {
					cid = fitem.Cid
				} else {
					cid = fitem.Fid
				}
			}
		}

		if flist, err := client.FetchList(cid); err != nil {
			gol.Errorf("Get file list failed: %v\n", err)
		} else {
			for _, fitem := range flist {
				var cid, fid string
				ftype := "文件"
				if fitem.Fid == "" {
					cid = fitem.Cid
					fid = fitem.Pid
					ftype = "目录"
				} else {
					cid = fitem.Fid
					fid = fitem.Cid
				}
				fmt.Printf("%s\t%s\t%s/%s\t%s\t%s\t%s\n", ftype, fitem.Name, cid, fid, fitem.Sha, fitem.Te, fitem.Pc)
			}
		}

	},
}

func init() {
	p115Cmd.PersistentFlags().StringP("cookie", "c", "", "pan cookie")
	p115Cmd.PersistentFlags().StringP("cid", "i", "", "pan target id")
	p115Cmd.PersistentFlags().StringP("cpath", "x", "", "pan target path")
	p115Cmd.PersistentFlags().StringP("fpath", "f", "", "directory tree file path")
	p115Cmd.PersistentFlags().StringP("base", "b", "", "media library directory path")
	p115Cmd.PersistentFlags().StringP("url", "u", "", "media url")
	p115Cmd.PersistentFlags().StringP("prefix", "p", "", "media url prefix path")

	p115Cmd.AddCommand(p115exportCmd)
	p115Cmd.AddCommand(p115treeifyCmd)
	p115Cmd.AddCommand(p115strmCmd)
	p115Cmd.AddCommand(p115attrCmd)
	p115Cmd.AddCommand(p115pathCmd)
	p115Cmd.AddCommand(p115cidCmd)
	p115Cmd.AddCommand(p115listCmd)

	rootCmd.AddCommand(p115Cmd)
}

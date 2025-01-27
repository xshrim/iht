package cmd

import (
	"fmt"
	"iht/pkg/cfg"
	"iht/pkg/pan"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xshrim/gol"
)

var p123Cmd = &cobra.Command{
	Use:   "p123",
	Short: "pan 123 commands",
	Long:  ``,
}

// 获取网盘token
var p123tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "get pan 123 token",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		gol.Info("Start to get pan 123 token")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
		}

		token, expiry, err := client.GetToken()
		if err != nil {
			gol.Error(err)
			return
		}

		gol.Info("Token: ", token)
		gol.Info("Expiry: ", expiry)

	},
}

// 获取指定文件路径
var p123pathCmd = &cobra.Command{
	Use:   "path",
	Short: "get 123 file path",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = cfg.Conf.P123.Token
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P123.Cid
		}

		gol.Info("Start to get file cid")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  token,
		}

		if fpath, err := client.FetchPath(cid); err != nil {
			gol.Errorf("Get cpath failed: %v\n", err)
		} else {
			fmt.Printf("文件路径: %s\n", fpath)
		}
	},
}

// 获取指定路径文件的ID
var p123cidCmd = &cobra.Command{
	Use:   "cid",
	Short: "get 123 file cid",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = cfg.Conf.P123.Token
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P123.Cpath
		}

		gol.Info("Start to get file cid")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  token,
		}

		if fobj, err := client.FetchObj(cpath); err != nil {
			gol.Errorf("Get Cid failed: %v\n", err)
			return
		} else {
			typ := "文件"
			if fobj.Type == 1 {
				typ = "目录"
			}

			var category string
			switch fobj.Category {
			case 0:
				category = "未知"
			case 1:
				category = "音频"
			case 2:
				category = "视频"
			case 3:
				category = "图片"
			}

			status := "正常"
			if fobj.Status > 100 {
				status = "驳回"
			}

			fmt.Printf("文件编号: %d\n", fobj.FileId)
			fmt.Printf("文件名称: %s\n", fobj.Filename)
			fmt.Printf("文件类型: %s\n", typ)
			fmt.Printf("文件大小: %d\n", fobj.Size)
			fmt.Printf("父级编号: %v\n", fobj.ParentFileId)
			fmt.Printf("文件哈希: %s\n", fobj.Etag)
			fmt.Printf("文件种类: %s\n", category)
			fmt.Printf("文件路径: %s\n", fobj.Path)
			fmt.Printf("文件状态: %s\n", status)
		}
	},
}

// 获取指定文件属性
var p123attrCmd = &cobra.Command{
	Use:   "attr",
	Short: "get 123 file attribute",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = cfg.Conf.P123.Token
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P123.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P123.Cpath
		}

		gol.Info("Start to get file attribute")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  token,
		}

		if cid == "" && cpath != "" {
			if fitem, err := client.FetchObj(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				cid = fmt.Sprintf("%d", fitem.FileId)
			}
		}

		if fobj, err := client.FetchAttr(cid); err != nil {
			gol.Errorf("Get cid failed: %v\n", err)
		} else {
			typ := "文件"
			if fobj.Type == 1 {
				typ = "目录"
			}

			status := "正常"
			if fobj.Status > 100 {
				status = "驳回"
			}

			fmt.Printf("文件编号: %d\n", fobj.FileID)
			fmt.Printf("文件名称: %s\n", fobj.Filename)
			fmt.Printf("文件类型: %s\n", typ)
			fmt.Printf("文件大小: %d\n", fobj.Size)
			fmt.Printf("父级编号: %v\n", fobj.ParentFileID)
			fmt.Printf("文件哈希: %s\n", fobj.Etag)
			fmt.Printf("文件路径: %s\n", fobj.Path)
			fmt.Printf("创建时间: %s\n", fobj.CreateAt)
			fmt.Printf("文件状态: %s\n", status)
		}

	},
}

// 获取目录文件列表
var p123listCmd = &cobra.Command{
	Use:   "list",
	Short: "get 123 file list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = cfg.Conf.P123.Token
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P123.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P123.Cpath
		}

		keyword, _ := cmd.Flags().GetString("keyword")
		if keyword == "" {
			keyword = cfg.Conf.P123.Keyword
		}

		gol.Info("Start to get file list")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  token,
		}

		if cid == "" && cpath != "" {
			if fitem, err := client.FetchObj(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				cid = fmt.Sprintf("%d", fitem.FileId)
			}
		}

		if cid == "" {
			cid = "0"
		}
		params := []string{cid, keyword}

		if flist, err := client.FetchList(params...); err != nil {
			gol.Errorf("Get file list failed: %v\n", err)
		} else {
			for _, fobj := range flist {
				var typ, category string
				switch fobj.Type {
				case 0:
					typ = "文件"
				case 1:
					typ = "目录"
				}

				switch fobj.Category {
				case 0:
					category = "未知"
				case 1:
					category = "音频"
				case 2:
					category = "视频"
				case 3:
					category = "图片"
				}

				status := "正常"
				if fobj.Status > 100 {
					status = "驳回"
				}

				fmt.Printf("%s\t%s\t%s\t%d/%d\t%d\t%s\t%s\t%s\n", typ, fobj.Filename, fobj.Path, fobj.FileId, fobj.ParentFileId, fobj.Size, fobj.Etag, category, status)
			}
		}

	},
}

// 文件重命名
var p123renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename 123 files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			id = cfg.Conf.P123.Id
		}

		secret, _ := cmd.Flags().GetString("secret")
		if secret == "" {
			secret = cfg.Conf.P123.Secret
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = cfg.Conf.P123.Token
		}

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P123.Cid
		}

		cpath, _ := cmd.Flags().GetString("cpath")
		if cpath == "" {
			cpath = cfg.Conf.P123.Cpath
		}

		cname, _ := cmd.Flags().GetString("cname")
		if cname == "" {
			cname = cfg.Conf.P123.Cname
		}

		rlist, _ := cmd.Flags().GetStringSlice("rlist")
		if len(rlist) == 0 {
			rlist = cfg.Conf.P123.Rlist
		}

		gol.Info("Start to rename files")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  token,
		}

		if cid == "" && cpath != "" {
			if fitem, err := client.FetchObj(cpath); err != nil {
				gol.Errorf("Get Cid failed: %v\n", err)
				return
			} else {
				cid = fmt.Sprintf("%d", fitem.FileId)
			}
		}

		if cid != "" && cname != "" {
			rlist = append(rlist, fmt.Sprintf("%s|%s", cid, cname))
		}

		var rnlist []string
		for _, ritem := range rlist {
			if strings.Contains(ritem, "|") {
				items := strings.Split(ritem, "|")
				if len(items) == 2 {
					cid = items[0]
					if _, err := strconv.Atoi(cid); err != nil {
						if fitem, err := client.FetchObj(cid); err != nil {
							continue
						} else {
							cid = fmt.Sprintf("%d", fitem.FileId)
						}
					}
					cname = items[1]

					rnlist = append(rnlist, fmt.Sprintf("%s|%s", cid, cname))
				}
			}
		}

		result := client.RenameList(rlist)
		fmt.Printf("%d个文件重命名成功, %d个文件重命名失败\n", len(result["succeed"]), len(result["failed"]))
	},
}

func init() {
	p123Cmd.PersistentFlags().StringP("id", "u", "", "pan client id")
	p123Cmd.PersistentFlags().StringP("secret", "s", "", "pan client secret")
	p123Cmd.PersistentFlags().StringP("cid", "i", "", "pan target id")
	p123Cmd.PersistentFlags().StringP("cpath", "x", "", "pan target path")
	p123Cmd.PersistentFlags().StringP("cname", "n", "", "new file name")
	p123Cmd.PersistentFlags().StringSliceP("rlist", "r", []string{}, "rename file list")
	p123Cmd.PersistentFlags().StringP("keyword", "k", "", "pan target keyword")

	p123Cmd.AddCommand(p123tokenCmd)
	p123Cmd.AddCommand(p123pathCmd)
	p123Cmd.AddCommand(p123cidCmd)
	p123Cmd.AddCommand(p123attrCmd)
	p123Cmd.AddCommand(p123listCmd)
	p123Cmd.AddCommand(p123renameCmd)

	rootCmd.AddCommand(p123Cmd)
}

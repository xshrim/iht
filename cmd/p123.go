package cmd

import (
	"fmt"
	"iht/pkg/cfg"
	"iht/pkg/pan"

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

		cid, _ := cmd.Flags().GetString("cid")
		if cid == "" {
			cid = cfg.Conf.P123.Cid
		}

		gol.Info("Start to get file list")

		client := &pan.P123{
			Id:     id,
			Secret: secret,
			Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzc3MDE4NjYsImlhdCI6MTczNzA5NzA2NiwiaWQiOjE4NDQwNjQzNjYsIm1haWwiOiJ4c2hyaW1AeWVhaC5uZXQiLCJuaWNrbmFtZSI6InhzaHJpbSIsInVzZXJuYW1lIjoxOTMyNzQ0OTM2OSwidiI6MH0.30ObtPN0PvpVhN8YhLL2XxZy7kWKpCq_KNfkKTOZmY0",
			Expiry: "2025-02-16T14:57:46+08:00",
		}

		if flist, err := client.FetchList(cid); err != nil {
			gol.Errorf("Get file list failed: %v\n", err)
		} else {
			for _, fobj := range flist {
				var ftype, fcategory string
				switch fobj.Type {
				case 0:
					ftype = "文件"
				case 1:
					ftype = "目录"
				}

				switch fobj.Category {
				case 0:
					fcategory = "未知"
				case 1:
					fcategory = "音频"
				case 2:
					fcategory = "视频"
				case 3:
					fcategory = "图片"
				}

				status := "正常"
				if fobj.Status > 100 {
					status = "驳回"
				}

				fmt.Printf("%s\t%s\t%d/%d\t%d\t%s\t%s\t%s\n", ftype, fobj.Filename, fobj.FileId, fobj.ParentFileId, fobj.Size, fobj.Etag, fcategory, status)
			}
		}

	},
}

func init() {
	p123Cmd.PersistentFlags().StringP("id", "u", "", "pan client id")
	p123Cmd.PersistentFlags().StringP("secret", "s", "", "pan client secret")
	p123Cmd.PersistentFlags().StringP("cid", "i", "", "pan target id")

	p123Cmd.AddCommand(p123tokenCmd)
	p123Cmd.AddCommand(p123listCmd)

	rootCmd.AddCommand(p123Cmd)
}

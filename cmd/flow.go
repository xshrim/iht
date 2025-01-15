package cmd

import (
	"fmt"
	"iht/pkg/cfg"
	"iht/pkg/flow"

	"github.com/spf13/cobra"
	"github.com/xshrim/gol/tk"
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "flow commands",
	Long:  ``,
}

// 列出工作流步骤
var flowlistCmd = &cobra.Command{
	Use:   "list",
	Short: "list flow steps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fdata, _ := cmd.Flags().GetString("file")
		if fdata == "" {
			fdata = tk.Jsonify(cfg.Conf.Flow)
		}

		f, err := flow.Load(fdata)
		if err != nil {
			fmt.Println("load flow failed:", err)
			return
		}

		fmt.Println(f.List())

	},
}

// 运行工作流
var flowrunCmd = &cobra.Command{
	Use:   "run",
	Short: "run flow",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		data, _ := cmd.Flags().GetStringSlice("data")

		fdata, _ := cmd.Flags().GetString("file")
		if fdata == "" {
			fdata = tk.Jsonify(cfg.Conf.Flow)
		}

		f, err := flow.Load(fdata)
		if err != nil {
			fmt.Println("load flow failed:", err)
			return
		}

		res, err := f.Run(data)
		if err != nil {
			fmt.Println("run flow failed:", err)
			return
		}

		for _, item := range res {
			fmt.Println(item)
		}
	},
}

func init() {
	flowCmd.PersistentFlags().StringP("file", "f", "", "flow file path")
	flowCmd.PersistentFlags().StringSliceP("data", "d", []string{}, "flow target data")

	flowCmd.AddCommand(flowlistCmd)
	flowCmd.AddCommand(flowrunCmd)

	rootCmd.AddCommand(flowCmd)
}

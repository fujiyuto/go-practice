/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	

	"my-cli/internal/jwt"
	"github.com/spf13/cobra"
)

// jwtCmd represents the jwt command
var (
	filename string
	generate bool
	analyze  bool

	jwtCmd = &cobra.Command{
		Use:   "jwt",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return jwt.Run(filename, generate, analyze)
		},
	}
)

func init() {
	rootCmd.AddCommand(jwtCmd)

	flags := jwtCmd.Flags()

	// ファイル名
	flags.StringVarP(&filename, "filename", "f", "jwt.json", "Read file name.")
	// jwtの生成
	flags.BoolVarP(&generate, "generate", "g", false, "Generate JWT Token.")
	// jwtの解析
	flags.BoolVarP(&analyze, "analyze", "a", false, "Analyze JWT Token.")
	
}

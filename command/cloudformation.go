package command

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yawn/alm/cloudformation"
	"github.com/yawn/alm/file"
)

var cloudformationCmd = &cobra.Command{

	Use:   "cloudformation",
	Short: "Get last messages from CloudFormation",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()

		cfg, err := config.LoadDefaultConfig(ctx)

		if err != nil {
			return errors.Wrapf(err, "failed to load default config")
		}

		writer := file.New("/tmp/alm")

		regions := []string{
			"eu-central-1",
			"eu-west-1",
		}

		for _, region := range regions {

			go func(cfg aws.Config, region string, writer *file.Writer) {

				cfg = cfg.Copy()
				cfg.Region = region

				c := cloudformation.New(cfg)

				for {

					targets, err := c.Discover(ctx, 7*24*time.Hour)

					if err != nil {
						log.Printf("[ERROR] %s", err)
					}

					for _, target := range targets {

						if err := c.Log(ctx, writer, target); err != nil {
							log.Printf("[ERROR] %s", err)
						}

					}

					time.Sleep(10 * time.Second)

				}

			}(cfg, region, writer)

		}

		select {}

	},
}

func init() {
	rootCmd.AddCommand(cloudformationCmd)
}

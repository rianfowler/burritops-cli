package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Clone code, build a container image, and publish it using Dagger",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Connect to the Dagger engine (which leverages Docker under the hood)
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
		if err != nil {
			log.Fatalf("failed to connect to dagger: %v", err)
		}
		defer client.Close()

		// Clone the repository using Dagger's Git functionality.
		repo := client.Git("https://github.com/rianfowler/burritops-cli").Branch("main")
		// Convert the Git reference to a directory using Tree().
		srcDir := repo.Tree()

		// Create a container that mounts the directory.
		container := client.Container().
			From("docker.io/library/golang:1.24").
			WithMountedDirectory("/src", srcDir).
			WithWorkdir("/src")

		// Build the image using srcDir as the build context.
		image := container.Build(srcDir)

		id, err := image.ID(ctx)
		if err != nil {
			log.Fatalf("failed to retrieve image ID: %v", err)
		}
		fmt.Println("Build successful! Image ID:", id)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

// cmd/packbuild.go
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var packBuildCmd = &cobra.Command{
	Use:   "pack-build",
	Short: "Build your app using the pack CLI in daemonless mode",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
		if err != nil {
			log.Fatalf("failed to connect to dagger: %v", err)
		}
		defer client.Close()

		// Clone your repository and get its directory tree.
		repo := client.Git("https://github.com/rianfowler/burritops-cli").Branch("main")
		srcDir := repo.Tree()

		// Create a container based on the pack CLI image.
		// Note: No Docker socket mount is needed in daemonless mode.
		container := client.Container().
			From("buildpacksio/pack:latest").
			WithMountedDirectory("/workspace", srcDir).
			WithWorkdir("/workspace").
			WithEnvVariable("PACK_LOG_LEVEL", "debug") // Optional for extra logging.

		// Execute the pack CLI command in daemonless mode.
		// The --exporter=tar flag tells pack to output an OCI image tarball.
		execContainer := container.WithExec([]string{
			"pack", "build", "burritops-app",
			"--builder", "paketobuildpacks/builder:base",
			"--path", ".",
			"--exporter=tar",
		})

		// Capture the combined output.
		output, err := execContainer.Stdout(ctx)
		if err != nil {
			stderr, _ := execContainer.Stderr(ctx)
			log.Fatalf("failed to execute pack build command: %v\nStderr output:\n%s", err, stderr)
		}

		fmt.Println("Pack build output:")
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(packBuildCmd)
}

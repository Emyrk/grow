package cmd

import "context"

func Run() {
	ctx := context.Background()
	_ = RootCmd.ExecuteContext(ctx)
}

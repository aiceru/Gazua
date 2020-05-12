package main

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		log.Fatal("error while google.DefaultClient()")
	}
	computeService, err := compute.New(client)

}

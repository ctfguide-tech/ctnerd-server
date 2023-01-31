package main

import (
    "context"
    "log"

    "github.com/containerd/containerd"
    "github.com/containerd/containerd/cio"
    "github.com/containerd/containerd/namespaces"
    "github.com/kataras/iris/v12"
)

const DEFAULT = "docker.io/library/redis:latest"

func main() {
    app := iris.Default()

    client, err := containerd.New("/run/containerd/containerd.sock")
    if err != nil {
        log.Fatal("Failed to start container daemon")
    }

	defer client.Close()

	app.Post("/containers", func(ctx iris.Context) {
        imageName := ctx.URLParamDefault("image", DEFAULT)

        // Pull Image
        image, err := client.Pull(context, imageName)
        if err != nil {
            ctx.Writef("Failed to pull image %s", imageName)
        }

        redis, err := client.NewContainer(context, image)
        if err != nil {
            ctx.Writef("Failed to create container %s", imageName)
        }
        ctx.Writef("Container Created!")
    })

    app.Get("/containers/{container}", func(ctx iris.Context) {
        name := ctx.Params().Get("container")
        ctx.Writef("Container Name: %s", name)
    })

    app.Listen(":8080")
}

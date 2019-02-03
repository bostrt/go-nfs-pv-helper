// nfspv pv-001 --server=example.com --size=1Gi --path=/mnt/nfs/test --reclaim-policy=Retain

package main

import (
	"log"
	"os"
	"text/template"

	"github.com/urfave/cli"
)

type PV struct {
	Name          string
	Server        string
	Size          string
	Path          string
	ReclaimPolicy string
}

func main() {
	app := cli.NewApp()
	app.Name = "OpenShift NFS PersistentVolume Helper"
	app.Usage = ""
	app.UsageText = "go-nfs-pv-helper [global options]"
	app.Version = "0.0.1"
	app.Description = "Generates PersistentVolume YAML from arguments"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "new volume name",
		},
		cli.StringFlag{
			Name:  "server",
			Usage: "NFS server hostname",
		},
		cli.StringFlag{
			Name:  "size",
			Usage: "size of volume",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "NFS mount path",
		},
		cli.StringFlag{
			Name:  "reclaim-policy",
			Usage: "reclaim policy for volume",
			Value: "Retain",
		},
	}

	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {

	name := c.String("name")
	server := c.String("server")
	size := c.String("size")
	path := c.String("path")
	reclaimPolicy := c.String("reclaim-policy")

	if name == "" {
		log.Fatal("PV name is required")
	}
	if size == "" {
		log.Fatal("PV size is required (e.g. 1Gi, 500Mi)")
	}

	if server == "" {
		log.Fatal("NFS server hostname required")
	}

	if path == "" {
		log.Fatal("NFS mount path is required")
	}

	if reclaimPolicy == "" {
		log.Fatal("reclaim policy is required (e.g. Retain, Recycle)")
	}

	pv := PV{
		Name:          name,
		Server:        server,
		Size:          size,
		Path:          path,
		ReclaimPolicy: reclaimPolicy,
	}

	doTemplate(pv)

	return nil
}

func doTemplate(pv PV) {
	tpl, err := template.ParseFiles("pv-template.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(os.Stdout, pv)
	if err != nil {
		log.Fatal(err)
	}
}

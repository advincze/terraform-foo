package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

type res struct {
	Must   string
	Option string
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"foo_res": &schema.Resource{
						Create: func(d *schema.ResourceData, m interface{}) error {
							log.Println("CREATE", d)
							id := "foo" + time.Now().Format("20060102150405")
							d.SetId(id)
							f, err := os.Create(id)
							if err != nil {
								return err
							}
							defer f.Close()

							return json.NewEncoder(f).Encode(res{
								Must:   d.Get("must").(string),
								Option: d.Get("option").(string),
							})
						},
						Read: func(d *schema.ResourceData, m interface{}) error {
							log.Println("READ", d)
							f, err := os.Open(d.State().ID)
							if err != nil {
								return err
							}
							defer f.Close()

							var r res
							err = json.NewDecoder(f).Decode(&r)
							if err != nil {
								return err
							}

							d.Set("must", r.Must)
							d.Set("option", r.Option)
							return nil
						},
						Update: func(d *schema.ResourceData, m interface{}) error {
							log.Println("UPDATE", d)
							f, err := os.OpenFile(d.State().ID, os.O_RDWR|os.O_TRUNC, 0666)
							if err != nil {
								return err
							}
							defer f.Close()

							return json.NewEncoder(f).Encode(res{
								Must:   d.Get("must").(string),
								Option: d.Get("option").(string),
							})
						},
						Delete: func(d *schema.ResourceData, m interface{}) error {
							log.Println("DELETE", d)
							return os.Remove(d.State().ID)
						},
						Exists: func(d *schema.ResourceData, m interface{}) (bool, error) {
							_, err := os.Stat(d.State().ID)
							if err != nil {
								if os.IsNotExist(err) {
									d.SetId("")
									return false, nil
								}
								return false, err
							}
							return true, nil
						},
						Schema: map[string]*schema.Schema{
							"must": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
							"option": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			}
		},
	})
}

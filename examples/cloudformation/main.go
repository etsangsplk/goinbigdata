package main

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
	"text/template"

	"github.com/upitau/goinbigdata/examples/cloudformation/s3"
)

func parseCfnTemplate(name, templateFile string, specTags *map[string]string) (string, error) {
	fmap := template.FuncMap{
		"Key":   getTagKey,
		"Value": getTagValue,
	}
	t := template.Must(template.New(name).Funcs(fmap).Parse(templateFile))

	var output bytes.Buffer
	err := generateCfnTemplate(t, &output, specTags)
	if err != nil {
		return "", err
	}
	return output.String(), err
}

type Tag struct {
	TagKey   string
	TagValue string
}

type Tags struct {
	Tags []Tag
}

func generateCfnTemplate(t *template.Template, output io.Writer, specTags *map[string]string) error {
	err := t.Execute(output, createDataObject(specTags))
	if err != nil {
		return err
	}
	return nil
}

func createDataObject(specTags *map[string]string) Tags {
	var tags Tags
	if specTags == nil {
		return tags
	}
	for k, v := range *specTags {
		tags.Tags = append(tags.Tags, Tag{TagKey: k, TagValue: v})
	}
	fmt.Printf("tags : %v", tags)
	return tags
}

func getTagKey(t Tag) string {
	return t.TagKey
}

func getTagValue(t Tag) string {
	return t.TagValue
}

func main() {
	specTags := map[string]string{"k8version": "1.8.0", "minikubeversion": "0.25.0"}
	//specTags := *new(map[string]string)
	out, err := parseCfnTemplate("s3CfnTemplate", s3.S3CfnTemplate, &specTags)
	fmt.Printf("output: %v error: %v", out, err)
}

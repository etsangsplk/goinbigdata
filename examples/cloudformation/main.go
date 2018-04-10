package main

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
	"text/template"

	"github.com/upitau/goinbigdata/examples/cloudformation/kinesis"
	//"github.com/upitau/goinbigdata/examples/cloudformation/s3"
	"github.com/upitau/goinbigdata/examples/cloudformation/tags"
)

func parseCfnTemplate(name, templateFile string, specTags *map[string]string) (string, error) {
	fmap := template.FuncMap{
		"Key":   tags.GetTagKey,
		"Value": tags.GetTagValue,
	}
	t := template.Must(template.New(name).Funcs(fmap).Parse(templateFile))

	var output bytes.Buffer
	err := generateCfnTemplate(t, &output, specTags)
	if err != nil {
		return "", err
	}
	return output.String(), err
}

func generateCfnTemplate(t *template.Template, output io.Writer, specTags *map[string]string) error {
	err := t.Execute(output, kinesis.CreateDataObject(specTags))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	specTags := map[string]string{"k8version": "1.8.0", "minikubeversion": "0.25.0"}
	//specTags := *new(map[string]string)
	// out, err := parseCfnTemplate("s3CfnTemplate", s3.S3CfnTemplate, &specTags)
	out, err := parseCfnTemplate("kinesisCfnTemplate", kinesis.KinesisCfnTemplate, &specTags)
	fmt.Printf("output: %v error: %v", out, err)

}

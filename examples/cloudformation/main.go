package main

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
	"text/template"

	//"github.com/aws/aws-sdk-go/aws"

	"github.com/upitau/goinbigdata/examples/cloudformation/kinesis"
	//"github.com/upitau/goinbigdata/examples/cloudformation/s3"
	"github.com/upitau/goinbigdata/examples/cloudformation/tags"
)

type cfnTemplate struct {
	dataObject   interface{}
	t            *template.Template
	templateBody string
}

func NewTemplateTransform(name, templateBody string, dataObject interface{}) *cfnTemplate {
	return &cfnTemplate{
		dataObject:   dataObject,
		t:            template.New(name),
		templateBody: templateBody,
	}
}

func (c *cfnTemplate) expandTemplateTags() *cfnTemplate {
	fmap := template.FuncMap{
		"Key":   tags.GetTagKey,
		"Value": tags.GetTagValue,
	}
	c.t = template.Must(c.t.Funcs(fmap).Parse(c.templateBody))
	return c
}

func (c *cfnTemplate) transform() (string, error) {
	var output bytes.Buffer
	err := c.evaluateTemplate(&output)
	if err != nil {
		return "", err
	}
	return output.String(), err
}

func (c *cfnTemplate) evaluateTemplate(output io.Writer) error {
	err := c.t.Execute(output, c.dataObject)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	specTags := map[string]string{"k8version": "1.8.0", "minikubeversion": "0.25.0"}
	//specTags := *new(map[string]string)
	// out, err := parseCfnTemplate("s3CfnTemplate", s3.S3CfnTemplate, &specTags)
	// getDataObject := kinesis.CreateDataObject
	// out, err := parseCfnTemplate("kinesisCfnTemplate", kinesis.KinesisCfnTemplate, getDataObject(&specTags))
	// fmt.Printf("output: %v error: %v", out, err)

	getDataObject := kinesis.CreateDataObject(&specTags)
	tpl, err := NewTemplateTransform("kinesisCfnTemplate", kinesis.KinesisCfnTemplate, getDataObject).
		expandTemplateTags().
		transform()
	//out := aws.String(tpl)
	fmt.Printf("output: %v error: %v", tpl, err)
}

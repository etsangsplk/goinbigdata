package cloudformation

import (
	"bytes"
	"io"
	"text/template"
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

func (c *cfnTemplate) expandTemplateTags(specTags *map[string]string) (string, error) {
	fmap := template.FuncMap{
		"Key":   getTagKey,
		"Value": getTagValue,
	}
	c.t = template.Must(c.t.Funcs(fmap).Parse(c.templateBody))

	return c.transform(specTags)
}

func (c *cfnTemplate) transform(specTags *map[string]string) (string, error) {
	var output bytes.Buffer
	err := c.evaluateTemplate(&output, specTags)
	if err != nil {
		return "", err
	}
	return output.String(), err
}

func (c *cfnTemplate) evaluateTemplate(output io.Writer, specTags *map[string]string) error {
	err := c.t.Execute(output, c.dataObject)
	if err != nil {
		return err
	}
	return nil
}

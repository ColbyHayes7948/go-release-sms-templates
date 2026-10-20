package main

import (
	"fmt"
	"time"
)

type Release struct{ Name, Phone string }

func ValidateRelease(r Release) error {
	if r.Name == "" || r.Phone == "" {
		return fmt.Errorf("release name and phone are required")
	}
	return nil
}

func PublishRelease(c *InfraiClient, r Release) error {
	if err := ValidateRelease(r); err != nil {
		return err
	}
	namespace := fmt.Sprintf("devtools-%s-%d", r.Name, time.Now().UnixNano())
	signature := "DevTools-" + namespace
	if err := c.CreateSignature(SignatureRequest{Name: signature}); err != nil {
		return err
	}
	template, err := c.CreateTemplate(TemplateRequest{Name: namespace, Body: "Release {{name}} is ready"})
	if err != nil {
		return err
	}
	return c.SendSMS(SendRequest{To: r.Phone, TemplateID: template.ID, TemplateVars: map[string]string{"name": r.Name}, From: signature})
}

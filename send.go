package sendinblue

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Address describes an email address
type Address struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Attachment describes an attachment to an email message
type Attachment struct {

	// Name is the filename to be asserted for this attachment
	Name string `json:"name"`

	// URL indicates an external reference from which the attachment content should be retrieved
	URL string `json:"url"`

	// Content declares the inline content of the attachment, encoded as a Base64 string
	Content string `json:"content"`
}

// InlineAttachment returns a new Attachment from a byte-wise reader source.  The content will be converted to a Base64 string inside the Attachment.
func InlineAttachment(name string, in io.Reader) (ret Attachment, err error) {
	ret.Name = name

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)

	_, err = io.Copy(enc, in)
	if err == nil {
		ret.Content = buf.String()
	}
	return
}

// Message describes an email message which should be sent
type Message struct {

	// Sender is the entity which is sending the email message.  REQUIRED.
	Sender *Address `json:"sender"`

	// To is the list of primary recipients of the email
	To []*Address `json:"to,omitempty"`

	// Bcc (blind carbon copy) is the list of recipients of the email which should not be disclosed to other recipients
	Bcc []*Address `json:"bcc,omitempty"`

	// Cc (carbon copy) is the list of secondary recipients of the email
	Cc []*Address `json:"cc,omitempty"`

	// HTMLContent is the HTML-formatted content of the email
	HTMLContent string `json:"htmlContent,omitempty"`

	// TextContent is the plain-text content of the email
	TextContent string `json:"textContent,omitempty"`

	// Subject is the subject of the email
	Subject string `json:"subject,omitempty"`

	// ReplyTo indicates that replies to this email should be sent to this address
	ReplyTo *Address `json:"replyTo,omitempty"`

	// Attachments describe any attachments which should be added to this email
	Attachments []*Attachment `json:"attachment,omitempty"` // documentation indicates attachment (singular) even though multiple attachments are allowed

	// Headers is the list of email headers which should be sent with the email message
	Headers map[string]string `json:"headers,omitempty"`

	// TemplateID indicates that the content of the email address should be taken from the indicated template instead of directly-included content
	TemplateID int64 `json:"templateId,omitempty"`

	// Params is the list of parameters which should be used to populate the template
	Params map[string]string `json:"params,omitempty"`

	// Tags are arbitrary labels which are applied to this email in order to facilitate organizational operations in SendInBlue
	Tags []string `json:"tags,omitempty"`
}

// Send transmits the email message to SendInBlue
func (m *Message) Send(apiKey string) error {
	url := "https://api.sendinblue.com/v3/smtp/email"

	data, err := json.Marshal(m)
	if err != nil {
		return errors.New("failed to encode message: " + err.Error())
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Add("api-key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("failed to transmit message: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("send failed: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

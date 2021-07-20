package codeship

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEncrypted() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEncryptedCreate,
		ReadContext:   noop,
		UpdateContext: resourceEncryptedCreate,
		DeleteContext: noop,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"aes_key": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"content": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"encrypted_content": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func noop(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceEncryptedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	plaintext := []byte(d.Get("content").(string))
	key, err := base64.StdEncoding.DecodeString(d.Get("aes_key").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return diag.FromErr(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return diag.FromErr(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	output := base64.StdEncoding.EncodeToString(ciphertext)

	d.Set("encrypted_content", output)
	d.SetId(output)

	return nil
}

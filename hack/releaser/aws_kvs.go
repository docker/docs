package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	kvstypes "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore/types"
)

// updateKeysBatchSize is the maximum number of put+delete items per
// UpdateKeys call. The service rejects larger batches.
const updateKeysBatchSize = 50

type AwsCloudfrontFunctionUpdateCmd struct {
	Region                string `kong:"name='region',env='AWS_REGION'"`
	FunctionName          string `kong:"name='function-name',env='AWS_CLOUDFRONT_FUNCTION_NAME'"`
	FunctionFile          string `kong:"name='function-file',env='AWS_CLOUDFRONT_FUNCTION_FILE'"`
	KvsARN                string `kong:"name='kvs-arn',env='AWS_CLOUDFRONT_KVS_ARN'"`
	RedirectsFile         string `kong:"name='redirects-file',env='REDIRECTS_FILE'"`
	RedirectsPrefixesFile string `kong:"name='redirects-prefixes-file',env='REDIRECTS_PREFIXES_FILE'"`
	DryRun                bool   `kong:"name='dry-run',env='DRY_RUN'"`
}

func (s *AwsCloudfrontFunctionUpdateCmd) Run() error {
	ctx := context.Background()

	desired, err := loadDesiredRedirects(s.RedirectsFile)
	if err != nil {
		return fmt.Errorf("load redirects: %w", err)
	}
	log.Printf("INFO: loaded %d redirect entries from %s", len(desired), s.RedirectsFile)

	funcCode, err := renderFunctionCode(s.FunctionFile, s.RedirectsPrefixesFile)
	if err != nil {
		return fmt.Errorf("render function code: %w", err)
	}

	if s.DryRun {
		log.Printf("INFO: dry run. Region=%s FunctionName=%s KvsARN=%s redirects=%d",
			s.Region, s.FunctionName, s.KvsARN, len(desired))
		log.Printf("INFO: function code (%d bytes):\n%s", len(funcCode), funcCode)
		return nil
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.Region),
		config.WithCredentialsProvider(awsV2Credentials()),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	if err := syncKVS(ctx, cloudfrontkeyvaluestore.NewFromConfig(cfg), s.KvsARN, desired); err != nil {
		return fmt.Errorf("sync KVS: %w", err)
	}

	if err := updateFunction(ctx, cloudfront.NewFromConfig(cfg), s.FunctionName, funcCode); err != nil {
		return fmt.Errorf("update function: %w", err)
	}

	return nil
}

// loadDesiredRedirects reads redirects.json (map[alias]target) and
// normalizes keys to match the function's lookup form: trailing slashes
// stripped, empty keys dropped.
func loadDesiredRedirects(path string) (map[string]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var in map[string]string
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		nk := normalizeKey(k)
		if nk == "" {
			continue
		}
		out[nk] = v
	}
	return out, nil
}

func normalizeKey(k string) string {
	for len(k) > 1 && k[len(k)-1] == '/' {
		k = k[:len(k)-1]
	}
	if k == "/" {
		return ""
	}
	return k
}

func renderFunctionCode(funcFile, prefixesFile string) (string, error) {
	tplBytes, err := os.ReadFile(funcFile)
	if err != nil {
		return "", err
	}
	prefixesBytes, err := os.ReadFile(prefixesFile)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(string(tplBytes))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct{ RedirectsPrefixesJSON string }{
		RedirectsPrefixesJSON: string(prefixesBytes),
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func syncKVS(ctx context.Context, svc *cloudfrontkeyvaluestore.Client, kvsARN string, desired map[string]string) error {
	desc, err := svc.DescribeKeyValueStore(ctx, &cloudfrontkeyvaluestore.DescribeKeyValueStoreInput{
		KvsARN: aws.String(kvsARN),
	})
	if err != nil {
		return fmt.Errorf("describe KVS: %w", err)
	}
	etag := *desc.ETag

	current := map[string]string{}
	pager := cloudfrontkeyvaluestore.NewListKeysPaginator(svc, &cloudfrontkeyvaluestore.ListKeysInput{
		KvsARN: aws.String(kvsARN),
	})
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("list KVS keys: %w", err)
		}
		for _, item := range page.Items {
			current[aws.ToString(item.Key)] = aws.ToString(item.Value)
		}
	}
	log.Printf("INFO: KVS currently holds %d keys; desired %d", len(current), len(desired))

	var puts []kvstypes.PutKeyRequestListItem
	for k, v := range desired {
		if cur, ok := current[k]; !ok || cur != v {
			puts = append(puts, kvstypes.PutKeyRequestListItem{
				Key:   aws.String(k),
				Value: aws.String(v),
			})
		}
	}
	var deletes []kvstypes.DeleteKeyRequestListItem
	for k := range current {
		if _, ok := desired[k]; !ok {
			deletes = append(deletes, kvstypes.DeleteKeyRequestListItem{
				Key: aws.String(k),
			})
		}
	}
	// Stable order for reproducible logs.
	sort.Slice(puts, func(i, j int) bool { return *puts[i].Key < *puts[j].Key })
	sort.Slice(deletes, func(i, j int) bool { return *deletes[i].Key < *deletes[j].Key })

	log.Printf("INFO: KVS diff: %d puts, %d deletes", len(puts), len(deletes))
	if len(puts) == 0 && len(deletes) == 0 {
		return nil
	}

	// Batch puts and deletes together, up to updateKeysBatchSize per call.
	for len(puts) > 0 || len(deletes) > 0 {
		batchPuts, batchDeletes := []kvstypes.PutKeyRequestListItem{}, []kvstypes.DeleteKeyRequestListItem{}
		remaining := updateKeysBatchSize

		take := min(remaining, len(puts))
		batchPuts = puts[:take]
		puts = puts[take:]
		remaining -= take

		take = min(remaining, len(deletes))
		batchDeletes = deletes[:take]
		deletes = deletes[take:]

		out, err := svc.UpdateKeys(ctx, &cloudfrontkeyvaluestore.UpdateKeysInput{
			KvsARN:  aws.String(kvsARN),
			IfMatch: aws.String(etag),
			Puts:    batchPuts,
			Deletes: batchDeletes,
		})
		if err != nil {
			return fmt.Errorf("update KVS keys: %w", err)
		}
		etag = *out.ETag
		log.Printf("INFO: applied batch (puts=%d deletes=%d)", len(batchPuts), len(batchDeletes))
	}
	return nil
}

func updateFunction(ctx context.Context, svc *cloudfront.Client, name, code string) error {
	desc, err := svc.DescribeFunction(ctx, &cloudfront.DescribeFunctionInput{
		Name:  aws.String(name),
		Stage: cftypes.FunctionStageDevelopment,
	})
	if err != nil {
		return fmt.Errorf("describe function: %w", err)
	}

	// Compare against currently-published code so we skip republishing
	// when nothing has changed.
	get, err := svc.GetFunction(ctx, &cloudfront.GetFunctionInput{
		Name:  aws.String(name),
		Stage: cftypes.FunctionStageLive,
	})
	if err == nil && bytes.Equal(get.FunctionCode, []byte(code)) {
		log.Printf("INFO: function %q LIVE code unchanged (sha256=%s); skipping publish",
			name, sha256hex(code))
		return nil
	}

	up, err := svc.UpdateFunction(ctx, &cloudfront.UpdateFunctionInput{
		Name:           aws.String(name),
		IfMatch:        desc.ETag,
		FunctionCode:   []byte(code),
		FunctionConfig: desc.FunctionSummary.FunctionConfig,
	})
	if err != nil {
		return fmt.Errorf("update function: %w", err)
	}
	log.Printf("INFO: function %q updated (sha256=%s)", name, sha256hex(code))

	if _, err := svc.PublishFunction(ctx, &cloudfront.PublishFunctionInput{
		Name:    aws.String(name),
		IfMatch: up.ETag,
	}); err != nil {
		return fmt.Errorf("publish function: %w", err)
	}
	log.Printf("INFO: function %q published to LIVE", name)
	return nil
}

func sha256hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func awsV2Credentials() aws.CredentialsProvider {
	return credentials.NewStaticCredentialsProvider(
		getEnvOrSecret("AWS_ACCESS_KEY_ID"),
		getEnvOrSecret("AWS_SECRET_ACCESS_KEY"),
		getEnvOrSecret("AWS_SESSION_TOKEN"),
	)
}
